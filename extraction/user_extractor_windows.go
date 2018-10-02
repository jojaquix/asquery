// +build windows

package extraction

import (
	"container/list"
	"reflect"
	"strconv"
	"syscall"
	"unsafe"
	//"golang.org/x/text/encoding/unicode"
	//"golang.org/x/sys/windows"
	//"golang.org/x/sys/windows/registry"
)

// windows implementation
type userExtractorWindows struct {
}

//func NewUserExtracor() UserExtractor {
//	return &userExtractorWindows{}
//}

func GetUsers() (list.List, error) {
	var results list.List
	processedSids := make([]string, 0)
	processLocalAccounts2(processedSids, &results)

	return results, nil
}

//std::string psidToString(PSID sid);
//int getUidFromSid(PSID sid);
//int getGidFromSid(PSID sid);

func processLocalAccounts(processedSids []string, results *list.List) {

	r := Row{}
	r["username"] = "james"
	r["uuid"] = "123"
	results.PushBack(r)

}

func processLocalAccounts2(processedSids []string, results *list.List) {

	//TODO good value for this
	const MAX_PREFERRED_LENGTH = ^DWORD(0)

	var dwUserInfoLevel DWORD = 1
	var dwNumUsersRead WORD = 0
	var dwTotalUsers DWORD = 0
	var resumeHandle DWORD = 0
	var ret NET_API_STATUS
	var userBuffer *BYTE = nil

	for {
		ret = NetUserEnum(nil,
			dwUserInfoLevel,
			0,
			&userBuffer,
			MAX_PREFERRED_LENGTH,
			&dwNumUsersRead,
			&dwTotalUsers,
			&resumeHandle)

		if (ret == NERR_Success || ret == ERROR_MORE_DATA) && userBuffer != nil {

			defer NetApiBufferFree((LPVOID)(unsafe.Pointer(userBuffer)))

			userInfo_ptr := uintptr(unsafe.Pointer(userBuffer))
			var userInfo_slide []USER_INFO_1
			uish := (*reflect.SliceHeader)(unsafe.Pointer(&userInfo_slide))
			uish.Data = userInfo_ptr
			uish.Len = int(dwNumUsersRead)
			uish.Cap = int(dwNumUsersRead)

			for ui, userInfo := range userInfo_slide {

				r := Row{}

				r["uuid"] = Data(strconv.Itoa(10 + ui))

				p := uintptr(unsafe.Pointer(userInfo.Usri1_name))
				var data []uint16
				sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
				sh.Data = p
				sh.Len = 255
				sh.Cap = 255
				username := syscall.UTF16ToString(data)

				r["username"] = Data(username)
				results.PushBack(r)
			}

		}

		if ret != ERROR_MORE_DATA {
			break
		}
	}

}
