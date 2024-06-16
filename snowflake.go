package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	SEQUENCE_NO_LENGTH    = 0xC
	MACHINE_ID_LENGTH     = 0x5
	DATA_CENTER_ID_LENGTH = 0x5
	TIMESTAMP_LENGTH      = 41
	PRESERVED_OFFSET      = 1
	MACHINE_ID_MASK       = (1<<MACHINE_ID_LENGTH - 1)
	DATA_CENTER_ID_MASK   = (1<<DATA_CENTER_ID_LENGTH - 1)
	SEQUENCE_ID_MASK      = (1<<SEQUENCE_NO_LENGTH - 1)
)

type SnowflakeId struct {
	bitset    uint64
	preBitset uint64
	bitShift  uint64
	mu        sync.Mutex
}

func NewSnowflakeId(centerId, machineId, sequenceNo uint64) (SnowflakeId, error) {
	/**
	 * Total 64 bit
	 * | 0 | 41 bit timestamp | 5 bit data center id | 5 bit machine id | 12 bit sequence id |
	 * We build the first 22 bits of known values, the 41 bit timestamp will leave it to generate
	 */
	if centerId >= (1 << DATA_CENTER_ID_LENGTH) {
		return *new(SnowflakeId), errors.New("exceed data center id limit, 31")
	}

	if machineId >= (1<<MACHINE_ID_LENGTH + 1) {
		return *new(SnowflakeId), errors.New("exceed machine id limit, 31")
	}

	if sequenceNo >= (1<<SEQUENCE_NO_LENGTH + 1) {
		return *new(SnowflakeId), errors.New("exceed sequence id limit, 4095")
	}
	return SnowflakeId{
		bitset:    0,
		preBitset: BuildPreBitMask(centerId, machineId, sequenceNo),
		bitShift:  DATA_CENTER_ID_LENGTH + MACHINE_ID_LENGTH + SEQUENCE_NO_LENGTH,
	}, nil
}

func (s *SnowflakeId) Generate() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bitset = (uint64(time.Now().UnixMilli()) << s.bitShift) | s.preBitset
}

func (s *SnowflakeId) GetTime() uint64 {
	return s.bitset >> s.bitShift
}

func (s *SnowflakeId) GetDataCenterId() uint64 {
	shift := s.bitShift - DATA_CENTER_ID_LENGTH
	return (s.bitset >> shift) & DATA_CENTER_ID_MASK
}
func (s *SnowflakeId) GetMachineId() uint64 {
	return (s.bitset >> SEQUENCE_NO_LENGTH) & MACHINE_ID_MASK
}

func (s *SnowflakeId) GetSequenceId() uint64 {
	return s.bitset & SEQUENCE_ID_MASK
}

func (s *SnowflakeId) GetId() uint64 {
	return s.bitset
}

func BuildPreBitMask(centerId, machineId, sequenceNo uint64) uint64 {
	cMask := centerId << (MACHINE_ID_LENGTH + SEQUENCE_NO_LENGTH)
	mMask := machineId << SEQUENCE_NO_LENGTH

	return (cMask | mMask) | sequenceNo
}
