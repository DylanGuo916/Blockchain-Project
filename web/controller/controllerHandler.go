package controller

import (
	"blc-demo/web/service"
	"blc-demo/web/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Application struct {
	Setup *service.ServiceSetup
}

// 进入查询页面
func (app *Application) QueryPage(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	ShowView(w, r, "PublicOption/queryPage.html", data)
}

// 根据ID查询信息
func (app *Application) FindDataByID(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	ID := r.FormValue("id")
	result, err := app.Setup.FindDataByID(ID)
	if err != nil {
		log.Println(err)
	}
	var d = service.Company{}
	err = json.Unmarshal(result, &d)
	fmt.Println(d)
	if err != nil {
		log.Println("unmarshal failed, err:", err)
	}

	data.Data = d

	ShowView(w, r, "PublicOption/queryResult.html", data)
}

func (app *Application) AddDataPage(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	if data.IsLogin {
		if data.IsStaff {
			ShowView(w, r, "StaffOption/addDataPage.html", data)
			return
		} else {
			data.Msg = "无权访问"
			ShowView(w, r, "index.html", data)
			return
		}
	} else if !data.IsLogin {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}
func (app *Application) AddData(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsStaff {
		//r.ParseMultipartForm(32 << 10)
		//获取表单输入
		id := r.FormValue("id")
		name := r.FormValue("company_name")
		legal := r.FormValue("legal")
		date := r.FormValue("date")
		score := r.FormValue("score")
		rank := r.FormValue("rank")

		//defer content.Close()

		d := service.Company{
			id,
			name,
			legal,
			date,
			score,
			rank,
		}

		_, err := app.Setup.SaveData(d)
		if err != nil {
			fmt.Println("err3:", err)
		}
		data.Data = d

		ShowView(w, r, "StaffOption/addSuccess.html", data)
		return
	} else if !data.IsStaff {
		ShowView(w, r, "index.html", data)
		return
	}
}
