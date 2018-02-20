package lib

type Error struct {
	cause       error
	t           int
	msgOverride string
}

const (
	ResourceUnreachable = iota
	UnsupportedContentType
	TransformationFailure
	EncodingFailure
	InvalidParams
)

var codeMap = map[int]int{
	ResourceUnreachable:    404,
	UnsupportedContentType: 400,
	TransformationFailure:  500,
	EncodingFailure:        500,
	InvalidParams:          400,
}

var msgMap = map[int]string{
	ResourceUnreachable:    "Unable to access specified url",
	UnsupportedContentType: "Content type is not supported, supported formats: jpeg, gif, png",
	TransformationFailure:  "Sorry, but something went wrong, our support engineers are already notified",
	EncodingFailure:        "Sorry, but something went wrong, our support engineers are already notified",
	InvalidParams:          "Request params are invalid, please, verify that url is a valid url, width and height are positive integers",
}

func NewError(cause error, t int, msgOverride ...string) Error {
	msg := ""
	if len(msgOverride) > 0 {
		msg = msgOverride[0]
	}

	return Error{
		cause:       cause,
		t:           t,
		msgOverride: msg,
	}
}

func (e Error) Error() string {
	if e.msgOverride != "" {
		return e.msgOverride
	}

	if e.cause != nil {
		return e.cause.Error()
	}

	return ""
}

func (e Error) Code() int {
	code, ok := codeMap[e.t]
	if ok {
		return code
	}

	return 500
}

const GenericMsg = "Sorry, but something went wrong, our support engineers are already notified"

func (e Error) Msg() string {
	if e.msgOverride != "" {
		return e.msgOverride
	}

	msg, ok := msgMap[e.t]
	if ok {
		return msg
	}

	return GenericMsg
}
