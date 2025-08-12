package storage

import "errors"


var (
	UrlNotFound = errors.New("url not found")
	ErrUrlExists = errors.New("url exists")
)

