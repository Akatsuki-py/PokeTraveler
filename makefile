ifdef COMSPEC
	EXE_EXT := .exe
else
	EXE_EXT := 
endif

.PHONY: run
run:
	go run ./cmd/main.go

build:
	go build -o poketraveler$(EXE_EXT) ./cmd/main.go