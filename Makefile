build:
	go build -buildmode=c-shared -o libGoReSym.so
	go build -buildmode=c-archive -o libGoReSym.a

clean:
	rm libGoReSym.h libGoReSym.so libGoReSym.a

test: build
	cd test && gcc main.c --static -L../ -lGoReSym  -lpthread && ./a.out

