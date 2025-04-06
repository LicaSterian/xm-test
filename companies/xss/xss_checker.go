package xss

import (
	"errors"

	"github.com/microcosm-cc/bluemonday"
)

var ErrFoundXSS = errors.New("found XSS in input")

func CheckForXSS(input string) error {
	p := bluemonday.UGCPolicy()

	sanitizedInput := p.Sanitize(input)
	if sanitizedInput != input {
		return ErrFoundXSS
	}
	return nil
}
