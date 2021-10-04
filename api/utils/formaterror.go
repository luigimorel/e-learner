package utils

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "email") {
		return errors.New("email already in use") 
	}

	if strings.Contains(err, "title"){
		return errors.New("title already in use ")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("incorrect password")
	}
	return errors.New("incorrect details")
}