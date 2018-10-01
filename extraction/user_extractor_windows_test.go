// +build windows

package extraction

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessLocalAccounts(t *testing.T) {
	assert := assert.New(t)

	users, err := GetUsers()

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(0, users.Len())

	t.Log("Listing accounts")

	for e := users.Front(); e != nil; e = e.Next() {
		v := e.Value.(Row)
		t.Log("Name: ", v["username"], " Uuid: ", v["uuid"])
	}
	//for i, v := range users {
	//	t.Log(i, "Name: ", v["name"], " Uuid: ", v["uuid"])
	//}

}
