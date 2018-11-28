// +build windows

package extraction

import (
	"container/list"
	//"syscall"
	//"golang.org/x/text/encoding/unicode"
	//"golang.org/x/sys/windows"
	//"golang.org/x/sys/windows/registry"
)

func GetPrograms() list.List {
	results := list.New()
	processed := make([]string, 0, 20)

	generateDummie(&processed, results)
	genPrograms(&processed, results)
	return *results
}

var kProgramKeys = [...]string{
	"HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
	"HKEY_LOCAL_MACHINE\\SOFTWARE\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
}

func genPrograms(processed *[]string, results *list.List) {

	programKeys := make([]string, 2)

	for _, v := range kProgramKeys {
		programKeys = append(programKeys, v)
	}

	userProgramKeys := make([]string, 0, 10)
	userProgramKeys = expandRegistryGlobs("HKEY_USERS\\%\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall", userProgramKeys)

	for _, v := range userProgramKeys {
		programKeys = append(programKeys, v)
	}

	for _, v := range programKeys {
		keyEnumPrograms(v, processed, results)
	}

}

func generateDummie(processed *[]string, results *list.List) {
	r := Row{}
	r["name"] = "jjanl"
	r["version"] = "123"
	results.PushBack(r)

}

func keyEnumPrograms(key string, processed *[]string, results *list.List) {

	keyResults, err := queryKey(key)
	if err != nil {
		//TODO log
		return
	}

	for k := keyResults.Front(); k != nil; k = k.Next() {
		row := k.Value.(Row)
		if row["type"].(string) != "subKey" {
			continue
		}

		fullProgramName := row["path"].(string)
		if Contains(*processed, fullProgramName) {
			continue
		}

		*processed = append(*processed, fullProgramName)
		appResults, _ := queryKey(key)

		for aKey := appResults.Front(); aKey != nil; aKey = aKey.Next() {
			aRow := aKey.Value.(Row)
			aRow = aRow

		}

		r := Row{}
		r["name"] = fullProgramName
		r["version"] = "comming soon"
		results.PushBack(r)

	}

}
