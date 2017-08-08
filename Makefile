BIN_DIR=bin

default: zippy

zippy: clean
	mkdir -p $(BIN_DIR)
	godep go build -o $(BIN_DIR)/zippy apps/*.go
	cp bin/zippy ${GOPATH}/bin/zippy

clean:
	rm -rf $(BIN_DIR)/*

test: start_zippy
	cd apps && godep go test
	@cd ..
	pkill -9 zippy || true

start_zippy: stop_zippy zippy
	pgrep zippy > /dev/null || (./bin/zippy > zippy.log 2>&1 &)

stop_zippy:
	pkill -9 zippy || true