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
		T: time.Date(2019, 3, 25, 11, 33, 59, 0, time.Local),
	}

	clone := new(MyTestStructAsIs)
	err := deepcopy.Copy(orig, clone, deepcopy.CopyAsIs(AsIsStruct{}))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\norig: %+v\ncopy: %+v", *orig, *clone)
}
