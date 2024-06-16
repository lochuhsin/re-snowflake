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
	MACHINE_ID_MASK       = (1<<MACHINE_ID_LENGTH - 1)
	DATA_CENTER_ID_MASK   = (1<<DATA_CENTER_ID_LENGTH - 1)
	SEQUENCE_ID_MASK      = (1<<SEQUENCE_NO_LENGTH - 1)
	BIT_SHIFT             = DATA_CENTER_ID_LENGTH + MACHINE_ID_LENGTH + SEQUENCE_NO_LENGTH
)

type Id struct {
	bitset uint64
}

type Source struct {
	preBitset uint64
	mu        sync.Mutex
}

func NewSource(centerId, machineId, sequenceNo uint64) (Source, error) {
	/**
	 * Total 64 bit
	 * | 0 | 41 bit timestamp | 5 bit data center id | 5 bit machine id | 12 bit sequence id |
	 * We build the first 22 bits of known values, the 41 bit timestamp will leave it to generate
	 */
	if centerId >= (1 << DATA_CENTER_ID_LENGTH) {
		return *new(Source), errors.New("exceed data center id limit, 31")
	}

	if machineId >= (1<<MACHINE_ID_LENGTH + 1) {
		return *new(Source), errors.New("exceed machine id limit, 31")
	}

	if sequenceNo >= (1<<SEQUENCE_NO_LENGTH + 1) {
		return *new(Source), errors.New("exceed sequence id limit, 4095")
	}
	return Source{
		preBitset: BuildPreBitMask(centerId, machineId, sequenceNo),
	}, nil
}

func (s *Source) Generate() Id {
	s.mu.Lock()
	defer s.mu.Unlock()
	return Id{
		bitset: (uint64(time.Now().UnixMilli()) << BIT_SHIFT) | s.preBitset,
	}
}

func (i *Id) GetTime() uint64 {
	return i.bitset >> BIT_SHIFT
}

func (i *Id) GetDataCenterId() uint64 {
	shift := BIT_SHIFT - DATA_CENTER_ID_LENGTH
	return (i.bitset >> shift) & DATA_CENTER_ID_MASK
}
func (i *Id) GetMachineId() uint64 {
	return (i.bitset >> SEQUENCE_NO_LENGTH) & MACHINE_ID_MASK
}

func (i *Id) GetSequenceNo() uint64 {
	return i.bitset & SEQUENCE_ID_MASK
}

func (i *Id) GetId() uint64 {
	return i.bitset
}

func BuildPreBitMask(centerId, machineId, sequenceNo uint64) uint64 {
	cMask := centerId << (MACHINE_ID_LENGTH + SEQUENCE_NO_LENGTH)
	mMask := machineId << SEQUENCE_NO_LENGTH

	return (cMask | mMask) | sequenceNo
}
