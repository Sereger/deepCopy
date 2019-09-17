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

	log.Printf("\norig: %+v\ncopy: %+v", *orig, *clone)
	log.Printf("\norig: %+v\ncopy: %+v", *orig.Pointer, *clone.Pointer)
}
