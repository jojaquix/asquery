package extraction

import (
	"container/list"
	"reflect"
	"unsafe"
)

type Data interface{}
type Row map[string]Data
type ColumnNames []string
type QueryData list.List

//InfoIterable ... is a generic closable interface for iterating over Infos
type InfoIterable interface {
	Next() (*InfoIterable, error)
	ForEach(func(*InfoIterable) error) error
	Close()
}

//this create on slide from existing memeory layout
func createSlide(unsafePtr unsafe.Pointer, slideProto interface{}, size int) reflect.Value {

	slideType := reflect.SliceOf(reflect.TypeOf(slideProto))
	slide := reflect.New(slideType)
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&slide))
	sh.Data = uintptr(unsafePtr)
	sh.Len = size
	sh.Cap = size
	return slide

}
