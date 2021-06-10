
package web

import (
	"fmt"
	"net/http"
	"blc-demo/web/controller"
)

func WebStart(app *controller.Application)  {


	fs := http.FileServer(http.Dir("./web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// 打开系统进入的页面
	http.HandleFunc("/", app.LoginView)

	http.HandleFunc("/backToHome", app.BackToHome)			// 返回首页
	http.HandleFunc("/registerPage", app.RegisterPage)
	http.HandleFunc("/register", app.Register)

	// 登陆
	http.HandleFunc("/loginPage", app.LoginView)
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/logout", app.Logout)
	// 查询
	http.HandleFunc("/queryPage", app.QueryPage)		// 转至查询信息页面
	http.HandleFunc("/findDataByID", app.FindDataByID)	// 根据id查询并转至查询结果页面
	//
	http.HandleFunc("/addDataPage", app.AddDataPage) // 显示添加信息页面
	http.HandleFunc("/addData", app.AddData)         // 提交修改请求并跳转添加成功提示页面

	fmt.Println("---------------------------------------------")
	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}

}
