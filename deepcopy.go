package deepcopy

import (
	"github.com/pkg/errors"
	"reflect"
	"time"
)

func Copy(val, target interface{}, opts ...option) (err error) {
	defer func() {
		// reflect паникует в случае ошибки, т.ч. делаем recover
		if recErr := recover(); recErr != nil {
			err = errors.Errorf("con't copy: [%v]", recErr)
		}
	}()

	v, t := reflect.ValueOf(val), reflect.ValueOf(target)
	if v.Type() != t.Type() {
		return errors.New("types not equal")
	}

	p := makeProcessor()
	for _, opt := range opts {
		opt(p)
	}
	p.reflectCopy(v, t)
	return
}

var CheckRecursive = func(p *processor) {
	p.cache = make(map[uintptr]reflect.Value)
}

func CopyAsIs(types ...interface{}) option {
	return func(p *processor) {
		for _, v := range types {
			t := reflect.ValueOf(v).Type()
			p.asIs[t] = struct{}{}
		}
	}
}

func makeProcessor() *processor {
	return &processor{
		asIs: map[reflect.Type]struct{}{
			reflect.TypeOf(time.Time{}): {},
		},
	}
}

type (
	processor struct {
		cache map[uintptr]reflect.Value
		asIs  map[reflect.Type]struct{}
	}

	option func(p *processor)
)

func (p *processor) reflectCopy(v, t reflect.Value) {
	switch v.Kind() {
	case reflect.Struct:
		vt := v.Type()
		if _, ok := p.asIs[vt]; ok {
			t.Set(v)
			return
		}

		for i := 0; i < v.NumField(); i++ {
			if !t.Field(i).CanSet() {
				continue
			}

			p.reflectCopy(v.Field(i), t.Field(i))
		}
	case reflect.Ptr:
		if v.IsNil() {
			return
		}

		if t.IsNil() {
			t.Set(reflect.New(v.Elem().Type()))
		}

		if p.cache != nil {
			ptr, ok := p.cache[v.Pointer()]
			if ok {
				t.Elem().Set(ptr.Elem())
			} else {
				p.cache[v.Pointer()] = t
				p.reflectCopy(v.Elem(), t.Elem())
			}
		} else {
			p.reflectCopy(v.Elem(), t.Elem())
		}
	case reflect.Slice:
		if v.IsNil() {
			return
		}

		t.Set(reflect.MakeSlice(v.Type(), v.Len(), v.Cap()))
		for i := 0; i < v.Len(); i++ {
			p.reflectCopy(v.Index(i), t.Index(i))
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			p.reflectCopy(v.Index(i), t.Index(i))
		}
	case reflect.Map:
		if v.IsNil() {
			return
		}

		t.Set(reflect.MakeMapWithSize(v.Type(), v.Len()))
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)

			nKey := reflect.New(key.Type()).Elem()
			p.reflectCopy(key, nKey)

			nVal := reflect.New(val.Type()).Elem()
			p.reflectCopy(val, nVal)

			t.SetMapIndex(nKey, nVal)
		}
	case reflect.Interface:
		if v.IsNil() || !v.Elem().IsValid() {
			return
		}
		nv := reflect.New(v.Elem().Type())
		p.reflectCopy(v.Elem(), nv.Elem())
		t.Set(nv.Elem())
	case reflect.Chan:
		if !v.IsNil() {
			t.Set(reflect.MakeChan(v.Type(), v.Cap()))
		}
	case reflect.Func:
		if v.IsNil() {
			return
		}

		t.Set(v)
	default:
		if v.Kind() < reflect.Array || v.Kind() == reflect.String {
			t.Set(v)
		}
	}
}
