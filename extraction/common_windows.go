// +build windows

package extraction

import (
	"container/list"
	"encoding/binary"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"reflect"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	kRegSep         = "\\"
	kDefaultRegName = "(Default)"
)

const (
	kRegMaxRecursiveDepth = 32
)

//the key numbers must be alone in one const declaration because iota use
const (
	HKEY_CLASSES_ROOT = 0x80000000 + iota
	HKEY_CURRENT_USER
	HKEY_LOCAL_MACHINE
	HKEY_USERS
	HKEY_PERFORMANCE_DATA
	HKEY_CURRENT_CONFIG
	HKEY_DYN_DATA
	HKEY_CURRENT_USER_LOCAL_SETTINGS
	HKEY_PERFORMANCE_TEXT    = 0x80000050
	HKEY_PERFORMANCE_NLSTEXT = 0x80000060
)

var kRegistryHives = map[string]registry.Key{
	"HKEY_CLASSES_ROOT":     registry.CLASSES_ROOT,
	"HKEY_CURRENT_USER":     registry.CURRENT_USER,
	"HKEY_CURRENT_CONFIG":   registry.CURRENT_CONFIG,
	"HKEY_LOCAL_MACHINE":    registry.LOCAL_MACHINE,
	"HKEY_USERS":            registry.USERS,
	"HKEY_PERFORMANCE_DATA": registry.PERFORMANCE_DATA,

	"HKEY_CURRENT_USER_LOCAL_SETTINGS": registry.Key(HKEY_CURRENT_USER_LOCAL_SETTINGS),
	"HKEY_PERFORMANCE_NLSTEXT":         registry.Key(HKEY_PERFORMANCE_NLSTEXT),
	"HKEY_PERFORMANCE_TEXT":            registry.Key(HKEY_PERFORMANCE_TEXT),
}

var kRegistryHives2 = map[string]string{
	"HKEY_CLASSES_ROOT":     "registry.Key(HKEY_CLASSES_ROOT)",
	"HKEY_CURRENT_USER":     "registry.Key(HKEY_CURRENT_USER)",
	"HKEY_CURRENT_CONFIG":   "registry.Key(HKEY_CURRENT_CONFIG)",
	"HKEY_LOCAL_MACHINE":    "registry.Key(HKEY_LOCAL_MACHINE)",
	"HKEY_USERS":            "registry.Key(HKEY_USERS)",
	"HKEY_PERFORMANCE_DATA": "registry.Key(HKEY_PERFORMANCE_DATA)",

	"HKEY_CURRENT_USER_LOCAL_SETTINGS": "registry.Key(HKEY_CURRENT_USER_LOCAL_SETTINGS)",
	"HKEY_PERFORMANCE_NLSTEXT":         "registry.Key(HKEY_PERFORMANCE_NLSTEXT)",
	"HKEY_PERFORMANCE_TEXT":            "registry.Key(HKEY_PERFORMANCE_TEXT)",
}

var kRegistryStringTypes = [...]int{
	syscall.REG_SZ, syscall.REG_MULTI_SZ, syscall.REG_EXPAND_SZ,
}

var kRegistryTypes = map[int]string{
	syscall.REG_BINARY:                   "REG_BINARY",
	syscall.REG_DWORD:                    "REG_DWORD",
	syscall.REG_DWORD_BIG_ENDIAN:         "REG_DWORD_BIG_ENDIAN",
	syscall.REG_EXPAND_SZ:                "REG_EXPAND_SZ",
	syscall.REG_LINK:                     "REG_LINK",
	syscall.REG_MULTI_SZ:                 "REG_MULTI_SZ",
	syscall.REG_NONE:                     "REG_NONE",
	syscall.REG_QWORD:                    "REG_QWORD",
	syscall.REG_SZ:                       "REG_SZ",
	syscall.REG_FULL_RESOURCE_DESCRIPTOR: "REG_FULL_RESOURCE_DESCRIPTOR",
	syscall.REG_RESOURCE_LIST:            "REG_RESOURCE_LIST",
}

var kClassKeys = [...]string{
	"HKEY_USERS\\%\\SOFTWARE\\Classes\\CLSID",
	"HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\CLSID",
}

var kClassExecSubKeys = [...]string{
	"InProcServer%", "InProcHandler%", "LocalServer%",
}

func explodeRegistryPath(path string) (rHive, rKey string) {
	toks := strings.Split(path, kRegSep)
	rHive = toks[0]
	toks = append(toks[1:])
	rKey = strings.Join(toks, kRegSep)
	return rHive, rKey
}

func queryKey(keyPath string) (*list.List, error) {
	var dataQuery *list.List
	hive, key := explodeRegistryPath(keyPath)

	base, ok := kRegistryHives[hive]
	if !ok {
		return nil, fmt.Errorf("Key not exists in Hives")
	}

	hkey, err := registry.OpenKey(base, key, syscall.KEY_READ)
	if err != nil {
		return nil, err
	}

	subKeyNames, err := hkey.ReadSubKeyNames(-1)
	if err != nil {
		return nil, err
	}

	dataQuery = list.New()

	if len(subKeyNames) > 0 {
		for i := 0; i < len(subKeyNames); i++ {
			subKey, err := registry.OpenKey(hkey, subKeyNames[i], syscall.KEY_READ)
			if err != nil {
				return nil, err
			}
			subKeyInfo, err := subKey.Stat()

			var r Row
			r = make(Row)
			r["key"] = keyPath
			r["type"] = "subKey"
			r["name"] = subKeyNames[i]
			r["path"] = keyPath + kRegSep + subKeyNames[i]
			r["mtime"] = subKeyInfo.ModTime()

			subKey.Close()

			dataQuery.PushBack(r)

		}
	}

	keyInfo, err := hkey.Stat()

	if keyInfo.ValueCount <= 0 {
		return dataQuery, nil
	}

	valueNames, err := hkey.ReadValueNames(-1)

	var buf []byte
	buf = make([]byte, keyInfo.MaxValueLen)

	for i := 0; i < len(valueNames); i++ {
		buf[0] = 0
		_, valtype, err := hkey.GetValue(valueNames[i], nil)

		if err != nil {
			return nil, err
		}

		var r Row
		r = make(Row)
		r["key"] = keyPath
		r["name"] = valueNames[i]
		r["path"] = keyPath + kRegSep + valueNames[i]
		r["mtime"] = ""

		valueTypeStr, ok := kRegistryTypes[int(valtype)]

		if !ok {
			r["type"] = "UNKNOWN"
		} else {
			r["type"] = valueTypeStr
		}

		switch valtype {
		//TODO other register values types as osquery does
		case registry.LINK:
			r["data"] = "No Implemented yet :("

		case registry.EXPAND_SZ, registry.SZ:
			strValue, _, err := hkey.GetStringValue(valueNames[i])
			if err != nil {
				return nil, err
			}
			r["data"] = strValue

		case registry.MULTI_SZ:
			_, _, err := hkey.GetValue(valueNames[i], buf)
			if err != nil {
				return nil, err
			}
			//TODO this is to naive check latter
			str := string(buf)
			r["data"] = str

			//			multiSzStrs := []string{}
			//			str := ""
			//			last_was_null:= false
			//			for i := 0; i< bytes; i++ {
			//
			//				if buf[i] != 0x00 {
			//					str = str + string(buf[i])
			//					last_was_null= false
			//				} else if last_was_null {
			//					multiSzStrs= append(multiSzStrs, str)
			//					str = ""
			//					last_was_null= true
			//				}
			//			}
			//			r["data"] = strings.Join(multiSzStrs,",")

		case registry.DWORD, registry.QWORD:
			intValue, _, err := hkey.GetIntegerValue(valueNames[i])
			if err != nil {
				return nil, err
			}
			r["data"] = string(intValue)

		case registry.BINARY:
			binValue, _, err := hkey.GetBinaryValue(valueNames[i])
			if err != nil {
				return nil, err
			}
			r["data"] = string(binValue)

		case registry.NONE:
			r["data"] = "(zero-length binary value)"

		default:
			r["data"] = ""
		}

		dataQuery.PushBack(r)

	}

	return dataQuery, nil
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

// uf16PtrToString creates a Go string from a pointer to a UTF16 encoded zero-terminated string.
// Such pointers are returned from the Windows API calls.
// The function creates a copy of the string.
func utf16PtrToString(wstr *uint16) string {
	if wstr != nil {
		for len := 0; ; len++ {
			ptr := unsafe.Pointer(uintptr(unsafe.Pointer(wstr)) + uintptr(len)*unsafe.Sizeof(*wstr)) // see https://golang.org/pkg/unsafe/#Pointer (3)
			if *(*uint16)(ptr) == 0 {
				return string(utf16.Decode(*(*[]uint16)(unsafe.Pointer(&reflect.SliceHeader{
					Data: uintptr(unsafe.Pointer(wstr)),
					Len:  len,
					Cap:  len,
				}))))
			}
		}
	}
	return ""
}

// utf16ToByte creates a byte array from a given UTF 16 char array.
func utf16ToByte(wstr []uint16) (result []byte) {
	result = make([]byte, len(wstr)*2)
	for i := range wstr {
		binary.LittleEndian.PutUint16(result[(i*2):(i*2)+2], wstr[i])
	}
	return
}

// utf16FromString creates a UTF16 char array from a string.
func utf16FromString(str string) []uint16 {
	out, err := syscall.UTF16FromString(str)
	if err != nil {
		return make([]uint16, 0)
	} else {
		return out
	}
}

func expandRegistryGlobs() {
	return
}
