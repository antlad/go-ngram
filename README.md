go-ngram
========

Ngram index for Go.

## Key features

* Append only. Data can't be deleted from index.
* GC friendly (all strings are pooled and compressed)
* Application agnostic (there is no notion of document or something that user needs to implement)
 

## Usage

```go
index := ngram.NewNgramIndex(ngram.SetN(3))
tokenId, err := index.Add("hello") 
str, err := index.GetString(tokenId)  // str == "hello"
resultsList, err := index.Search("world")
```

## TODO:

* Smoothing functions (Laplace etc)

[![GoDoc](https://godoc.org/github.com/Lazin/go-ngram?status.png)](https://godoc.org/github.com/Lazin/go-ngram)
[![Coverage Status](https://coveralls.io/repos/Lazin/go-ngram/badge.png)](https://coveralls.io/r/Lazin/go-ngram)
