// +build windows

package extraction

import (
	"strings"
	"reflect"
	"syscall"
	"unsafe"
	"fmt"
	"golang.org/x/sys/windows/registry"
)


const (
	kRegSep = "\\"
	HKEY_CLASSES_ROOT = 0x80000000 + iota
	HKEY_CURRENT_USER
	HKEY_LOCAL_MACHINE
	HKEY_USERS
	HKEY_PERFORMANCE_DATA
	HKEY_CURRENT_CONFIG
	HKEY_DYN_DATA
	HKEY_CURRENT_USER_LOCAL_SETTINGS
	HKEY_PERFORMANCE_TEXT = 0x80000050
	HKEY_PERFORMANCE_NLSTEXT = 0x80000060
)

  var kRegistryHives = map[string]int {
    "HKEY_CLASSES_ROOT": 				HKEY_CLASSES_ROOT,
    "HKEY_CURRENT_CONFIG":  			HKEY_CURRENT_CONFIG,
    "HKEY_CURRENT_USER": 				HKEY_CURRENT_USER,
    "HKEY_CURRENT_USER_LOCAL_SETTINGS": HKEY_CURRENT_USER_LOCAL_SETTINGS,
    "HKEY_LOCAL_MACHINE": 				HKEY_LOCAL_MACHINE,
    "HKEY_PERFORMANCE_DATA":			HKEY_PERFORMANCE_DATA,
    "HKEY_PERFORMANCE_NLSTEXT": 		HKEY_PERFORMANCE_NLSTEXT,
    "HKEY_PERFORMANCE_TEXT": 			HKEY_PERFORMANCE_TEXT,
    "HKEY_USERS": 						HKEY_USERS,
}


func explodeRegistryPath( path string) (rHive, rKey string) {
	toks := strings.Split(path, kRegSep);
	rHive = toks[0];
	toks = append(toks[1:])
	rKey = strings.Join(toks, kRegSep)
	return rHive,rKey
}

func queryKey(keyPath string) (Data, error) {
	var data Data
	hive, key := explodeRegistryPath(keyPath)

	val, ok := kRegistryHives[hive]
	if !ok {
		return data, fmt.Errorf("Key not exists in Hives")
	}
	
	relkey  := registry.Key(val)

	hkey, err := registry.OpenKey(relkey, key, syscall.KEY_READ)
	if err != nil {
		return data, err
	}

	hkey = hkey



	return nil,nil
}



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



