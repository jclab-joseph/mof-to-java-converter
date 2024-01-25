PWD := $(shell pwd)

.PHONY: clean all
default: all

clean:
	cd omi/Unix && make clean || true
	rm -rf omi/Unix/output
	rm omi/Unix/GNUmakefile || true

omi/Unix/GNUmakefile:
	cd omi/Unix && ./configure --disable-rtti --disable-templates --disable-localsession --disable-indication --disable-shell --disable-encryption --disable-auth

omi/Unix/output/lib/libpal.a: omi/Unix/GNUmakefile
	cd omi/Unix && make -j4 -C pal

omi/Unix/output/lib/libmof.a: omi/Unix/GNUmakefile
	cd omi/Unix && make -j4 -C mof

converter.exe: omi/Unix/output/lib/libpal.a omi/Unix/output/lib/libmof.a
	CGO_CFLAGS="-I$(PWD)/omi/Unix/common" CGO_LDFLAGS="-L$(PWD)/omi/Unix/output/lib" go build -o converter.exe ./cmd/converter/

all: converter.exe
