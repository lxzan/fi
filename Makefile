test:
	go test ./...

cover:
	go test -coverprofile=./bin/cover.out --cover ./...

bench:
	go test -benchmem -run=^$$ -bench ^Benchmark github.com/lxzan/fi