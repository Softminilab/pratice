package model

import (
	"github.com/pkg/errors"
)

type User struct {
	ID         uint64 `borm:"id" json:"id" name:"id"`                                                  // 用户表 id
	Name       string `borm:"name" json:"name" name:"name"`                                              // 用户名称
	Email      string `borm:"email" json:"email" name:"email"`                           // 用户email
}

//func QueryUser() ([]User, error) {
//	var users []User
//	t := b.Table(db, "tbl_users").Debug()
//	_, err := t.Select(&users, b.Where("1=1"))
//	if err != nil {
//		return nil, errors.Wrap(err, "query user is failed")
//	}
//	return users, err
//}

func QueryUser() (User, error) {
	user := User{}
	sql := "select * from tbl_users"
	err := db.QueryRow(sql).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return user, errors.Wrap(err, "model: query user is failed")
	}
	return user, err
}