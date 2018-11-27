// +build windows

package extraction

import (
	"container/list"
	"golang.org/x/sys/windows"
	"unsafe"
	//"syscall"
	//"golang.org/x/text/encoding/unicode"
	//"golang.org/x/sys/windows"
	//"golang.org/x/sys/windows/registry"
)

const (
	kRegProfilePath = "HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\ProfileList"
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
	processLocalAccounts2(processedSids, &results)
	processRoamingAccounts(processedSids, &results)
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

			defer NetApiBufferFree((LPVOID)(unsafe.Pointer(userBuffer)))
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

				processedSids = append(processedSids, sidString)

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

				results.PushBack(r)
			}

		}

		if ret != ERROR_MORE_DATA {
			break
		}
	}

}

func processRoamingAccounts(processedSids []string, results *list.List) {

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
		if findSid(sidString, processedSids) {
			continue
		}

		r := Row{}

		r["uuid"] = sidString
		r["directory"] = getUserHomeDir(sidString)

		sid, err := windows.StringToSid(sidString)
		if err != nil {
			return
		}

		account, domain, accType, err := sid.LookupAccount("")
		if err != nil {
			return
		}

		account = account
		domain = domain
		accType = accType

		//	var userInfoLevel DWORD = 3
		//	var userBuff *BYTE
		//	ret = NetUserGetInfo(nil,
		//		(*WSTR)(),
		//		userInfoLevel,
		//		&userBuff)
		//

		if !findSidInWellKnwon(sidString) {
			r["type"] = "roaming"
		} else {
			r["type"] = "special"
		}

		//TODO
		r["shell"] = "C:\\Windows\\System32\\cmd.exe"
		r["description"] = ""

		r["username"] = account

		results.PushBack(r)

	}

}
