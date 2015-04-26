MRUBY_COMMIT ?= 1.1.0

all: libmruby.a
	go test

clean:
	rm -rf vendor
	rm -f libmruby.a

libmruby.a: vendor/mruby
	cd vendor/mruby && ${MAKE}
	cp vendor/mruby/build/host/lib/libmruby.a .

vendor/mruby:
	mkdir -p vendor
	git clone https://github.com/mruby/mruby.git vendor/mruby
	cd vendor/mruby && git reset --hard && git clean -fdx
	cd vendor/mruby && git checkout ${MRUBY_COMMIT}

.PHONY: all clean libmruby.a test
