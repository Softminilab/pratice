package controller

import (
	"fmt"
	"github.com/pkg/errors"
	"go_code/model"
)

type UserController struct {}

func (u *UserController) QueryUsers()  {
	user, err := model.QueryUser()
	if err != nil {
		fmt.Printf("origin eror: %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack track: \n%+v\n", err)
		return
	}
	fmt.Print(user)
}