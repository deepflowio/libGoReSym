# libGoReSym

C wrapper for [GoReSym](https://github.com/mandiant/GoReSym)

## BUILD

### GCC

```bash
make gcc-build
make install
make clean
```

### musl

Your Go needs this [Patch](https://go-review.googlesource.com/c/go/+/334991/).

Go relies on the feature of glibc, which does not belong to the C language standard. And musl does not implement this feature, so it is necessary to patch the go language to fix this error. This patch itself also has some problems, which will cause some go functions to be abnormal, but Does not affect this library.

```bash
make musl-build
make install
make clean
```
