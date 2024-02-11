package queue

import (
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
}

func newBuildQueue() *buildQueue {
	buildSettingsChannel := make(chan docker.BuildSettings, 100)
	workingChannel := make(chan bool, 1)
	buildLogsChannel := make(chan string, 100)
	return &buildQueue{
		buildSettingsChannel: buildSettingsChannel,
		workingChannel:       workingChannel,
		BuildLogsChannel:     buildLogsChannel,
	}
}

func (e *buildQueue) Work() {
	// for {
	// 	select {
	// 	case buildSettings := <-e.buildSettingsChannel:
	// 		// Enqueue message to workingChannel to avoid miscalculation in queue size.
	// 		e.workingChannel <- true

	// 		// Let's assume this time sleep is send email process
	// 		docker.BuildAndUploadImage(buildSettings, e.buildLogsChannel)

	// 		<-e.workingChannel
	// 	}
	// }
	for buildSettings := range e.buildSettingsChannel {
		e.workingChannel <- true
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
