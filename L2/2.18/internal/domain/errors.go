package domain

import "errors"

var (
	ErrEventNotFound  = errors.New("event not found")
	ErrInvalidDate    = errors.New("invalid date")
	ErrInvalidUserID  = errors.New("invalid user id")
	ErrInvalidEventID = errors.New("invalid event ID")
	ErrInvalidText    = errors.New("invalid event text")
)
