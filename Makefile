local:
	goreleaser release --snapshot --clean
test-stdin:
	go test ./cmd -run TestStdin  -count=1 -v
test-fp:
	go test ./cmd -run TestFp  -count=1 -v