# 1BRC

one billion rows challenge [https://github.com/gunnarmorling/1brc](https://github.com/gunnarmorling/1brc)

## rules

- no external library used (is `golang.org/x/*` allowed?)

## generate file

```bash
go run cmd/create/main.go
```

## run baseline

```bash
go run cmd/std/main.go
```

## My fastest implementation

### what I did

- memory map zero copy & lock free reading
- custom split (zero copy)
- custom parsefloat
- concurrent workers
- xxhash fast hash algorithm

runs on i5-6400 ~20s

### run

```bash
go run main.go
```

## any idea make it faster?

issue are welcome ^_^
