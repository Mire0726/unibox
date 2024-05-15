package cerror

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Error struct {
	error
	code        Code
	clientMsg   string
	reasonCodes []ReasonCode
	stack       *stack
}

type errorResponse struct {
	Code        Code     `json:"code"`
	Detail      string   `json:"detail"`
	ReasonCodes []string `json:"reasonCodes,omitempty"`
	Stack       string   `json:"-"`
	Status      int      `json:"-"`
}

// CustomHTTPErrorHandler : the custom error handlers for echo
func CustomHTTPErrorHandler(err error, c echo.Context) {
	c.Response().Header().Set(echo.HeaderContentType, "application/problem+json")
	errorResponse := errorFormatter(err)

	if !c.Response().Committed {
		err = c.JSON(errorResponse.Status, errorResponse)
		if err != nil {
			panic(fmt.Errorf("JSONのエンコードに失敗: %w", err))
		}
	}
}

func errorFormatter(err error) errorResponse {
	resp := errorResponse{
		Code:   Unknown,
		Detail: err.Error(),
		Status: http.StatusInternalServerError,
	}

	// カスタムエラー`Error`型の場合、詳細情報を取得
	var customErr *Error
	if errors.As(err, &customErr) {
		resp.Code = customErr.Code()
		resp.Detail = customErr.ClientMsg()
		resp.Status = GetHTTPStatus(err)
		resp.ReasonCodes = customErr.ReasonCodes()
		resp.Stack = customErr.StackTrace()
	} else if he, ok := err.(*echo.HTTPError); ok { //nolint:errorlint
		// EchoのHTTPエラーの場合
		resp.Status = he.Code
		resp.Detail = fmt.Sprintf("%v", he.Message)
		resp.Code = MapHTTPErrorToCode(he.Code)
	}

	return resp
}

func (e *Error) Error() string {
	return e.error.Error()
}

func (e *Error) Code() Code {
	return e.code
}

func (e *Error) ClientMsg() string {
	return e.clientMsg
}

func (e *Error) ReasonCodes() []string {
	if len(e.reasonCodes) == 0 {
		// NOTE:
		// if reasonCode is empty, return default reason code
		// I want to alert that reason code is not set
		return []string{string(RC00000)}
	}

	reasonCodes := make([]string, 0, len(e.reasonCodes))
	for _, rc := range e.reasonCodes {
		reasonCodes = append(reasonCodes, string(rc))
	}

	return reasonCodes
}

func (e *Error) StackTrace() string {
	if e.stack == nil {
		return ""
	}

	return fmt.Sprintf("%+v", e.stack)
}

func (e *Error) Unwrap() error {
	return e.error
}

// GetHttpStatusFromErrCode returns http status code.
func GetHTTPStatusFromErrCode(c Code) int {
	v, ok := codeMap[c]
	if ok {
		return v.httpStatus
	}

	return http.StatusInternalServerError
}

// GetHTTPStatus returns http status code.
func GetHTTPStatus(err error) int {
	c := GetCode(err)

	return GetHTTPStatusFromErrCode(c)
}

// GetCode takes error and returns Code accordingly.
func GetCode(err error) Code {
	if err == nil {
		return OK
	}
	var e *Error
	if errors.As(err, &e) {
		return e.code
	}

	return Unknown
}

// New returns cerror with given message
// if you want to add params to msg, please use `fmt.Sprintf()`
func New(msg string, opts ...Option) error {
	err := &Error{
		error: errors.New(msg), //nolint:goerr113
		code:  Unknown,
		stack: callers(),
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

func Wrap(err error, prefix string, opts ...Option) error {
	werr := &Error{
		error: fmt.Errorf("%s: %w", prefix, err),
	}

	if cerr, ok := As(err); ok {
		werr.stack = cerr.stack
		werr.code = cerr.code
		werr.clientMsg = cerr.clientMsg
		werr.reasonCodes = cerr.reasonCodes
	} else {
		werr.stack = callers()
		werr.code = Unknown
	}

	for _, opt := range opts {
		opt(werr)
	}

	return werr
}

func As(err error) (*Error, bool) {
	cerr := &Error{}
	if ok := errors.As(err, &cerr); !ok {
		return nil, false
	}

	return cerr, true
}

// Is checks if error is of Code c.
func Is(err error, c Code) bool {
	return GetCode(err) == c
}

func IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
