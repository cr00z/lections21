package domain

import "errors"

var (
	ErrorIncorrectUsername    = errors.New("incorrect username")
	ErrorUnknownUser          = errors.New("unknown user")
	ErrorEmptyAuthHeader      = errors.New("empty auth header")
	ErrorInvalidAuthHeader    = errors.New("invalid auth header")
	ErrorTokenIsEmpty         = errors.New("token is empty")
	ErrorInvalidSigningMethod = errors.New("invalid signing method")
	ErrorInvalidTokenClaims   = errors.New("invalid token claims")
	ErrorUserIdNotFound       = errors.New("user id not found")
	ErrorUserIdInvalidType    = errors.New("user id invalid type")
)
