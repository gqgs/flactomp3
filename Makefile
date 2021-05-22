
install: flactomp3

generate:
	go generate ./...

flactomp3: generate
	go install ./cmd/flactomp3/
