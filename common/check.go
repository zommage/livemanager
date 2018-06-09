package common

import (
	"fmt"
	"regexp"
)

// 以数字和字母开头,包含下划线和扛
func NumLetterLine(str string) error {
	if str == "" {
		return fmt.Errorf("It is nil")
	}

	if len(str) > 64 {
		return fmt.Errorf("The length of str is invalid, less 64")
	}

	regexpStr := "^[a-zA-Z0-9][a-zA-Z0-9_-]*$"

	regCom, err := regexp.Compile(regexpStr)
	if err != nil {
		return fmt.Errorf("expression of regexp=%v is err: %v", regexpStr, err)
	}

	matchFlag := regCom.MatchString(str)
	if !matchFlag {
		return fmt.Errorf("only start number or char, contain number,char,line")
	}

	return nil
}
