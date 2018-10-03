// +build windows

package extraction

import (
	"reflect"
	"syscall"
	"unsafe"
)

func StringFromLPWSTR(source LPWSTR, size int) string {
	p := uintptr(unsafe.Pointer(source))
	var data []uint16
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = p
	sh.Len = 255
	sh.Cap = 255
	the_string := syscall.UTF16ToString(data)
	return the_string

}

//TODO is it possible using reflection ?

func CreateUserInfo1SlideFromLPBYTE(source *BYTE, size int) []USER_INFO_1 {
	userInfo_ptr := uintptr(unsafe.Pointer(source))
	var userInfoSlide []USER_INFO_1
	uish := (*reflect.SliceHeader)(unsafe.Pointer(&userInfoSlide))
	uish.Data = userInfo_ptr
	uish.Len = size
	uish.Cap = size
	return userInfoSlide
}

func CreateUserInfo3SlideFromLPBYTE(source *BYTE, size int) []USER_INFO_3 {
	userInfo_ptr := uintptr(unsafe.Pointer(source))
	var userInfoSlide []USER_INFO_3
	uish := (*reflect.SliceHeader)(unsafe.Pointer(&userInfoSlide))
	uish.Data = userInfo_ptr
	uish.Len = size
	uish.Cap = size
	return userInfoSlide
}

func CreateUserInfo4SlideFromLPBYTE(source *BYTE, size int) []USER_INFO_4 {
	userInfo_ptr := uintptr(unsafe.Pointer(source))
	var userInfoSlide []USER_INFO_4
	uish := (*reflect.SliceHeader)(unsafe.Pointer(&userInfoSlide))
	uish.Data = userInfo_ptr
	uish.Len = size
	uish.Cap = size
	return userInfoSlide
}



