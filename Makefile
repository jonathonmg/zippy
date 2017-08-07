BIN_DIR=bin


default: all

all: zippy

zippy: clean
	mkdir -p $(BIN_DIR)
	godep go build -o $(BIN_DIR)/zippy apps/*.go
	cp bin/zippy ${GOPATH}/bin/zippy

clean:
	rm -rf $(BIN_DIR)/*