package main

import (
	"github.com/Sereger/deepCopy"
	"log"
	"unsafe"
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

	log.Printf("\norig pointer: %v\ncopy pointer: %v", unsafe.Pointer(orig), unsafe.Pointer(clone))
	log.Printf("\norig: %+v\ncopy: %+v", *orig.Pointer, *clone.Pointer)
}
