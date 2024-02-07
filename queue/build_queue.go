package queue

import (
	"github.com/Pineapple217/harbor-hawk/docker"
)

type buildQueue struct {
	buildSettingsChannel chan BuildSettings
	workingChannel       chan bool
}

func NewBuild() *buildQueue {
	buildSettingsChannel := make(chan BuildSettings, 100)
	workingChannel := make(chan bool, 100)
	return &buildQueue{
		buildSettingsChannel: buildSettingsChannel,
		workingChannel:       workingChannel,
	}
}

// Logical flow from the queue
func (e *buildQueue) Work() {
	for {
		select {
		case <-e.buildSettingsChannel:
			// Enqueue message to workingChannel to avoid miscalculation in queue size.
			e.workingChannel <- true

			// Let's assume this time sleep is send email process
			docker.BuildAndUploadImage()

			<-e.workingChannel
		}
	}
}

// Size is a function to get the size of email queue
func (e *buildQueue) Size() int {
	return len(e.buildSettingsChannel) + len(e.workingChannel)
}

// Enqueue is a function to enqueue email string into email channel
func (e *buildQueue) Enqueue(emailString string) {
	e.buildSettingsChannel <- emailString
}
