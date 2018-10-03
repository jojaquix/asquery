// +build windows

package extraction

import (
	"container/list"
	"unsafe"
	//"golang.org/x/text/encoding/unicode"
	//"golang.org/x/sys/windows"
	//"golang.org/x/sys/windows/registry"
)

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
			var dwDetailedUserInfoLevel DWORD = 4;
			var userLvl4Buff *BYTE
			ret = NetUserGetInfo(nil,
									(*WSTR)(unsafe.Pointer(userInfo.usri3_name)),
									dwDetailedUserInfoLevel,
									&userLvl4Buff);
	
			if ret != NERR_Success || userLvl4Buff == nil {
				if userLvl4Buff != nil {
					NetApiBufferFree((LPVOID)(unsafe.Pointer(userLvl4Buff)));
				}
				//TODO loging
				continue;
			}

			// Will return empty string on fail
			userInfo4Slide := CreateUserInfo4SlideFromLPBYTE(userLvl4Buff,1)
			sid := userInfo4Slide[0].usri4_user_sid
			sidString, err := sid.String()
			if (err != nil) {
				sidString = ""
			}

			processedSids= append(processedSids, sidString);
	

			r := Row{}
			r["uuid"] = sidString				
			r["username"] = StringFromLPWSTR(userInfo.usri3_name, 255)
			r["uid"] = int64(userInfo.usri3_user_id)
			r["gid"] = int64(userInfo.usri3_primary_group_id)
			r["uid_signed"] = r["uid"]
			r["gid_signed"] = r["gid"]
			r["description"] =	StringFromLPWSTR(userInfo4Slide[0].usri4_comment, 2048)
			r["directory"] = ""
			r["shell"] = "C:\\Windows\\System32\\cmd.exe";
			r["type"] = "local";

			results.PushBack(r)
			}

		}

		if ret != ERROR_MORE_DATA {
			break
		}
	}

}
