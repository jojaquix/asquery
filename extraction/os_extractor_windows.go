// +build windows

package extraction

import (
	"fmt"
	"io"

	"github.com/StackExchange/wmi"
)

// windows OsExtractor implementation
type wOsExtractor struct {
	count int8
}

func NewOsExtractor() OsExtractor {
	return &wOsExtractor{}
}

func (wOsExtractor) getOsVersion() (*OsVersionInfo, error) {

	type wmiResult struct {
		Caption string
		Version string
		Name    string
	}

	var dst []wmiResult
	err := wmi.Query("SELECT Name, Version, Caption FROM Win32_OperatingSystem", &dst)
	if err != nil {
		return nil, fmt.Errorf("wmi_win_os_extractor error")
	}

	return &OsVersionInfo{Name: dst[0].Caption, Version: dst[0].Version}, nil
}

func (osExtractor *wOsExtractor) Next() (*OsVersionInfo, error) {
	osVersionInfo, err := osExtractor.getOsVersion()

	if err != nil || osExtractor.count > 0 {
		err = io.EOF
		return nil, err
	}

	osExtractor.count++
	return osVersionInfo, nil
}

func (*wOsExtractor) ForEach(func(*OsVersionInfo) error) error {
	//not implemented yet
	return nil
}
