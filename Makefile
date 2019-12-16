MRUBY_COMMIT ?= 2.1.0

all: test

clean:
	rm -rf vendor
	rm -f libmruby.a

gofmt:
	@echo "Checking code with gofmt.."
	gofmt -s *.go >/dev/null

lint:
	sh golint.sh

megacheck:
	go get honnef.co/go/tools/cmd/megacheck
	GO111MODULE=off megacheck ./...

libmruby.a: vendor/mruby
	cd vendor/mruby && MRUBY_CONFIG=../../build_config.rb ${MAKE}
	cp vendor/mruby/build/host/lib/libmruby.a .

vendor/mruby:
	mkdir -p vendor
	git clone https://github.com/mruby/mruby.git vendor/mruby
	cd vendor/mruby && git reset --hard && git clean -fdx
	cd vendor/mruby && git checkout ${MRUBY_COMMIT}

test: libmruby.a gofmt lint
	go test -v

.PHONY: all clean libmruby.a test lint
