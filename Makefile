MRUBY_COMMIT ?= 1.2.0
MRUBY_VENDOR_DIR ?= mruby-build

all: libmruby.a test

clean:
	rm -rf ${MRUBY_VENDOR_DIR}
	rm -f libmruby.a

gofmt:
	@echo "Checking code with gofmt.."
	gofmt -s *.go >/dev/null

lint:
	sh golint.sh

megacheck:
	go get honnef.co/go/tools/cmd/megacheck
	GO111MODULE=off megacheck ./...

libmruby.a: ${MRUBY_VENDOR_DIR}/mruby
	cd ${MRUBY_VENDOR_DIR}/mruby && ${MAKE}
	cp ${MRUBY_VENDOR_DIR}/mruby/build/host/lib/libmruby.a .

${MRUBY_VENDOR_DIR}/mruby:
	mkdir -p ${MRUBY_VENDOR_DIR}
	git clone https://github.com/mruby/mruby.git ${MRUBY_VENDOR_DIR}/mruby
	cd ${MRUBY_VENDOR_DIR}/mruby && git reset --hard && git clean -fdx
	cd ${MRUBY_VENDOR_DIR}/mruby && git checkout ${MRUBY_COMMIT}

test: gofmt lint
	go test -v

.PHONY: all clean libmruby.a test lint
