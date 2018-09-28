// +build windows

package extraction

import ()

// windows implementation
type userExtractorWindows struct {
}

func NewUserExtracor() UserExtractor {
	return &userExtractorWindows{}
}

func (userExtractorWindows) GetUsers() ([]UserInfo, error) {

	return nil, nil
}
