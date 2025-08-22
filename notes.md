# WTF-GO // GOtchats
## A run down of some strange quirks and intricacies of the Go language. 
- Based loosely off `100 Go Mistakes And How to Avoid Them`

### Memory Leaks
#### Initialisation:
- #NOTE: maybe put something more interactive up top?
- Initialising both maps and slices requires careful attention.
- If we know the potential length of an array or map beforehand we should initialise

```go
func getBar(foos []Foo) []Bar {
  bars := make([]Bar, 0) // Go runtime forces the length here

  for _, foo := range foos {
    bars = append(bars, fooToBar(foo))
  }

  return bars
}
```
- a slice grows by doubling its size until it contains 1,024 elements, after which it grows by 25%.
  - therefore from 0->1024 we get 11 new backing arrays in 

```bash
BenchmarkEmptySlice-8                577           2075027 ns/op
BenchmarkCapSlice-8                 3793            312279 ns/op
BenchmarkLenSlice-8                 3792            312003 ns/op
```

#### Deletion:

- 
- The Go Map implementation is a pointer to the folloowing struct:

  ```go
type hmap struct {
  // Note: the format of the hmap is also encoded in cmd/compile/internal/reflectdata/reflect.go.
  // Make sure this stays in sync with the compiler's definition.
  count     int // # live cells == size of map.  Must be first (used by len() builtin)
  flags     uint8
  B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
  noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
  hash0     uint32 // hash seed

  buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
  oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
  nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

  extra *mapextra // optional fields
  ```
#### Slices:
- 

