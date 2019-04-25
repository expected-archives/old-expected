package collector

import (
	"sync"
)

// packets is a list of list of packet.
// A packet is a []byte representing a stats.Stats
// A list of packet is all metrics collected on each
// running containers in this daemon.
var packets [][][]byte

var mutex = sync.Mutex{}

func AddPackets(s [][]byte) {
	mutex.Lock()
	defer mutex.Unlock()
	packets = append(packets, s)
}

func GetPackets() [][][]byte {
	mutex.Lock()
	defer mutex.Unlock()
	p := packets
	packets = make([][][]byte, 0)
	return p
}
