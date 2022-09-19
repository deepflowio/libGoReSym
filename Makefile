gcc-build:
	/usr/local/go/bin/go build -buildmode=c-shared -o libGoReSym.so
	/usr/local/go/bin/go build -buildmode=c-archive -o libGoReSym.a
gcc-test: gcc-build
	cd test && gcc main.c --static -L../ -lGoReSym  -lpthread && ./a.out

musl-build:
	CC=musl-gcc /usr/local/musl-go/bin/go build -buildmode=c-shared -o libGoReSym.so
	CC=musl-gcc /usr/local/musl-go/bin/go build -buildmode=c-archive -o libGoReSym.a

musl-test: musl-build
	cd test && musl-gcc main.c --static -L../ -lGoReSym  -lpthread && ./a.out

