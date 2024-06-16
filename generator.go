package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	SEQUENCE_NO_LENGTH    = 12
	MACHINE_ID_LENGTH     = 5
	DATA_CENTER_ID_LENGTH = 5
	TIMESTAMP_LENGTH      = 41
	SNOWFLAKE_ID_LENGTH   = 64
	PRESERVED_OFFSET      = 1
)

type Bitset struct {
	bits   []uint
	length int
	offset int
}

func NewBitset(length int) Bitset {
	return Bitset{
		bits:   make([]uint, length),
		length: length,
		offset: 1,
	}
}

func (b *Bitset) Set(bits []uint) (n int, err error) {
	if b.offset+len(bits) > b.length {
		return -1, errors.New("byte range out of bounds")
	}
	n = copy(b.bits[b.offset:], bits)
	b.offset += len(bits)
	return n, nil
}

func (b *Bitset) Getbitset() []uint {
	return b.bits
}

type SnowflakeId struct {
	Bitset
	datacenter int
	machine    int
	sequenceNo int
	mu         sync.Mutex
}

func NewSnowflakeId(datacenter, machine, sequenceNo int) SnowflakeId {
	return SnowflakeId{
		Bitset:     NewBitset(SNOWFLAKE_ID_LENGTH),
		datacenter: datacenter,
		machine:    machine,
		sequenceNo: sequenceNo,
	}
}
func (s *SnowflakeId) Generate() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	timestamp := time.Now().UnixMilli()

	bits, err := itobitset(int(timestamp))
	if err != nil {
		return err
	}
	s.Set(bits)

	bits, err = itobitset(s.datacenter)
	if err != nil {
		return err
	}
	if len(bits) > DATA_CENTER_ID_LENGTH {
		return errors.New("the value of data center should not exceed 31")
	}
	if len(bits) != DATA_CENTER_ID_LENGTH {
		nbits := make([]uint, DATA_CENTER_ID_LENGTH)
		copy(nbits[DATA_CENTER_ID_LENGTH-len(bits):], bits)
		bits = nbits
	}
	s.Set(bits)

	bits, err = itobitset(s.machine)
	if err != nil {
		return err
	}
	if len(bits) > MACHINE_ID_LENGTH {
		return errors.New("the value of data center should not exceed 31")
	}
	if len(bits) != MACHINE_ID_LENGTH {
		nbits := make([]uint, MACHINE_ID_LENGTH)
		copy(nbits[MACHINE_ID_LENGTH-len(bits):], bits)
		bits = nbits
	}
	s.Set(bits)

	bits, err = itobitset(s.sequenceNo)
	if err != nil {
		return err
	}
	if len(bits) > SEQUENCE_NO_LENGTH {
		return errors.New("the value of data center should not exceed 31")
	}
	if len(bits) != SEQUENCE_NO_LENGTH {
		nbits := make([]uint, SEQUENCE_NO_LENGTH)
		copy(nbits[SEQUENCE_NO_LENGTH-len(bits):], bits)
		bits = nbits
	}
	s.Set(bits)
	return nil
}

func (s *SnowflakeId) Time() uint {
	return bitsettoi(s.bits[PRESERVED_OFFSET : TIMESTAMP_LENGTH+PRESERVED_OFFSET])
}

func (s *SnowflakeId) DataCenter() uint {
	return bitsettoi(s.bits[PRESERVED_OFFSET+TIMESTAMP_LENGTH : PRESERVED_OFFSET+TIMESTAMP_LENGTH+DATA_CENTER_ID_LENGTH])
}

func (s *SnowflakeId) Machine() uint {
	return bitsettoi(s.bits[PRESERVED_OFFSET+TIMESTAMP_LENGTH+DATA_CENTER_ID_LENGTH : PRESERVED_OFFSET+TIMESTAMP_LENGTH+DATA_CENTER_ID_LENGTH+MACHINE_ID_LENGTH])
}

func (s *SnowflakeId) Sequence() uint {
	return bitsettoi(s.bits[PRESERVED_OFFSET+TIMESTAMP_LENGTH+DATA_CENTER_ID_LENGTH+MACHINE_ID_LENGTH:])
}

func (s *SnowflakeId) Int64() uint64 {
	return uint64(bitsettoi(s.bits))
}

func bitsettoi(bits []uint) uint {
	var val uint = 0
	for _, v := range bits {
		val = val<<1 + v
	}
	return val
}

func itobitset(val int) (bits []uint, err error) {
	if val < 0 {
		return nil, errors.New("invalid number")
	}
	bits = []uint{}
	for val > 0 {
		bit := uint(val) % 2
		bits = append([]uint{bit}, bits...)
		val /= 2
	}
	return bits, nil
}
