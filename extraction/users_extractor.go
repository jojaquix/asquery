package extraction

//Here interfaces for users and groups based in osquery

type UserInfo struct {
	Uid         uint64
	Gid         uint64
	Uid_signed  int64
	Gid_signed  int64
	Username    string
	Description string
	Directory   string
	Shell       string
	Uuid        string
	Ttype       string
}

//type UserExtractor interface {
//	Next() (*UserInfo, error)
//	ForEach(func(*UserInfo) error) error
//}

type UserExtractor interface {
	GetUsers() ([]UserInfo, error)
}
