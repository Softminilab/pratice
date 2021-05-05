package main

import (
	"fmt"
	controller2 "go_code/chapter_1/controller"
	model2 "go_code/chapter_1/model"
	"log"
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("configs.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Can't get configs file")
	}

	if viper.GetBool("debug") {
		log.Println("Services RUN on Debug model")
	}
}

func getdsn() string {
	//serAddr := viper.GetString("server.addr")

	//redisAddr := viper.GetString("redis.addr")
	//redisPwd := viper.GetString("redis.password")

	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	//val.Add("parseTime", "1")
	val.Add("loc", "Asia/Shanghai")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	return dsn
}

func main() {
	err := model2.InitDB(getdsn())
	if err != nil {
		fmt.Sprintf("origin error: %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Sprintf("stack track: %+v\n", err)
	}

	defer func() {
		err := model2.Close()
		if err != nil {
			fmt.Sprintf("origin error: %T %v\n", errors.Cause(err), errors.Cause(err))
			fmt.Sprintf("stack track: %+v\n", err)
		}
	}()

	userCtl := new(controller2.UserController)
	userCtl.QueryUsers()
}
