package queue

import (
	"context"

	"github.com/Pineapple217/harbor-hawk/broadcast"
	"github.com/Pineapple217/harbor-hawk/docker"
)

var (
	activeBuildQueue *buildQueue
)

func GetBuildQueue() *buildQueue {
	return activeBuildQueue
}

func InitBuildQueue() {
	activeBuildQueue = newBuildQueue()
	go activeBuildQueue.Work()
	// TODO safe shutdown
}

type buildQueue struct {
	buildSettingsChannel chan docker.BuildSettings
	workingChannel       chan bool
	BuildLogsChannel     chan string
	Broadcaster          *broadcast.BroadcastServer
}

func newBuildQueue() *buildQueue {
	buildSettingsChannel := make(chan docker.BuildSettings, 100)
	workingChannel := make(chan bool, 1)
	buildLogsChannel := make(chan string, 1)
	broadcaster := broadcast.NewBroadcastServer(context.Background(), buildLogsChannel)

	return &buildQueue{
		buildSettingsChannel: buildSettingsChannel,
		workingChannel:       workingChannel,
		BuildLogsChannel:     buildLogsChannel,
		Broadcaster:          &broadcaster,
	}
}

func (e *buildQueue) Work() {
	for buildSettings := range e.buildSettingsChannel {
		e.workingChannel <- true
		// j, _ := json.Marshal(p)
		e.BuildLogsChannel <- "\x1B[1;3;31mSTART BUILD\x1B[0m\\r\\n"
		docker.BuildAndUploadImage(buildSettings, e.BuildLogsChannel)
		e.BuildLogsChannel <- "\\r\\n"
		e.BuildLogsChannel <- "\x1B[1;3;31mEND BUILD\x1B[0m\\r\\n"
		<-e.workingChannel
	}
}

func (e *buildQueue) Size() int {
	return len(e.buildSettingsChannel) + len(e.workingChannel)
}

func (e *buildQueue) Enqueue(buildSettings docker.BuildSettings) {
	e.buildSettingsChannel <- buildSettings
}
