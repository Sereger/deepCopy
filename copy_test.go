package deepcopy

import (
	"github.com/pkg/errors"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

type (
	tStruct struct {
		X          int
		Y          int
		Primitives *primitives
		Simple     simple
		Intef      interface{}
		I          int
		Func       func() error
		Ch         chan string
		z          int
		Nil        nilStructs
		Rec        *recurciveA
		T          time.Time
		D          time.Duration
		MyPriv     myPrivareStruct // will be copied as it is
	}

	nilStructs struct {
		Intef      interface{}
		Ch         chan struct{}
		Func       func() error
		Pointer    *simple
		PointerInt *int
		Slice      []int
	}

	primitives struct {
		Bts          []byte
		PointerBts   *[]byte
		Sl           []int
		SlPoints     []*int
		PointInt     *int
		Int          int
		Str          string
		MpStr        map[int]string
		MpPointer    map[string]*int
		MpSlice      map[uint][]int
		MpKeyPointer map[*int]interface{}
		Arr          [4]string
	}

	simple struct {
		A, B string
	}

	recurciveA struct {
		B   *recurciveB
		Val string
	}

	recurciveB struct {
		A   *recurciveA
		Val string
	}

	myPrivareStruct struct {
		a, b int
	}
)

func makePrimitives() *primitives {
	return &primitives{
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
		MpKeyPointer: map[*int]interface{}{
			intPointer(): "xxx", intPointer(): 42, intPointer(): intPointer(), intPointer(): struct{}{},
		},
		Arr: [4]string{"a", "b", "c", "d"},
	}
}

func intPointer() *int {
	v := rand.Int()
	return &v
}

func testStruct() *tStruct {
	ts := &tStruct{
		X: 99, Y: 100,
		Primitives: makePrimitives(),
		Simple:     simple{"XXX", "YYYY"},
		Intef: &struct {
			X, U int
		}{99, 500},
		I: 100,
		z: -1,
		Func: func() error {
			return errors.New("test error")
		},
		Ch:     make(chan string, 10),
		T:      time.Now(),
		D:      time.Hour,
		MyPriv: myPrivareStruct{1, 2},
	}

	recA := &recurciveA{Val: "A"}
	recB := &recurciveB{Val: "B"}
	recA.B = recB
	recB.A = recA

	ts.Rec = recA
	return ts
}

func TestError(t *testing.T) {
	v := testStruct()
	n := new(struct{ X, Y string })
	err := Copy(v, n, CheckRecursive)
	if err == nil {
		t.Fatal("haven't error, but expected")
	}

	if err.Error() != "types not equal" {
		t.Fatal("incorrect error")
	}
}

func TestSimple(t *testing.T) {
	v := testStruct()
	n := new(tStruct)
	err := Copy(v, n, CheckRecursive, CopyAsIs(myPrivareStruct{}))
	if err != nil {
		t.Fatal(err)
	}

	if v.X != n.X {
		t.Errorf("prop X incorrect: %v, expect %v", n.X, v.X)
	}

	if v.Y != n.Y {
		t.Errorf("prop Y incorrect: %v, expect %v", n.Y, v.Y)
	}

	if v.Simple.A != n.Simple.A {
		t.Errorf("prop Simple.A incorrect: %v, expect %v", n.Simple.A, v.Simple.A)
	}

	if v.Simple.B != n.Simple.B {
		t.Errorf("prop Simple.B incorrect: %v, expect %v", n.Simple.B, v.Simple.B)
	}

	if v.Primitives == n.Primitives {
		t.Error("pointer Primitives equal")
	}

	if v.Primitives == n.Primitives {
		t.Error("pointer Primitives equal")
	}

	if !reflect.DeepEqual(v.Intef, n.Intef) {
		t.Error("interface Intef are not equal")
	}

	for i, val := range v.Primitives.Bts {
		nVal := n.Primitives.Bts[i]
		if val != nVal {
			t.Errorf("Bts[%d] not equals: %v, expect: %v", i, nVal, val)
		}
	}

	if v.Primitives.PointerBts == n.Primitives.PointerBts {
		t.Error("pointer Primitives.PointerBts equal")
	}

	for i, val := range *v.Primitives.PointerBts {
		bts := *n.Primitives.PointerBts
		if val != bts[i] {
			t.Errorf("PointerBts[%d] not equals: %v, expect: %v", i, bts[i], val)
		}
	}

	for i, val := range v.Primitives.Sl {
		nVal := n.Primitives.Sl[i]
		if val != nVal {
			t.Errorf("Primitives.Sl[%d] not equals: %v, expect: %v", i, nVal, val)
		}
	}

	for i, val := range v.Primitives.SlPoints {
		nVal := n.Primitives.SlPoints[i]
		if val == nVal {
			t.Errorf("pointer Primitives.SlPoints[%d] equal", i)
			if *val != *nVal {
				t.Errorf("values Primitives.SlPoints[%d] not equal: %v, expect: %v", i, nVal, val)
			}
		}
	}

	if v.Primitives.PointInt == n.Primitives.PointInt {
		t.Error("pointer Primitives.PointInt equal")
	}

	if v.Primitives.Int != n.Primitives.Int {
		t.Errorf("prop Primitives.Int incorrect: %v, expect %v", n.Primitives.Int, v.Primitives.Int)
	}

	if v.Primitives.Str != n.Primitives.Str {
		t.Errorf("prop Primitives.Str incorrect: %v, expect %v", n.Primitives.Str, v.Primitives.Str)
	}

	for key, val := range v.Primitives.MpStr {
		nVal, ok := n.Primitives.MpStr[key]
		if !ok {
			t.Errorf("haven't value in map v.Primitives.MpStr: %d", key)
		}

		if val != nVal {
			t.Errorf("values v.Primitives.MpStr[%d] not equal: %v, expect: %v", key, nVal, val)
		}
	}

	for key, val := range v.Primitives.MpPointer {
		nVal, ok := n.Primitives.MpPointer[key]
		if !ok {
			t.Errorf("haven't value in map v.Primitives.MpPointer: %s", key)
		}

		if val == nVal {
			t.Errorf("pointers Primitives.MpPointer[%s] equal: %v == %v", key, val, nVal)
		}

		if *val != *nVal {
			t.Errorf("values v.Primitives.MpPointer[%s] not equal: %v, expect: %v", key, nVal, val)
		}
	}

	for i, val := range v.Primitives.MpSlice {
		nVal, ok := n.Primitives.MpSlice[i]
		if !ok {
			t.Errorf("haven't value in map v.Primitives.MpSlice: %d", i)
		}

		for j, subVal := range val {
			nSubVal := nVal[j]

			if subVal != nSubVal {
				t.Errorf("values v.Primitives.MpSlice[%d][%d] not equal: %v, expect: %v", i, j, nSubVal, subVal)
			}
		}
	}

	for i, val := range v.Primitives.Arr {
		nVal := n.Primitives.Arr[i]

		if val != nVal {
			t.Errorf("values v.Primitives.Arr[%d] not equal: %v, expect: %v", i, nVal, val)
		}
	}

	err1, err2 := v.Func(), n.Func()
	if err2 == nil {
		t.Error("failed call n.Func")
	} else if err1.Error() != err2.Error() {
		t.Errorf("incorrect result of fn n.Func: [%s], expect [%s]", err2.Error(), err1.Error())
	}

	n.Ch <- "xxx"
	chVal, ok := <-n.Ch
	if !ok {
		t.Error("can't receive from n.Ch")
	} else {
		if chVal != "xxx" {
			t.Errorf("receive incorrect value: [%s]", chVal)
		}
	}

	if v.Rec == n.Rec {
		t.Error("pointer Rec equal")
	}

	if v.Rec.Val != n.Rec.Val {
		t.Errorf("prop Rec.Val incorrect: %v, expect %v", n.Rec.Val, v.Rec.Val)
	}
	if v.Rec.B.Val != n.Rec.B.Val {
		t.Errorf("prop Rec.B.Val incorrect: %v, expect %v", n.Rec.B.Val, v.Rec.B.Val)
	}
	if n.Rec.B.A.B.A != n.Rec.B.A {
		t.Errorf("incorrect recursive pointers: %v != %v", n.Rec.B.A, n.Rec)
	}

	if !v.T.Equal(n.T) {
		t.Errorf("prop T incorrect: %s, expect %s", n.T, v.T)
	}

	if v.D != n.D {
		t.Errorf("prop D incorrect: %s, expect %s", n.D, v.D)
	}

	if v.MyPriv.a != n.MyPriv.a {
		t.Errorf("prop MyPriv.a incorrect: %v, expect %v", n.MyPriv.a, v.MyPriv.a)
	}

	if v.MyPriv.b != n.MyPriv.b {
		t.Errorf("prop MyPriv.b incorrect: %v, expect %v", n.MyPriv.b, v.MyPriv.b)
	}

	t.Logf("\n%+v\n%+v", v, n)
}
