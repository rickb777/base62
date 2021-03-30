# base62

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/base62)
[![Build Status](https://api.travis-ci.org/rickb777/base62.svg?branch=master)](https://travis-ci.org/rickb777/base62/builds)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/base62)](https://goreportcard.com/report/github.com/rickb777/base62)
[![Issues](https://img.shields.io/github/issues/rickb777/base62.svg)](https://github.com/rickb777/base62/issues)

base62 is a compact and fast implementation of `Base62` encoding/decoding algorithm,
which is inspired by the [java implementation by glowfall](https://github.com/glowfall/base62).

As this package was written, I did not find good base62 implementation on github.
The following packages can only encode/decode integers:

    - https://github.com/gravityblast/go-base62/
    - https://github.com/kare/base62
    - https://github.com/catinello/base62

While this package encodes/decodes bytes of arbitrary length, as well as integers.

This `Base62` implementation is much faster than `big.Int` based implementation,
and is not much slower than typical `Base64` implementations. See the following
benchmark results.

## Benchmark

```text
Benchmark_Encode-12             10000000               193 ns/op
Benchmark_Decode-12             20000000                68.4 ns/op

Benchmark_Encode_BigInt-12       1000000              1070 ns/op

Benchmark_Encode_Base64-12      20000000                59.7 ns/op
Benchmark_Decode_Base64-12      20000000                62.3 ns/op

Benchmark_Encode_Integer-12     30000000                44.4 ns/op
Benchmark_Decode_Integer-12    200000000                 9.63 ns/op
```

## Status

This is a maintenance fork from github.com/jxskiss/base62.
