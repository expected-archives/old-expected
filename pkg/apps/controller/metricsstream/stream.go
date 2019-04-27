package metricsstream

import (
	"github.com/expectedsh/expected/pkg/protocol"
	"sync"
)

var streams = make(map[string]map[string]protocol.Controller_GetContainerMetricsServer)
var mutex sync.Mutex

func AddStream(id, streamId string, stream protocol.Controller_GetContainerMetricsServer) {
	mutex.Lock()
	defer mutex.Unlock()
	str, ok := streams[id]
	if ok {
		str[streamId] = stream
	} else {
		str = map[string]protocol.Controller_GetContainerMetricsServer{streamId: stream}
	}
	streams[id] = str
}

func RemoveStream(id, streamId string) {
	mutex.Lock()
	defer mutex.Unlock()

	str, ok := streams[id]
	if ok {
		delete(str, streamId)
	}
}

func Send(id string, metric []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	str, ok := streams[id]
	if ok {
		reply := protocol.GetContainerMetricsReply{Message: metric}
		for _, v := range str {
			_ = v.Send(&reply)
		}
	}
}
