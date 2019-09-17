package deepcopy

import (
	"bytes"
	"encoding/json"
	"testing"
)

type (
	tStructJson struct {
		X          int
		Y          int
		Primitives *primitJson
		Simple     simpleJson
		Intef      interface{}
		I          int
		z          int
	}

	primitJson struct {
		Bts        []byte
		PointerBts *[]byte
		Sl         []int
		SlPoints   []*int
		PointInt   *int
		Int        int
		Str        string
		MpStr      map[int]string
		MpPointer  map[string]*int
		MpSlice    map[uint][]int
		Arr        [4]string
	}

	simpleJson struct {
		A, B string
	}
)

func makePrimitivesJson() *primitJson {
	return &primitJson{
		Bts:        []byte{11, 23, 99, 100},
		PointerBts: &[]byte{11, 23, 99, 100},
		Sl:         []int{1, 2, 3, 4},
		SlPoints:   []*int{intPointer(), intPointer(), intPointer(), intPointer()},
		PointInt:   intPointer(),
		Int:        *intPointer(),
		Str:        "xxxxx",
		MpStr:      map[int]string{1: "x", 2: "U", 3: "axadwa"},
		MpPointer:  map[string]*int{"a": intPointer(), "b": intPointer(), "rrr": intPointer()},
		MpSlice: map[uint][]int{
			1: {1, 2, 3},
			7: {4, 7, 3, 9},
			9: {7, 3, 9, 6, 1},
		},
		Arr: [4]string{"a", "b", "c", "d"},
	}
}

func testStructJson() *tStructJson {
	return &tStructJson{
		X: 99, Y: 100,
		Primitives: makePrimitivesJson(),
		Simple:     simpleJson{"XXX", "YYYY"},
		Intef: &struct {
			X, U int
		}{99, 500},
		I: 100,
		z: -1,
	}
}

func BenchmarkReflectCopy(b *testing.B) {
	v := testStructJson()
	b.Run("reflect-simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			n := new(tStructJson)
			_ = Copy(v, n)
		}
	})
	b.Run("reflect-recursive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			n := new(tStructJson)
			_ = Copy(v, n, CheckRecursive)
		}
	})
	b.Run("jsonEnc-Dec", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bbuf := bytes.NewBuffer([]byte{})
			n := new(tStructJson)
			err := json.NewEncoder(bbuf).Encode(v)
			if err != nil {
				b.Fatal(err)
			}

			err = json.NewDecoder(bbuf).Decode(n)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
