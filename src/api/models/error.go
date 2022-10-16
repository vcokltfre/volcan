package models

/*
Error categories:
- 1000: Generic errors
- 11xx: Authentication errors
*/

const (
	ErrorGeneric      ErrorCode = 1000
	ErrorInvalidToken ErrorCode = 1100
)

type ErrorCode int

type Error struct {
	Error string    `json:"error"`
	Code  ErrorCode `json:"code"`
}
