# deepCopy
### Description
This package provides functionality for making copies of some specific objects like pointers, maps channels and etc by value.

Example:
```go
package main

import (
	"github.com/Sereger/deepCopy"
	"log"
)

type (
	MyPointer struct {
		A, B string
	}

	MyTestStruct struct {
		Pointer *MyPointer
	}
)

func main() {
	orig := &MyTestStruct{
		Pointer: &MyPointer{
			A: "AAAA",
			B: "BBBB",
		},
	}

	clone := new(MyTestStruct)
	err := deepcopy.Copy(orig, clone)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\norig: %+v\ncopy: %+v\n", *orig, *clone)
	log.Printf("\norig: %+v\ncopy: %+v", *orig.Pointer, *clone.Pointer)
}

```

output:
```bash
orig: {Pointer:0xc00008e000}
copy: {Pointer:0xc00008e020}

orig: {A:AAAA B:BBBB}
copy: {A:AAAA B:BBBB}
```

### Performance
You can use `json.Encode / json.Decode` (or some else) for attain same results, but it's not effectively.
Example:
```bash
make bench

BenchmarkReflectCopy/reflect-simple-4             158290              7068 ns/op            2952 B/op         76 allocs/op
BenchmarkReflectCopy/reflect-recursive-4          140791              8427 ns/op            3906 B/op         79 allocs/op
BenchmarkReflectCopy/jsonEnc-Dec-4                 55137             21934 ns/op            7529 B/op        131 allocs/op
``` 

### Others features
#### Copy recursive pointers
With this package you can copy structures with cyclic dependencies using option `CheckRecursive`.
Example:
```go
package main

import (
	"github.com/Sereger/deepCopy"
	"log"
)

type (
	MyPointerRecursive struct {
		A, B       string
		RecPointer *MyTestStructRecursive
	}

	MyTestStructRecursive struct {
		Pointer *MyPointerRecursive
	}
)

func main() {
	orig := &MyTestStructRecursive{
		Pointer: &MyPointerRecursive{
			A: "AAAA",
			B: "BBBB",
		},
	}
	orig.Pointer.RecPointer = orig

	clone := new(MyTestStructRecursive)
	err := deepcopy.Copy(orig, clone, deepcopy.CheckRecursive)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\norig: %+v\ncopy: %+v", *orig, *clone)
	log.Printf("\norig: %+v\ncopy: %+v", *orig.Pointer, *clone.Pointer)
}
```

output:
```bash
orig pointer: 0xc00000e018
copy pointer: 0xc00000e020

orig: {A:AAAA B:BBBB RecPointer:0xc00000e018}
copy: {A:AAAA B:BBBB RecPointer:0xc00000e028}
```

#### Copy `as is`
Sometimes we need to copy struct 'as is', for example, `time.Time` or some own structs with private fields.

```go
package main

import (
	"github.com/Sereger/deepCopy"
	"log"
	"time"
)

type (
	AsIsStruct struct {
		a, b string
	}

	MyTestStructAsIs struct {
		Pointer AsIsStruct
		T       time.Time
	}
)

func main() {
	orig := &MyTestStructAsIs{
		Pointer: AsIsStruct{
			a: "AAAA",
			b: "BBBB",
		},
		T: time.Date(1988, 3, 25, 11, 33, 59, 0, time.Local),
	}

	clone := new(MyTestStructAsIs)
	err := deepcopy.Copy(orig, clone, deepcopy.CopyAsIs(AsIsStruct{}))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\norig: %+v\ncopy: %+v", *orig, *clone)
}
```  

output:
```bash
orig: {Pointer:{a:AAAA b:BBBB} T:2019-03-25 11:33:59 +0300 MSK}
copy: {Pointer:{a:AAAA b:BBBB} T:2019-03-25 11:33:59 +0300 MSK}
```