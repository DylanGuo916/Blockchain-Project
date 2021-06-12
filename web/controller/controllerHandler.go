package controller

import (
	"blc-demo/web/service"
	"blc-demo/web/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
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
	var score float64
	data.Data = d
	for i := 0; i < len(d.Score); i++ {
		score += d.Score[i]
	}
	data.Score, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", score/float64(len(d.Score))), 64)
	if score >= 95.0 {
		data.Rank = "A+"
	} else if score >= 90.0 {
		data.Rank = "A"
	} else if score >= 85.0 {
		data.Rank = "A-"
	} else if score >= 80.0 {
		data.Rank = "B+"
	} else if score >= 75.0 {
		data.Rank = "B"
	} else if score >= 75.0 {
		data.Rank = "B-"
	}
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
func SwitchTimeStampToData(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)
	return t.Format("2006-01-02 15:04:05")
}
func (app *Application) AddData(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsStaff {
		//r.ParseMultipartForm(32 << 10)
		//获取表单输入
		id := r.FormValue("id")
		name := r.FormValue("company_name")
		legal := r.FormValue("legal")
		score, _ := strconv.ParseFloat(r.FormValue("score"), 64)
		rank := r.FormValue("rank")

		//defer content.Close()

		d := service.Company{
			id,
			name,
			legal,
			SwitchTimeStampToData(time.Now().Unix()),
			[]float64{score},
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
