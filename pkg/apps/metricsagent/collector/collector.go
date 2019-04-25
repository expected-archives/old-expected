package collector

import (
	"sync"
)

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
