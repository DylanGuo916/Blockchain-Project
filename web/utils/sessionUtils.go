package utils

import (
	"blc-demo/web/dao"
	"blc-demo/web/model"
	"blc-demo/web/service"
	"fmt"
	"net/http"
)

func DeleteSession(r *http.Request) *struct {
	Sess         *model.Session
	FailedLogin  bool
	IsLogin      bool
	IsSuperAdmin bool
	IsAdmin      bool
	IsUser       bool
	IsStaff      bool
	Msg          string
} {
	data := &struct {
		Sess         *model.Session
		FailedLogin  bool
		IsLogin      bool
		IsSuperAdmin bool
		IsAdmin      bool
		IsUser       bool
		IsStaff      bool
		Msg          string
	}{
		Sess:         nil,
		FailedLogin:  false,
		IsLogin:      false,
		IsSuperAdmin: false,
		IsAdmin:      false,
		IsUser:       false,
		IsStaff:      false,
		Msg:          "",
	}
	//获取cookie
	cookie, _ := r.Cookie("user")
	if cookie != nil {
		//获取cookie的值
		cookieValue := cookie.Value
		//删除session
		dao.DeleteSession(cookieValue)
	}
	fmt.Println("---------------------------------------------")
	fmt.Println("Session已删除，正在退出登录")
	return data
}

func CheckLogin(r *http.Request) *struct {
	Sess         *model.Session
	FailedLogin  bool
	IsLogin      bool
	IsSuperAdmin bool
	IsAdmin      bool
	IsUser       bool
	IsStaff      bool
	Msg          string
	Admin        []*model.User
	User         []*model.User
	Staff        []*model.User
	Action       []*model.Action
	Data          service.Company
	CodePath  string
	Score float64
	Rank string
} {

	fmt.Println("---------------------------------------------")
	fmt.Println("默认参数已就绪")
	data := &struct {
		Sess         *model.Session
		FailedLogin  bool
		IsLogin      bool
		IsSuperAdmin bool
		IsAdmin      bool
		IsUser       bool
		IsStaff      bool
		Msg          string
		Admin        []*model.User
		User         []*model.User
		Staff        []*model.User
		Action       []*model.Action
		Data          service.Company
		CodePath  string
		Score float64
		Rank string
	}{
		Sess:         nil,
		FailedLogin:  false,
		IsLogin:      false,
		IsSuperAdmin: false,
		IsAdmin:      false,
		IsUser:       false,
		IsStaff:      false,
		Msg:          "",
		Admin:        nil,
		User:         nil,
		Staff:        nil,
		Action:       nil,
		Data:          service.Company{},
		CodePath:  "",
	}

	//获取cookie
	cookie, _ := r.Cookie("user")
	if cookie != nil {
		//获取cookie的值
		cookieValue := cookie.Value
		//在数据库查询cookieValue对应的session
		session := dao.GetSession(cookieValue)

		if session.UserID > 0 {
			fmt.Println("用户已登录")
			if session.Role == "超级管理员" {
				data.IsLogin = true
				data.IsSuperAdmin = true
			} else if session.Role == "管理员" {
				data.IsLogin = true
				data.IsAdmin = true
			} else if session.Role == "用户" {
				data.IsLogin = true
				data.IsUser = true
			} else if session.Role == "员工" {
				data.IsLogin = true
				data.IsStaff = true
			}
			data.Sess = session
			return data
		} else {
			fmt.Println("用户未登录")
			data.IsLogin = false
			data.Sess = nil
			DeleteSession(r)
			return data
		}
	}
	DeleteSession(r)
	return data
}
