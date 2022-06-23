package adon

import "errors"

var (
	ErrInvalidFunctionArguments = errors.New("arguments are not match with function signature")
	ErrInvalidValueKind         = errors.New("value of kind is invalid")
)
