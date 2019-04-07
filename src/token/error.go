package token

const (
	ERROR_TOKEN_EXPIRED           = 1000
	ERROR_TOKEN_INVALIDE          = 1001
	ERROR_TOKEN_PERMISSION_DENIED = 1002
)

var errCodes map[int16]string

func init() {
	errCodes = map[int16]string{
		ERROR_TOKEN_EXPIRED:           "error_token_expired",
		ERROR_TOKEN_INVALIDE:          "error_token_invalide",
		ERROR_TOKEN_PERMISSION_DENIED: "error_token_permission_denied",
	}
}

func getErrorMsg(id int16) string {
	return errCodes[id]
}

func NewError(id int16) Error {
	return Error{id, getErrorMsg(id)}
}
