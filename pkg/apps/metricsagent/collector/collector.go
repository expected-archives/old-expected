package collector

import (
	"sync"
)

// packets is a list of list of packet.
// A packet is a []byte representing a metric.Metric
// A list of packet is all metrics collected on each
// running containers in this daemon.
var packets [][][]byte

var mutex = sync.RWMutex{}

func AddPackets(s [][]byte) {
	mutex.Lock()
	defer mutex.Unlock()
	packets = append(packets, s)
}

func GetPackets() [][][]byte {
	mutex.RLock()
	p := packets
	mutex.RUnlock()

	mutex.Lock()
	packets = make([][][]byte, 0)
	mutex.Unlock()
	return p
}
