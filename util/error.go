package util

import "fmt"

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func NewErrorMessage(code int, message string, err error) error {
	if message == "" {
		return fmt.Errorf("%d_%v", code, err)
	} else if err == nil {
		return fmt.Errorf("%d_%s", code, message)
	}

	return fmt.Errorf("%d_%s_%v", code, message, err)
}
