# build `wkt` binary
build:
    @go build -o ./bin/wkt ./cmd/wkt/main.go

# run 'wkt'
run:
	@just build
	@./bin/wkt


