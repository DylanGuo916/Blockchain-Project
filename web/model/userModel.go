package model

type User struct {
	Id				int
	Username		string
	Password		string
	//Logincounter	int64
	Role			string
	Phone			string
	Status     		string			// 0 正常状态， 1 删除
	Createtime 		string
}
