package stats

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/google/uuid"
	"math"
	"strings"
)

type Stats struct {
	ID          uuid.UUID // 16
	Memory      uint32    // 32
	NetInput    float32   // 32
	NetOutput   float32   // 32
	BlockInput  uint32    // 32
	BlockOutput uint32    // 32
	Cpu         float32   // 32
	// 320 bits
	// 40 bytes
}

func (s *Stats) String() string {
	return fmt.Sprintf(""+
		"ID: %s\n"+
		"Memory: %d\n"+
		"Net: %.2f / %.2f\n"+
		"Block: %d / %d\n"+
		"Cpu: %.2f %%", s.ID, s.Memory, s.NetInput, s.NetOutput, s.BlockInput, s.BlockOutput, s.Cpu)
}

func (s *Stats) UnmarshalBinary(data []byte) error {
	if len(data) != 40 {
		return errors.New("unable to unmarshal stats")
	}

	seq := 16

	uid, err := uuid.FromBytes(data[:seq])
	if err != nil {
		return err
	}

	s.ID = uid

	// Memory
	s.Memory = binary.LittleEndian.Uint32(data[plus(&seq, 0) : seq+4])

	// Networks
	s.NetInput = math.Float32frombits(binary.LittleEndian.Uint32(data[plus(&seq, 4) : seq+4]))
	s.NetOutput = math.Float32frombits(binary.LittleEndian.Uint32(data[plus(&seq, 4) : seq+4]))

	// Blocks
	s.BlockInput = binary.LittleEndian.Uint32(data[plus(&seq, 0) : seq+4])
	s.BlockOutput = binary.LittleEndian.Uint32(data[plus(&seq, 0) : seq+4])

	// Cpu
	s.Cpu = math.Float32frombits(binary.LittleEndian.Uint32(data[plus(&seq, 4) : seq+4]))

	return nil
}

func (s Stats) MarshalBinary() (data []byte, err error) {
	var b []byte
	buffer := bytes.NewBuffer(b)

	//uuid
	buid, err := s.ID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buffer.Write(buid)

	// Memory
	bits32toBuffer(buffer, s.Memory)

	// Networks
	bits32toBuffer(buffer, math.Float32bits(s.NetInput))
	bits32toBuffer(buffer, math.Float32bits(s.NetOutput))

	// Blocks
	bits32toBuffer(buffer, s.BlockInput)
	bits32toBuffer(buffer, s.BlockOutput)

	// Cpu
	bits32toBuffer(buffer, math.Float32bits(s.Cpu))

	return buffer.Bytes(), nil
}

func bits32toBuffer(buffer *bytes.Buffer, u32 uint32) {
	buffer.Write([]byte{
		byte(u32),
		byte(u32 >> 8),
		byte(u32 >> 16),
		byte(u32 >> 24),
	})
}

func plus(n *int, add int) int {
	*n += add
	return *n
}

func FromDockerStats(t types.StatsJSON, uuid uuid.UUID) Stats {
	s := Stats{
		ID:     uuid,
		Memory: uint32(t.MemoryStats.Usage),
		Cpu:    calculateCPUPercentUnix(t.PreCPUStats.CPUUsage.TotalUsage, t.PreCPUStats.SystemUsage, t),
	}
	s.BlockInput, s.BlockOutput = calculateBlockIO(t.BlkioStats)
	s.NetInput, s.NetOutput = calculateNetwork(t.Networks)

	return s
}

func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v types.StatsJSON) float32 {
	var (
		cpuPercent = float32(0.0)
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float32(v.CPUStats.CPUUsage.TotalUsage) - float32(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float32(v.CPUStats.SystemUsage) - float32(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float32(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}

func calculateBlockIO(blkio types.BlkioStats) (blkRead uint32, blkWrite uint32) {
	for _, bioEntry := range blkio.IoServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			blkRead = uint32(blkRead + uint32(bioEntry.Value))
		case "write":
			blkWrite = blkWrite + uint32(bioEntry.Value)
		}
	}
	return
}

func calculateNetwork(network map[string]types.NetworkStats) (float32, float32) {
	var rx, tx float32

	for _, v := range network {
		rx += float32(v.RxBytes)
		tx += float32(v.TxBytes)
	}
	return rx, tx
}
