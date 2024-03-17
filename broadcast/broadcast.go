package broadcast

import "context"

type BroadcastServer interface {
	Subscribe() <-chan string
	CancelSubscription(<-chan string)
}

type broadcastServer struct {
	source         <-chan string
	listeners      []chan string
	addListener    chan chan string
	removeListener chan (<-chan string)
}

func (s *broadcastServer) Subscribe() <-chan string {
	newListener := make(chan string, 100)
	s.addListener <- newListener
	return newListener
}

func (s *broadcastServer) CancelSubscription(channel <-chan string) {
	s.removeListener <- channel
}

func (s *broadcastServer) serve(ctx context.Context) {
	defer func() {
		for _, listener := range s.listeners {
			if listener != nil {
				close(listener)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case newListener := <-s.addListener:
			s.listeners = append(s.listeners, newListener)
		case listenerToRemove := <-s.removeListener:
			for i, ch := range s.listeners {
				if ch == listenerToRemove {
					// replace item we want to remove with last item
					s.listeners[i] = s.listeners[len(s.listeners)-1]
					// remove last item
					s.listeners = s.listeners[:len(s.listeners)-1]
					close(ch)
					break
				}
			}
		case val, ok := <-s.source:
			if !ok {
				return
			}
			for _, listener := range s.listeners {
				if listener != nil {
					select {
					case listener <- val:
					case <-ctx.Done():
						return
					}

				}
			}
		}
	}
}

func NewBroadcastServer(ctx context.Context, source <-chan string) BroadcastServer {
	service := &broadcastServer{
		source:         source,
		listeners:      make([]chan string, 0),
		addListener:    make(chan chan string),
		removeListener: make(chan (<-chan string)),
	}
	go service.serve(ctx)
	return service
}
