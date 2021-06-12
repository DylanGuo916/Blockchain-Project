package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"blc-demo/web/dao"
	"blc-demo/web/model"
	"blc-demo/web/utils"
	"time"
)

//进入首页
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	ShowView(w, r, "index.html", data)
}

// 返回首页
func (app *Application) BackToHome(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	ShowView(w, r, "index.html", data)
}

// 进入注册界面
func (app *Application) RegisterPage(w http.ResponseWriter, r *http.Request) {
	ShowView(w, r, "AccountRelated/register.html", nil)
}

// 随机数字6位
func GenValidateCode(width int) string {
	numeric := [10]byte{0,1,2,3,4,5,6,7,8,9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
// 注册添加用户信息
func (app *Application) Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("loginName")
	Password := r.FormValue("password")
	phone := r.FormValue("tel")
	role := "员工"
	statue := "正常"

	if username == "" || Password == "" || phone == "" {
		ShowView(w, r, "AccountRelated/register.html", nil)
		fmt.Println("---------> err 1")
		return
	} else {

		ID := dao.QueryUserWithUsername(username)
		fmt.Println("id", ID)
		if ID > 0 {
			ShowView(w, r, "AccountRelated/register.html", nil)
			fmt.Println("---------> err 2")
			return
		}
		password := utils.MD5(Password)
		createtime := utils.SwitchTimeStampToData(time.Now().Unix())

		fmt.Println(createtime)

		user := model.User{
			Username:   username,
			Password:   password,
			Role:       role,
			Phone:      phone,
			Status:     statue,
			Createtime: createtime,
		}
		_, err := dao.InsertUser(user) //user表插入记录

		if err != nil {
			ShowView(w, r, "AccountRelated/register.html", nil)
			fmt.Println("---------> err 3")
			return
		} else {
			ShowView(w, r, "AccountRelated/login.html", nil)
			return
		}

	}
}

// 进入登录界面
func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	if data.IsLogin {
		ShowView(w, r, "index.html", data)
		return
	}
	ShowView(w, r, "AccountRelated/login.html", data)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("---------------------------------------------")
	fmt.Println("默认参数已就绪")
	fmt.Println("---------------------------------------------")
	//获取表格信息
	username := r.FormValue("loginName")
	Password := r.FormValue("password")
	password := utils.MD5(Password)
	fmt.Println("---------------------------------------------")
	fmt.Println("前端表格读取完成")

	//返回完整的用户信息
	user := dao.FindUserByUsernameAndPassword(username, password)
	fmt.Println("---------------------------------------------")
	fmt.Println("用户", user.Username, "查询结果已传回，正在核查")

	if user.Id == 0 {
		data.FailedLogin = true
		fmt.Println("---------------------------------------------")
		fmt.Println("用户名或密码错误，登陆失败，以未登录状态返回首页")
		data.Msg = "用户名或密码错误"
		ShowView(w, r, "AccountRelated/login.html", data)
		return

	} else if user.Status == "异常" {
		data.FailedLogin = true
		fmt.Println("---------------------------------------------")
		fmt.Println(user.Role, user.Username, "账户受限，登陆失败，以未登录状态返回首页")
		data.Msg = user.Role + user.Username + "账户受限，登陆失败，请联系管理员"
		ShowView(w, r, "index.html", data)
		return

	} else if user.Status == "正常" {
		uuid := utils.CreateUUID()
		session := &model.Session{
			SessionID:  uuid,
			UserID:     user.Id,
			UserName:   user.Username,
			PassWord:   user.Password,
			Role:       user.Role,
			Phone:      user.Phone,
			Status:     user.Status,
			CreateTime: user.Createtime,
		}

		_ = dao.AddSession(session)

		fmt.Println("---------------------------------------------")
		fmt.Println("Session已设置")

		cookie := http.Cookie{
			Name:     "user",
			Value:    uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		fmt.Println("---------------------------------------------")
		fmt.Println("Cookie已送往浏览器")
		if user.Role == "员工" {
			data.IsStaff = true
		}
		data.IsLogin = true
		data.Sess = session
		fmt.Println("---------------------------------------------")
		fmt.Println("默认参数已更新")

			ShowView(w, r, "index.html", data)

	}
}

// 退出登陆
func (app *Application) Logout(w http.ResponseWriter, r *http.Request) {
	data := utils.DeleteSession(r)
	fmt.Println(data.Msg)
	ShowView(w, r, "AccountRelated/login.html", data)
}

func (app *Application) About(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	ShowView(w, r, "PublicOption/about.html", data)
}


func (app *Application) ModifyStatus(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	if data.IsLogin {
		userID, _ := strconv.ParseInt(r.FormValue("userID"), 10, 64)
		userStatus := r.FormValue("userStatus")
		userRole := r.FormValue("userRole")

		if userStatus == "正常" {
			dao.UpdateUser(userID, "异常")
		} else if userStatus == "异常" {
			dao.UpdateUser(userID, "正常")
		}

		if userRole == "用户" {
			users, _ := dao.QueryAllUser()
			data.User = users
			if data.IsAdmin {
				ShowView(w, r, "BackStage/bsUserMana.html", data)
				return
			} else if data.IsSuperAdmin {
				ShowView(w, r, "SuperBackStage/sbsUserMana.html", data)
				return
			}
		} else if userRole == "员工" {
			fmt.Println("---------------------------------------------")
			fmt.Println("查询已有职员")
			staffs, _ := dao.QueryAllStaff()
			data.Staff = staffs
			ShowView(w, r, "BackStage/bsStaffMana.html", data)
		}
	} else {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}
