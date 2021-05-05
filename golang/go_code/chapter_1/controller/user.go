package controller

import (
	"fmt"
	model2 "go_code/chapter_1/model"

	"github.com/pkg/errors"
)

type UserController struct{}

func (u *UserController) QueryUsers() {
	user, err := model2.QueryUser()
	if err != nil {
		fmt.Printf("origin eror: %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack track: \n%+v\n", err)
		return
	}
	fmt.Print(user)
}
