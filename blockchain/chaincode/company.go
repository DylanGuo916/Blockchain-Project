package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
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
		"1992-01-01",
		"1000",
		"A",
	}

	flag := PutP(stub, p)
	if flag != true {
		fmt.Println("写入信息失败")
	}

	return shim.Success(nil)
}

func (s *CompanyChaincode) saveInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var np Company
	err := json.Unmarshal([]byte(args[0]), &np)
	if err != nil {
		return shim.Error("Unmarshal np failed")
	}

	flag := PutP(stub, np)
	if !flag {
		return shim.Error("Add data failed")
	}

	return shim.Success([]byte("Add NP succeed"))
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
