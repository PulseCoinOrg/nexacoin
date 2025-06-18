build:
	go build -o ./gnex ./cmd

run:
	go run cmd/main.go 

clean:
	rm ./gnex && rm -rf ./chaindb-output