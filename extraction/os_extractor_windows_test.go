// +build windows

package extraction

import (
	"testing"

	"github.com/StackExchange/wmi"
)

type Win32_Process struct {
	Name string
}

func TestWmiGetOsInfo(t *testing.T) {
	type Result struct {
		Caption string
		Version string
		Name    string
	}

	var dst []Result
	err := wmi.Query("SELECT Name, Version, Caption FROM Win32_OperatingSystem", &dst)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Listing results")
	for i, v := range dst {
		t.Log(i, v.Name, ";", v.Caption, ";", v.Version)
	}

}
func TestWmiList(t *testing.T) {
	//assert := assert.New(t)

	var dst []Win32_Process
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)

	if err != nil {
		t.Fatal(err)
	}

	//t.Log("Listing results")
	//for i, v := range dst {
	//	t.Log(i, v.Name)
	//}

}
