package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"time"
)

// 该文件实现使用链码相关API对账本状态进行具体操作的函数们

func PutP(stub shim.ChaincodeStubInterface, p Company) bool {

	b, err := json.Marshal(p)
	if err != nil {
		return false
	}
	err = stub.PutState(p.Id, b)
	if err != nil {
		return false
	}
	return true
}

func (s *CompanyChaincode) initLedger(stub shim.ChaincodeStubInterface) pb.Response {

	p := Company{
		"452145165233223223",
		"Wahaha",
		"宗庆后",
		SwitchTimeStampToData(time.Now().Unix()),
		[]float64{100},
		"A",
	}

	flag := PutP(stub, p)
	if flag != true {
		fmt.Println("写入信息失败")
	}

	return shim.Success(nil)
}
func SwitchTimeStampToData(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)
	return t.Format("2006-01-02 15:04:05")
}
func GetCompanyInfo(stub shim.ChaincodeStubInterface, entityId string) (Company, bool) {

	var company Company
	b, err := stub.GetState(entityId)
	if err != nil || b == nil {
		return company, false
	} // 有错误 或者 Id不存在[id不存在GetState()返回 nil, nil]

	err = json.Unmarshal(b, &company)
	if err != nil {
		return company, false
	}

	return company, true
}
func (s *CompanyChaincode) saveData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var np Company
	err := json.Unmarshal([]byte(args[0]), &np)
	if err != nil {
		return shim.Error("Unmarshal np failed")
	}
	company, exist := GetCompanyInfo(stub, np.Id)

	if exist {
		np.Score = append(company.Score,np.Score[0])
	}
	flag := PutP(stub, np)
	if !flag {
		return shim.Error("Add data failed")
	}
	return shim.Success([]byte("Add Company succeed"))
}

func (s *CompanyChaincode) queryInfoById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("incorrect nums of  args, expecting 1")
	}
	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("query failed according to ic")
	}
	if result == nil {
		return shim.Error("get nothing according to ic")
	}
	return shim.Success(result)
}
