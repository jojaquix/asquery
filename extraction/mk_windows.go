package extraction

import (
	"syscall"
)

// to genarate
// compile mksyscall_windows using go build .\mksyscall_windows.go
// and then use the executable over this file: .\mksyscall_windows.exe  -output net_windows.go .\mk_windows.go

type (
	LPVOID         uintptr	
	LMSTR          *uint16
	WORD           uint16
	DWORD          uint32
	BYTE           byte
	LPBYTE         *byte
	LPDWORD        *uint32
	LPWSTR         *uint16
	LPCWSTR        *uint16
	WSTR           uint16	
	NET_API_STATUS DWORD


	USER_INFO_1 struct {
		Usri1_name         LPWSTR
		Usri1_password     LPWSTR
		Usri1_password_age DWORD
		Usri1_priv         DWORD
		Usri1_home_dir     LPWSTR
		Usri1_comment      LPWSTR
		Usri1_flags        DWORD
		Usri1_script_path  LPWSTR
	}


	USER_INFO_3 struct {
		  usri3_name				LPWSTR 
		  usri3_password			LPWSTR 
		  usri3_password_age		DWORD  
		  usri3_priv				DWORD  
		  usri3_home_dir			LPWSTR 
		  usri3_comment				LPWSTR 
		  usri3_flags				DWORD  
		  usri3_script_path			LPWSTR 
		  usri3_auth_flags			DWORD  
		  usri3_full_name			LPWSTR 
		  usri3_usr_comment			LPWSTR 
		  usri3_parms				LPWSTR 
		  usri3_workstations		LPWSTR 
		  usri3_last_logon			DWORD  
		  usri3_last_logoff			DWORD  
		  usri3_acct_expires		DWORD  
		  usri3_max_storage			DWORD  
		  usri3_units_per_week		DWORD  
		  usri3_logon_hours			*BYTE  
		  usri3_bad_pw_count		DWORD  
		  usri3_num_logons			DWORD  
		  usri3_logon_server		LPWSTR 
		  usri3_country_code		DWORD  
		  usri3_code_page			DWORD  
		  usri3_user_id				DWORD  
		  usri3_primary_group_id	DWORD  
		  usri3_profile				LPWSTR 
		  usri3_home_dir_drive		LPWSTR 
		  usri3_password_expired	DWORD  
	}


	USER_INFO_4 struct {
		usri4_name				LPWSTR
		usri4_password			LPWSTR
		usri4_password_age		DWORD 
		usri4_priv				DWORD 
		usri4_home_dir			LPWSTR
		usri4_comment			LPWSTR
		usri4_flags				DWORD 
		usri4_script_path		LPWSTR
		usri4_auth_flags		DWORD 
		usri4_full_name			LPWSTR
		usri4_usr_comment		LPWSTR
		usri4_parms				LPWSTR
		usri4_workstations		LPWSTR
		usri4_last_logon		DWORD 
		usri4_last_logoff		DWORD 
		usri4_acct_expires		DWORD 
		usri4_max_storage		DWORD 
		usri4_units_per_week	DWORD 
		usri4_logon_hours		*BYTE 
		usri4_bad_pw_count		DWORD 
		usri4_num_logons		DWORD 
		usri4_logon_server		LPWSTR
		usri4_country_code		DWORD 
		usri4_code_page			DWORD 
		usri4_user_sid			*syscall.SID  
		usri4_primary_group_id	DWORD 
		usri4_profile			LPWSTR
		usri4_home_dir_drive	LPWSTR
		usri4_password_expired	DWORD 
	}	

	GROUP_USERS_INFO_0 struct {
		Grui0_name LPWSTR
	}

	USER_INFO_1003 struct {
		Usri1003_password LPWSTR
	}
)

const (
	// from LMaccess.h

	USER_PRIV_GUEST = 0
	USER_PRIV_USER  = 1
	USER_PRIV_ADMIN = 2

	UF_SCRIPT                          = 0x0001
	UF_ACCOUNTDISABLE                  = 0x0002
	UF_HOMEDIR_REQUIRED                = 0x0008
	UF_LOCKOUT                         = 0x0010
	UF_PASSWD_NOTREQD                  = 0x0020
	UF_PASSWD_CANT_CHANGE              = 0x0040
	UF_ENCRYPTED_TEXT_PASSWORD_ALLOWED = 0x0080

	UF_TEMP_DUPLICATE_ACCOUNT    = 0x0100
	UF_NORMAL_ACCOUNT            = 0x0200
	UF_INTERDOMAIN_TRUST_ACCOUNT = 0x0800
	UF_WORKSTATION_TRUST_ACCOUNT = 0x1000
	UF_SERVER_TRUST_ACCOUNT      = 0x2000

	UF_ACCOUNT_TYPE_MASK = UF_TEMP_DUPLICATE_ACCOUNT |
		UF_NORMAL_ACCOUNT |
		UF_INTERDOMAIN_TRUST_ACCOUNT |
		UF_WORKSTATION_TRUST_ACCOUNT |
		UF_SERVER_TRUST_ACCOUNT

	UF_DONT_EXPIRE_PASSWD                     = 0x10000
	UF_MNS_LOGON_ACCOUNT                      = 0x20000
	UF_SMARTCARD_REQUIRED                     = 0x40000
	UF_TRUSTED_FOR_DELEGATION                 = 0x80000
	UF_NOT_DELEGATED                          = 0x100000
	UF_USE_DES_KEY_ONLY                       = 0x200000
	UF_DONT_REQUIRE_PREAUTH                   = 0x400000
	UF_PASSWORD_EXPIRED                       = 0x800000
	UF_TRUSTED_TO_AUTHENTICATE_FOR_DELEGATION = 0x1000000
	UF_NO_AUTH_DATA_REQUIRED                  = 0x2000000
	UF_PARTIAL_SECRETS_ACCOUNT                = 0x4000000
	UF_USE_AES_KEYS                           = 0x8000000

	UF_SETTABLE_BITS = UF_SCRIPT |
		UF_ACCOUNTDISABLE |
		UF_LOCKOUT |
		UF_HOMEDIR_REQUIRED |
		UF_PASSWD_NOTREQD |
		UF_PASSWD_CANT_CHANGE |
		UF_ACCOUNT_TYPE_MASK |
		UF_DONT_EXPIRE_PASSWD |
		UF_MNS_LOGON_ACCOUNT |
		UF_ENCRYPTED_TEXT_PASSWORD_ALLOWED |
		UF_SMARTCARD_REQUIRED |
		UF_TRUSTED_FOR_DELEGATION |
		UF_NOT_DELEGATED |
		UF_USE_DES_KEY_ONLY |
		UF_DONT_REQUIRE_PREAUTH |
		UF_PASSWORD_EXPIRED |
		UF_TRUSTED_TO_AUTHENTICATE_FOR_DELEGATION |
		UF_NO_AUTH_DATA_REQUIRED |
		UF_USE_AES_KEYS |
		UF_PARTIAL_SECRETS_ACCOUNT

	FILTER_TEMP_DUPLICATE_ACCOUNT    = (0x0001)
	FILTER_NORMAL_ACCOUNT            = (0x0002)
	FILTER_INTERDOMAIN_TRUST_ACCOUNT = (0x0008)
	FILTER_WORKSTATION_TRUST_ACCOUNT = (0x0010)
	FILTER_SERVER_TRUST_ACCOUNT      = (0x0020)

	LG_INCLUDE_INDIRECT = (0x0001)

	NERR_Success = 0

	ERROR_MORE_DATA = 234 // dderror

)

//sys NetApiBufferFree(Buffer LPVOID) (status NET_API_STATUS) = netapi32.NetApiBufferFree
//NetUserAdd(servername LMSTR, level DWORD, buf LPBYTE, parm_err LPDWORD) (status NET_API_STATUS) = netapi32.NetUserAdd
//NetUserChangePassword(domainname LPCWSTR, username LPCWSTR, oldpassword LPCWSTR, newpassword LPCWSTR) (status NET_API_STATUS) = netapi32.NetUserChangePassword
//NetUserDel(servername LPCWSTR, username LPCWSTR) (status NET_API_STATUS) = netapi32.NetUserDel
//sys NetUserEnum(servername *WSTR, level DWORD, filter DWORD, bufptr **BYTE, prefmaxlen DWORD, entriesread *WORD, totalentries *DWORD, resume_handle *DWORD) (status NET_API_STATUS) = netapi32.NetUserEnum
//NetUserGetGroups(servername LPCWSTR, username LPCWSTR, level DWORD, bufptr *LPBYTE, prefmaxlen DWORD, entriesread LPDWORD, totalentries LPDWORD) (status NET_API_STATUS) = netapi32.NetUserGetGroups
//NetUserSetGroups(servername LPCWSTR, username LPCWSTR, level DWORD, buf LPBYTE, num_entries DWORD) (status NET_API_STATUS) = netapi32.NetUserSetGroups
//NetUserSetInfo(servername LPCWSTR, username LPCWSTR, level DWORD, buf LPBYTE, parm_err LPDWORD) (status NET_API_STATUS) = netapi32.NetUserSetInfo
//sys NetUserGetInfo(servername *WSTR, username *WSTR, level DWORD, bufptr **BYTE) (status NET_API_STATUS) = netapi32.NetUserGetInfo
