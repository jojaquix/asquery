// +build windows

package extraction

import (
	"container/list"
	"golang.org/x/sys/windows"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
	//"syscall"
	//"golang.org/x/text/encoding/unicode"
	//"golang.org/x/sys/windows"
	//"golang.org/x/sys/windows/registry"
)

const (
	kRegProfilePath   = "HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\ProfileList"
	NERR_UserNotFound = 2221
)

var kWellKnownSids = [...]string{
	"S-1-5-1",
	"S-1-5-2",
	"S-1-5-3",
	"S-1-5-4",
	"S-1-5-6",
	"S-1-5-7",
	"S-1-5-8",
	"S-1-5-9",
	"S-1-5-10",
	"S-1-5-11",
	"S-1-5-12",
	"S-1-5-13",
	"S-1-5-18",
	"S-1-5-19",
	"S-1-5-20",
	"S-1-5-21",
	"S-1-5-32",
}

func findSid(sid string, where []string) bool {
	for _, v := range where {
		if v == sid {
			return true
		}
	}
	return false
}

func findSidInWellKnwon(sid string) bool {
	for _, v := range kWellKnownSids {
		if v == sid {
			return true
		}
	}
	return false
}

func getUserHomeDir(sid string) string {

	keyResult, err := queryKey(kRegProfilePath + kRegSep + sid)
	if err != nil {
		return ""
	}

	for k := keyResult.Front(); k != nil; k = k.Next() {

		row := k.Value.(Row)
		if row["name"] == "ProfileImagePath" {
			return string(row["data"].(string))
		}
	}

	return ""

}

func GetUsers() (list.List, error) {
	var results list.List
	processedSids := make([]string, 0)
	processLocalAccounts2(&processedSids, &results)
	processRoamingAccounts(&processedSids, &results)
	return results, nil
}

func getUidFromSid(sid *windows.SID) int64 {
	var userInfoLevel DWORD = 3
	var userBuffer *BYTE
	var uid int64 = -1
	account, _, _, err := sid.LookupAccount("")
	if err != nil {
		return 0
	}

	accountUtf16Slide := utf16FromString(account)
	accountUtf16SlideHeader := (*reflect.SliceHeader)(unsafe.Pointer(&accountUtf16Slide))
	ret := NetUserGetInfo(nil,
		(*WSTR)((unsafe.Pointer(accountUtf16SlideHeader.Data))),
		userInfoLevel,
		&userBuffer)

	if ret == NERR_UserNotFound {
		sidStr, err := sid.String()
		if err == nil {
			toks := strings.Split(sidStr, "-")
			toks = toks
			value, err := strconv.ParseInt(toks[len(toks)-1], 10, 64)
			if err == nil {
				return value
			}
		}
	} else if ret == NERR_Success {
		userInfo3Slide := CreateUserInfo3SlideFromLPBYTE(userBuffer, 1)
		uid = int64(userInfo3Slide[0].usri3_user_id)
	}

	if ret == NERR_Success && userBuffer != nil {
		NetApiBufferFree((LPVOID)(unsafe.Pointer(userBuffer)))
	}

	return uid
}

func getGidFromSid(sid *windows.SID) int64 {
	var userInfoLevel DWORD = 3
	var userBuffer *BYTE
	var gid int64 = -1
	account, _, _, err := sid.LookupAccount("")
	if err != nil {
		return 0
	}

	accountUtf16Slide := utf16FromString(account)
	accountUtf16SlideHeader := (*reflect.SliceHeader)(unsafe.Pointer(&accountUtf16Slide))
	ret := NetUserGetInfo(nil,
		(*WSTR)((unsafe.Pointer(accountUtf16SlideHeader.Data))),
		userInfoLevel,
		&userBuffer)

	if ret == NERR_UserNotFound {
		sidStr, err := sid.String()
		if err == nil {
			toks := strings.Split(sidStr, "-")
			toks = toks
			value, err := strconv.ParseInt(toks[len(toks)-1], 10, 64)
			if err == nil {
				return value
			}
		}
	} else if ret == NERR_Success {
		userInfo3Slide := CreateUserInfo3SlideFromLPBYTE(userBuffer, 1)
		gid = int64(userInfo3Slide[0].usri3_primary_group_id)
	}

	if ret == NERR_Success && userBuffer != nil {
		NetApiBufferFree((LPVOID)(unsafe.Pointer(userBuffer)))
	}

	return gid
}

func processLocalAccounts(processedSids []string, results *list.List) {

	r := Row{}
	r["username"] = "james"
	r["uuid"] = "123"
	results.PushBack(r)

}

func processLocalAccounts2(processedSids *[]string, results *list.List) {

	const MAX_PREFERRED_LENGTH = ^DWORD(0)

	var dwUserInfoLevel DWORD = 3
	var dwNumUsersRead WORD
	var dwTotalUsers DWORD
	var resumeHandle DWORD
	var ret NET_API_STATUS
	var userBuffer *BYTE

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

			userInfoSlide := CreateUserInfo3SlideFromLPBYTE(userBuffer, int(dwNumUsersRead))

			for _, userInfo := range userInfoSlide {

				// User level 4 contains the SID value
				var dwDetailedUserInfoLevel DWORD = 4
				var userLvl4Buff *BYTE
				ret = NetUserGetInfo(nil,
					(*WSTR)(unsafe.Pointer(userInfo.usri3_name)),
					dwDetailedUserInfoLevel,
					&userLvl4Buff)

				if ret != NERR_Success || userLvl4Buff == nil {
					if userLvl4Buff != nil {
						NetApiBufferFree((LPVOID)(unsafe.Pointer(userLvl4Buff)))
					}
					//TODO loging
					continue
				}

				// Will return empty string on fail
				userInfo4Slide := CreateUserInfo4SlideFromLPBYTE(userLvl4Buff, 1)
				sid := userInfo4Slide[0].usri4_user_sid
				sidString, err := sid.String()
				if err != nil {
					sidString = ""
				}

				*processedSids = append(*processedSids, sidString)

				r := Row{}
				r["uuid"] = sidString
				r["username"] = StringFromLPWSTR(userInfo.usri3_name, 255)
				r["uid"] = int64(userInfo.usri3_user_id)
				r["gid"] = int64(userInfo.usri3_primary_group_id)
				r["uid_signed"] = r["uid"]
				r["gid_signed"] = r["gid"]
				r["description"] = StringFromLPWSTR(userInfo4Slide[0].usri4_comment, 2048)
				r["directory"] = getUserHomeDir(sidString)
				r["shell"] = "C:\\Windows\\System32\\cmd.exe"
				r["type"] = "local"

				if userLvl4Buff != nil {
					NetApiBufferFree((LPVOID)(unsafe.Pointer(userLvl4Buff)))
				}

				results.PushBack(r)
			}
			if userBuffer != nil {
				NetApiBufferFree((LPVOID)(unsafe.Pointer(userBuffer)))
			}
		}

		if ret != ERROR_MORE_DATA {
			break
		}
	}

}

func processRoamingAccounts(processedSids *[]string, results *list.List) {

	keyResult, err := queryKey(kRegProfilePath)
	if err != nil {
		return
	}

	for k := keyResult.Front(); k != nil; k = k.Next() {
		row := k.Value.(Row)
		if row["type"].(string) != "subKey" {
			continue
		}

		sidString := row["name"].(string)
		if findSid(sidString, *processedSids) {
			continue
		}

		r := Row{}

		r["uuid"] = sidString
		r["directory"] = getUserHomeDir(sidString)

		sid, err := windows.StringToSid(sidString)
		if err != nil {
			return
		}

		r["uid"] = getUidFromSid(sid)
		r["gid"] = getGidFromSid(sid)
		r["uid_signed"] = r["uid"]
		r["gid_signed"] = r["gid"]

		if !findSidInWellKnwon(sidString) {
			r["type"] = "roaming"
		} else {
			r["type"] = "special"
		}

		//TODO
		r["shell"] = "C:\\Windows\\System32\\cmd.exe"
		r["description"] = ""

		account, _, _, err := sid.LookupAccount("")
		if err != nil {
			r["username"] = ""
		} else {
			r["username"] = account
		}

		results.PushBack(r)
	}

}
