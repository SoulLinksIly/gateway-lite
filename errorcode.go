package main

import (
	"errors"
)

// error message
var (
	// ErrorS2S string = "s2s is not supported"
	// ErrorHandshake string = "handshake error"
	ErrorS2S           = errors.New("s2s is not supported")
	ErrorHandshake     = errors.New("handshake error")
	ErrorNewUser       = errors.New("new user")
	ErrorLogin         = errors.New("user or password not correct")
	ErrNodeTooShort    = errors.New("node name too short, it must >=8 charaters")
	ErrReadDB          = errors.New("read db error")
	ErrUserOrPwdNotSet = errors.New("user or password not set")
)
