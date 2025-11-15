package dto

type BaseResponse[T any] struct {
	Result  *T      `json:"result,omitempty"`
	Message *string `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}

func ToResponse[T any](result T) *BaseResponse[T] {
	return &BaseResponse[T]{
		Result: &result,
	}
}

func ToResponseWithMessage[T any](result T, message string) *BaseResponse[T] {
	return &BaseResponse[T]{
		Result:  &result,
		Message: &message,
	}
}

func ToError(err error) *BaseResponse[any] {
	e := err.Error()
	return &BaseResponse[any]{
		Result:  nil,
		Message: nil,
		Error:   &e,
	}
}

func ToErrorWithMessage(err error, message string) *BaseResponse[any] {
	e := err.Error()
	return &BaseResponse[any]{
		Result:  nil,
		Message: &message,
		Error:   &e,
	}
}
