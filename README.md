samp: reservoir sampling on the command line
============================================

samp performs stream sampling on the command line. It's the lossy cat(1)
you've always wanted.

It emits items in the same order as they were found on input. It reads the
entire input before it can return anything, so ensure you feed it with a
stream that ends.

Install: `go get github.com/pteichman/samp`

```
Usage of samp:
  -0    use NUL characters instead of line delimiters
  -k int
        maximum number of items to pass (default 1)
  -s int
        random seed: -1 for current time (default -1)
```
