package models

const (
	ErrorInvalidToken ErrorCode = 1000
)

type ErrorCode int

type Error struct {
	Error string    `json:"error"`
	Code  ErrorCode `json:"code"`
}
