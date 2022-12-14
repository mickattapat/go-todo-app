package main

import (
	"golang-todo-app-atp/apis"
	"golang-todo-app-atp/setting"

	"github.com/sirupsen/logrus"
)

func main() {
	err := setting.InitTimeZone()
	if err != nil {
		logrus.Errorln(err)
	}

	db, err := setting.InitDatabase()
	if err != nil {
		logrus.Errorln(err)
	}
	apis.InitWebApi(db)
}
