build:
	go build -buildmode=c-shared -o libGoReSym.so
clean:
	rm libGoReSym.h libGoReSym.so
