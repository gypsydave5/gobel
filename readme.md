# Gobel

A _WIP_ implementation of Paul Graham's [Bel language][bel] written in Go.

~Right now it's a calculator that can do conditionals.~

Now featuring actual lambdas.

## The bad REPL

```shell
$ go build -o repl cmd/repl/main.go
$ ./repl
```

## Run the tests

```shell
$ go test ./...
```

## Influence

- The [Bel language][bel], obviously.
- Peter Norvig's [_(How to Write a (Lisp) Interpreter (in Python))_](https://norvig.com/lispy.html)

[bel]: (http://paulgraham.com/bel.html)
