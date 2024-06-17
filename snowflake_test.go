package snowflake

import (
	"math/rand"
	"testing"
	"time"
)

func Test_SnowflakeInvalidParameterRange(t *testing.T) {
	source := rand.NewSource(time.Now().UnixMilli())
	r := rand.New(source)
	_, err := NewSource(uint64(r.Intn(100)+32), 0, 0)
	if err == nil {
		t.Error("should have failed")
	}

	_, err = NewSource(0, uint64(r.Intn(100)+32), 0)
	if err == nil {
		t.Error("should have failed")
	}

	_, err = NewSource(0, 0, uint64(r.Intn(100)+4096))
	if err == nil {
		t.Error("should have failed")
	}
}

func Test_SnowflakeValidParameterRange(t *testing.T) {
	source := rand.NewSource(time.Now().UnixMilli())
	r := rand.New(source)
	_, err := NewSource(uint64(r.Intn(32)), uint64(r.Intn(32)), uint64(r.Intn(4096)))
	if err != nil {
		t.Error("should'n fail")
	}
}

func Test_SnowflakeIdPartitions(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	dataCenterId := uint64(r.Intn(32))
	machineId := uint64(r.Intn(32))
	sequenceNo := uint64(r.Intn(4096))
	source, err := NewSource(dataCenterId, machineId, sequenceNo)
	id := source.Generate()
	if err != nil {
		t.Error(err)
	}

	if val := id.GetDataCenterId(); val != dataCenterId {
		t.Error("data center id doesn't matched", val, dataCenterId)
	}

	if id.GetMachineId() != machineId {
		t.Error("machine id doesn't matched")
	}

	if id.GetSequenceNo() != sequenceNo {
		t.Error("sequenceNo id doesn't matched")
	}
}

// Benchmarking
func Benchmark_GenerateId(b *testing.B) {
	id, _ := NewSource(31, 31, 4095)
	for i := 0; i < b.N; i++ {
		id.Generate()
	}
}

func Benchmark_GetDataCenterId(b *testing.B) {
	source, _ := NewSource(31, 31, 4095)
	id := source.Generate()
	for i := 0; i < b.N; i++ {
		id.GetDataCenterId()
	}
}

func Benchmark_GetMachineId(b *testing.B) {
	source, _ := NewSource(31, 31, 4095)
	id := source.Generate()
	for i := 0; i < b.N; i++ {
		id.GetMachineId()
	}
}

func Benchmark_GetSequenceNo(b *testing.B) {
	source, _ := NewSource(31, 31, 4095)
	id := source.Generate()
	for i := 0; i < b.N; i++ {
		id.GetSequenceNo()
	}
}
