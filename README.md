# Patternopedia Examples

Runnable Go code examples for patterns at [patternopedia.com](https://patternopedia.com).

## Table of Contents

- [Concurrency Examples](#concurrency-examples)

---

## Concurrency Examples

| # | Pattern | File | Run |
|---|---------|------|-----|
| 1.1 | [Goroutine](https://patternopedia.com/patterns/pt1-concurrency-patterns/goroutine/) | `concurrency/1.1-goroutine.go` | `go run concurrency/1.1-goroutine.go` |
| 2.1 | [Channel](https://patternopedia.com/patterns/pt1-concurrency-patterns/channel/) — Unbuffered | `concurrency/2.1-channel-unbuffered.go` | `go run concurrency/2.1-channel-unbuffered.go` |
| 2.2 | [Channel](https://patternopedia.com/patterns/pt1-concurrency-patterns/channel/) — Closure & Range | `concurrency/2.2-channel-closure-range.go` | `go run concurrency/2.2-channel-closure-range.go` |
| 3.1 | [Buffered Channel](https://patternopedia.com/patterns/pt1-concurrency-patterns/buffered-channel/) — Basics | `concurrency/3.1-buffered-channel-basics.go` | `go run concurrency/3.1-buffered-channel-basics.go` |
| 3.2 | [Buffered Channel](https://patternopedia.com/patterns/pt1-concurrency-patterns/buffered-channel/) — Close & Drain | `concurrency/3.2-buffered-channel-close-drain.go` | `go run concurrency/3.2-buffered-channel-close-drain.go` |
| 4.1 | [Done Channel](https://patternopedia.com/patterns/pt1-concurrency-patterns/done-channel/) — Basic | `concurrency/4.1-done-channel-basic.go` | `go run concurrency/4.1-done-channel-basic.go` |
| 4.2 | [Done Channel](https://patternopedia.com/patterns/pt1-concurrency-patterns/done-channel/) — Cleanup | `concurrency/4.2-done-channel-cleanup.go` | `go run concurrency/4.2-done-channel-cleanup.go` |

## Usage

```bash
git clone https://github.com/patternopedia/patternopedia-examples.git
cd patternopedia-examples
go run concurrency/1.1-goroutine.go
```

## License

Code examples may be used freely without attribution.
