// +build windows

package extraction

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpandRegistryGlobs(t *testing.T) {
	assert := assert.New(t)

	userProgramKeys := make([]string, 0, 10)
	userProgramKeys = expandRegistryGlobs("HKEY_USERS\\%\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
		userProgramKeys)
	assert.NotEqual(0, len(userProgramKeys))

	for _, v := range userProgramKeys {
		t.Log(v)
	}

}

func TestPrograms(t *testing.T) {
	assert := assert.New(t)

	data := GetPrograms()

	assert.NotEqual(0, data.Len())

	t.Log("Listing Programs")

	for e := data.Front(); e != nil; e = e.Next() {
		v := e.Value.(Row)
		t.Log(v["name"], "|", v["version"], "|", v["publisher"], "|", v["install_source"], "|", v["identifying_number"], "|")
	}
	//for i, v := range users {
	//	t.Log(i, "Name: ", v["name"], " Uuid: ", v["uuid"])
	//}

}
