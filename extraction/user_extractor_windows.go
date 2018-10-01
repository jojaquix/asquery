// +build windows

package extraction

import (
	"container/list"
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
	processLocalAccounts(processedSids, &results)

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
