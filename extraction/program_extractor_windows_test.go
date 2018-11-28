// +build windows

package extraction

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrograms(t *testing.T) {
	assert := assert.New(t)

	data := GetPrograms()

	assert.NotEqual(0, data.Len())

	t.Log("Listing Programs")

	for e := data.Front(); e != nil; e = e.Next() {
		v := e.Value.(Row)
		t.Log(v["name"], v["version"])
	}
	//for i, v := range users {
	//	t.Log(i, "Name: ", v["name"], " Uuid: ", v["uuid"])
	//}

}
