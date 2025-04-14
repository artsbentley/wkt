# build `wkt` binary
build:
    @go build -o ./bin/wkt ./cmd/wkt/main.go

copy: 
	@go build -o ./bin/wkt ./cmd/wkt/main.go
	cp ./bin/wkt ~/.config/scripts/wkt 

# run 'wkt'
run:
	@just build
	@./bin/wkt


