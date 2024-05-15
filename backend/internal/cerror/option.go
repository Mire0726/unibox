package cerror

import "fmt"

type Option func(*Error)

func WithCode(code Code) Option {
	return func(e *Error) {
		e.code = code
	}
}

func WithInternalCode() Option {
	return func(e *Error) {
		e.code = Internal
	}
}

func WithInvalidArgumentCode() Option {
	return func(e *Error) {
		e.code = InvalidArgument
	}
}

func WithNotFoundCode() Option {
	return func(e *Error) {
		e.code = NotFound
	}
}

func WithAlreadyExistsCode() Option {
	return func(e *Error) {
		e.code = AlreadyExists
	}
}

func WithUnauthorizedCode() Option {
	return func(e *Error) {
		e.code = Unauthorized
	}
}

func WithPostgreSQLCode() Option {
	return func(e *Error) {
		e.code = PostgreSQL
	}
}

func WithInOpportuneTimeCode() Option {
	return func(e *Error) {
		e.code = InOpportuneTime
	}
}

func WithMailCode() Option {
	return func(e *Error) {
		e.code = Mail
	}
}

func WithNoRowsCode() Option {
	return func(e *Error) {
		e.code = NoRows
	}
}

func WithClientMsg(format string, args ...any) Option {
	return func(e *Error) {
		e.clientMsg = fmt.Sprintf(format, args...)
	}
}

func WithReasonCode(rc ReasonCode) Option {
	return func(e *Error) {
		e.reasonCodes = append(e.reasonCodes, rc)
	}
}

func WithFirebaseCode() Option {
	return func(e *Error) {
		e.code = Firebase
	}
}

func WithEncodingJSONCode() Option {
	return func(e *Error) {
		e.code = EncodingJSON
	}
}

func WithIOCode() Option {
	return func(e *Error) {
		e.code = IO
	}
}

func WithDoExternalHTTPRequestCode() Option {
	return func(e *Error) {
		e.code = DoExternalHTTPRequest
	}
}

func WithCreateExternalHTTPRequestCode() Option {
	return func(e *Error) {
		e.code = CreateExternalHTTPRequest
	}
}

func WithTimeParseCode() Option {
	return func(e *Error) {
		e.code = TimeParse
	}
}

func WithTimeLoadLocationCode() Option {
	return func(e *Error) {
		e.code = TimeLoadLocation
	}
}
