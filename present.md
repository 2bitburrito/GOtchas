---
theme: 
  name: "catppuccin-mocha"
---
<!-- pause -->
A Few Small
---
<!-- pause -->
```python +no_background

   █████████     ███████     █████             █████                       
  ███░░░░░███  ███░░░░░███  ░░███             ░░███                        
 ███     ░░░  ███     ░░███ ███████    ██████  ░███████    ██████    █████ 
░███         ░███      ░███░░░███░    ███░░███ ░███░░███  ░░░░░███  ███░░  
░███    █████░███      ░███  ░███    ░███ ░░░  ░███ ░███   ███████ ░░█████ 
░░███  ░░███ ░░███     ███   ░███ ███░███  ███ ░███ ░███  ███░░███  ░░░░███
 ░░█████████  ░░░███████░    ░░█████ ░░██████  ████ █████░░████████ ██████ 
  ░░░░░░░░░     ░░░░░░░       ░░░░░   ░░░░░░  ░░░░ ░░░░░  ░░░░░░░░ ░░░░░░  
```
<!-- pause -->
<!-- alignment: center -->
I learned from reading this book:
![](images/cover2.jpg)
<!-- speaker_note: Intro self -->
<!-- speaker_note: Teiva Harsanyi -->
<!-- end_slide -->
![image:w:55%](images/meme3.jpg)
<!-- speaker_note: easy to learn / hard to master -->
<!-- end_slide -->

<!-- alignment: center -->
 Memory and CPU Performance of Slices and Maps
---
<!-- alignment: left -->
## Initialisation Mistakes

---
<!-- alignment: center -->
<!-- column_layout: [1, 1] -->

<!-- pause -->
<!-- column: 0 -->
```go {all|2|4-6|8|all} +line_numbers
func getBar(foos []Foo) []Bar {
  bars := make([]Bar, 0) 

  for _, foo := range foos {
    bars = append(bars, fooToBar(foo))
  }

  return bars
}
```

<!-- pause -->
<!-- column: 1 -->
![](images/go_q.jpg)
<!-- reset_layout -->
<!-- pause -->
A slice grows by doubling its backing array until it contains 1,024 elements, after which it grows by 25%.
<!-- pause -->
So from 0->1024 we get 11 new backing arrays
<!-- speaker_note: when we exceed the backing array len a new arr is made -->
<!-- speaker_note: Impacts both performance and memory -->
<!-- end_slide -->

<!-- alignment: center -->
 Memory and CPU Performance of Slices and Maps
---


<!-- alignment: center -->
So what can we do about this?

---

<!-- column_layout: [1, 1] -->

<!-- speaker_note: solutions involve knowing len of new slice-->
<!-- speaker_note: go test ./1-slice-init -bench=. -->
<!-- column: 0 -->
<!-- pause -->
![](images/array1.png)
```go {all|2-3|all} +line_numbers
func convertGivenCapacity(foos []Foo) []Bar {
	n := len(foos)
	bars := make([]Bar, 0, n)

	for _, foo := range foos {
		bars = append(bars, fooToBar(foo))
	}
	return bars
}
```
<!-- pause -->
<!-- column: 1 -->
![](images/array2.png)
```go {all|5-7|all} +line_numbers
func convertGivenLength(foos []Foo) []Bar {
  n := len(foos)
  bars := make([]Bar, n)

  for i, foo := range foos {
    bars[i] = fooToBar(foo)
  }
  return bars
}
```
<!-- pause -->
<!-- 
speaker_note: |
  - Mention Maps
-->

<!-- reset_layout -->
Same with Maps:
```go
barMap := make(map[string]Bar, len(foos))
```
<!-- end_slide -->

<!-- alignment: center -->
 Memory and CPU Performance of Slices and Maps
---
<!-- alignment: left -->
## Deletion Mistakes
---
<!-- column_layout: [1, 1] -->
<!-- column: 0 -->
```go {all|2-3|5-7|9-13|all} +line_numbers 
func main() {
	n := 1_000_000
	m := make(map[int]Foo, n)

	for i := range n { // Adds 1 million elements
		m[i] = Foo{}
	}

	for i, foo := range m { // Deletes if foo is expired
    if foo.isExpired() {
			delete(m, i)
		}
	}
}
```
<!-- pause -->
<!-- column: 1 -->
![](images/go_q.jpg)
<!-- end_slide -->
<!-- alignment: center -->
 Memory and CPU Performance of Slices and Maps
---
<!-- alignment: left -->
## Deletion Mistakes
---
![](images/map-buckets1.png)
<!-- end_slide -->
<!-- alignment: center -->
 Memory and CPU Performance of Slices and Maps
---
<!-- alignment: left -->
## Deletion Mistakes
---
![](images/map-buckets2.png)
<!-- end_slide -->
<!-- alignment: center -->
 Memory and CPU Performance of Slices and Maps
---
<!-- alignment: left -->
## Deletion Mistakes
---
```go {all|2-4|6-9|11-13|15-17|all} +line_numbers 
func main() {
	n := 1_000_000
	m := make(map[int][128]byte)
	printAlloc()

	for i := range n { // Adds 1 million elements
		m[i] = [128]byte{}
	}
	printAlloc()

	for i := range n { // Deletes 1 million elements
		delete(m, i)
	}

	runtime.GC() // Triggers a manual GC
	printAlloc()
	runtime.KeepAlive(m) // Keeps a reference to m so that the map isn’t collected
}
```
<!-- pause -->
<!-- alignment: center -->
After adding 1 million elements, there are 262,144 buckets.
<!-- speaker_note: go run ./2-map-deletes -->

<!-- end_slide -->
<!-- alignment: center -->
 Memory and CPU Performance of Slices and Maps
---
What can we do?

---
```go {all|3|all}+line_numbers
func main() {
	n := 1_000_000
	m := make(map[int]*[128]byte)
	printAlloc()

	for i := range n { // Adds 1 million elements
		m[i] = &[128]byte{}
	}
	printAlloc()

	for i := range n { // Deletes 1 million elements
		delete(m, i)
	}

	runtime.GC() // Triggers a manual GC
	printAlloc()
	runtime.KeepAlive(m) // Keeps a reference to m so that the map isn’t collected
}
```
<!-- speaker_note: go run ./2-map-deletes-pointers -->
<!-- speaker_note: also could copy map however this could balloon mem -->
<!-- end_slide -->

<!-- alignment: center -->
Go Routines
---
<!-- alignment: left -->
## The cost of concurrency

---
<!-- column_layout: [1, 1] -->
<!-- column: 0 -->
```go {all|6-9|all} +line_numbers 
func sequentialMergesort(s []int) {
	if len(s) <= 1 {
		return
	}

	middle := len(s) / 2
	sequentialMergesort(s[:middle])
	sequentialMergesort(s[middle:])
	merge(s, middle)
}
```
<!-- column: 1 -->
```go {all|7-18|all} +line_numbers 
func parallelMergesortV1(s []int) {
	if len(s) <= 1 {
		return
	}
	middle := len(s) / 2

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		parallelMergesortV1(s[:middle])
	}()

	go func() {
		defer wg.Done()
		parallelMergesortV1(s[middle:])
	}()

	wg.Wait()
	merge(s, middle)
}
```
<!-- speaker_note: for 10_000 what is faster? -->
<!-- speaker_note: go test ./3-concurrency -bench="Benchmark_mergesort" -->
<!-- pause -->

<!-- column: 0 -->
![image:w:50%](images/go_q.jpg)
<!-- end_slide -->
<!-- alignment: center -->
Premature Go Routines
---
<!-- alignment: left -->
## The cost of concurrency

---
<!-- column_layout: [1, 1] -->
<!-- column: 1 -->
![image:w:80%](images/routine_meme.jpg)
<!-- column: 0 -->
- While extremely lightweight
<!-- speaker_note: around 2kb per "thread" -->
- Go Routines aren't free
<!-- speaker_note: switching context is still expensive -->
<!-- end_slide -->
<!-- alignment: center -->
Premature Go Routines
---
<!-- alignment: left -->
---

<!-- column_layout: [1, 1] -->
<!-- column: 0 -->
So is parallelism useless here?

<!-- column: 1 -->
<!-- pause -->
```go {all|1|8-9|10-28|all} +line_numbers
const max = 2048 // Defines the threshold

func parallelMergesortV2(s []int) {
    if len(s) <= 1 {
        return
    }

    if len(s) <= max {
        sequentialMergesort(s) 
    } else {
        middle := len(s) / 2

        var wg sync.WaitGroup
        wg.Add(2)

        go func() {
            defer wg.Done()
            parallelMergesortV2(s[:middle])
        }()

        go func() {
            defer wg.Done()
            parallelMergesortV2(s[middle:])
        }()

        wg.Wait()
        merge(s, middle)
    }
}
```
<!-- speaker_note: go test ./3-concurrency -bench=. -->
<!-- end_slide -->
<!-- alignment: center -->
CPU L1 Cache Optimizations
---
<!-- alignment: left -->
## In the weeds now...
---

<!-- column_layout: [ 1, 1 ] -->

<!-- column: 0 -->
```go +line_numbers
type Input struct {
	a int64
	b int64
}

type Result struct {
	sumA int64
	sumB int64
}
```
<!-- pause -->
<!-- column: 1 -->
```go {all|7-19|all} +line_numbers
func count1(inputs []Input) Result1 {
	wg := sync.WaitGroup{}
	wg.Add(2)

	result := Result1{}

	go func() {
		for i := range len(inputs) {
			result.sumA += inputs[i].a
		}
		wg.Done()
	}()

	go func() {
		for i := range len(inputs) {
			result.sumB += inputs[i].b
		}
		wg.Done()
	}()

	wg.Wait()
	return result
}
```
<!-- pause -->
<!-- column: 0 -->
![](images/go_q.jpg)
<!-- end_slide -->
<!-- alignment: center -->
CPU L1 Cache Optimizations
---
So what's happening here?

![image:w:100%](images/cacheline_1.png)
<!-- end_slide -->
<!-- alignment: center -->
CPU L1 Cache Optimizations
---
So what's happening here?

![image:w:100%](images/cacheline_2.png)
<!-- speaker_note: 7/8 times -->
<!-- end_slide -->
<!-- alignment: center -->
CPU L1 Cache Optimizations
---
What can we do? 
<!-- pause -->
```go {all|3|all} +line_numbers
type Result2 struct {
	sumA int64
	_    [56]byte
	sumB int64
}
```
<!-- speaker_note: go test ./4-false-sharing -bench=. -->
<!-- end_slide -->
<!-- jump_to_middle -->

Thank You
---
<!-- speaker_note: read the book -->
<!-- speaker_note: jobs? -->
