package main

import (
	"blc-demo/web"
	"blc-demo/web/cliInit"
	"blc-demo/web/controller"
	"blc-demo/web/dao"
)

func main() {
	//Web
	dao.InitMysql()

	app := controller.Application{
		cliInit.CliInit(),
	}

	//defer cliInit.SDK.Close()

	web.WebStart(&app)
}