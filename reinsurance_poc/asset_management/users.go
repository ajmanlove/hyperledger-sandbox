package main

import "errors"

type UserManager struct{}

func (a *UserManager) init() error {
	return errors.New("init() not implemented")
}
