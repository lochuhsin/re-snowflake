.PHONY: test
test: 
	go test ./...

.PHONY: build
build: 
	go build ./...

.PHONY: build-race
build-race: 
	go build -race ./...

.PHONY: profile
profile:
	go test ./... -run=none -bench=. -benchmem -benchtime=2s -memprofile=mem.pprof -cpuprofile=cpu.pprof -blockprofile=block.pprof
 