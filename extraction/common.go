package extraction

import (
	"reflect"
	"unsafe"
)

type Data interface{}
type Row map[string]Data
type ColumnNames []string

//InfoIterable ... is a generic closable interface for iterating over Infos
type InfoIterable interface {
	Next() (*InfoIterable, error)
	ForEach(func(*InfoIterable) error) error
	Close()
}

//this create on slide from existing memory layout
func createSlide(unsafePtr unsafe.Pointer, slideProto interface{}, size int) reflect.Value {

	slideType := reflect.SliceOf(reflect.TypeOf(slideProto))
	slide := reflect.New(slideType)
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&slide))
	sh.Data = uintptr(unsafePtr)
	sh.Len = size
	sh.Cap = size
	return slide
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
