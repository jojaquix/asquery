// +build windows

package extraction

import (
	"container/list"
	"regexp"
	//"syscall"
	//"golang.org/x/text/encoding/unicode"
	//"golang.org/x/sys/windows"
	//"golang.org/x/sys/windows/registry"
)

func GetPrograms() list.List {
	results := list.New()
	processed := make([]string, 0, 20)

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

		// Query additional information about the program
		appResults, _ := queryKey(fullProgramName)
		r := Row{}
		//default values
		r["name"] = ""
		r["version"] = ""
		r["install_location"] = ""
		r["install_source"] = ""
		r["language"] = ""
		r["publisher"] = ""
		r["uninstall_string"] = ""
		r["install_date"] = ""
		r["identifying_number"] = ""

		// Attempt to derive the program identifying GUID

		expression := regexp.MustCompile(`({[a-fA-F0-9]+-[a-fA-F0-9]+-[a-fA-F0-9]+-[a-fA-F0-9]+-[a-fA-F0-9]+})$`)

		identifyingNumber := expression.FindAllString(fullProgramName, -1)

		if len(identifyingNumber) > 0 {
			r["identifying_number"] = identifyingNumber[0]
		}

		for aKey := appResults.Front(); aKey != nil; aKey = aKey.Next() {
			aRow := aKey.Value.(Row)
			name := aRow["name"].(string)

			if len(identifyingNumber) == 0 && name == "BundleIdentifier" {
				r["identifying_number"] = aRow["data"].(string)
			}
			if name == "DisplayName" {
				r["name"] = aRow["data"].(string)
			}
			if name == "DisplayVersion" {
				r["version"] = aRow["data"].(string)
			}
			if name == "InstallLocation" {
				r["install_location"] = aRow["data"].(string)
			}
			if name == "InstallSource" {
				r["install_source"] = aRow["data"].(string)
			}
			if name == "Language" {
				r["language"] = aRow["data"].(string)
			}
			if name == "Publisher" {
				r["publisher"] = aRow["data"].(string)
			}
			if name == "UninstallString" {
				r["uninstall_string"] = aRow["data"].(string)
			}
			if name == "InstallDate" {
				r["install_date"] = aRow["data"].(string)
			}
		}

		results.PushBack(r)
	}

}
