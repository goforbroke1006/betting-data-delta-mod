test-run:
	go test ./...

test-cover:
	go test ./... -cover

test-bench:
	go test ./... -bench=.