package common

import "fmt"

type ServerError struct {
	Code    RETURN_CODE
	Message string
}

func (this ServerError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", this.Code, this.Message)
}

func (this ServerError) GetCode() RETURN_CODE {
	return this.Code
}

func (this ServerError) GetMessage() string {
	return this.Message
}
