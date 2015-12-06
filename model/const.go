package model

import (
	"errors"
)

import ()

var (
	NO_SUCH_RECORD = errors.New("no such record")

	NO_USER_RELATED = errors.New("no user related")
)

const (
	USER_AUTHORITY_NORMAL = 0
	USER_AUTHORITY_ADMIN  = 1

	USER_STATUS_NORMAL = 0

	POST_STATUS_NORMAL = 0

	PAGE_TYPE_DEFAULT = 0
	PAGE_TYPE_NAVBAR = 1
)
