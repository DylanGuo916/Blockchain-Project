package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type CompanyChaincode struct {
}



func main() {
	// Create a new  Smart Contract
	err := shim.Start(new(CompanyChaincode))
	if err != nil {
		fmt.Printf("Error starting Oil chaincode: %s", err)
	}
}

// 实现 Init 方法, 实例化账本时使用。
func (s *CompanyChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	s.initLedger(stub)
	return shim.Success(nil)
}

// 实现 Invoke 方法
func (s *CompanyChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// 获取函数名称、参数
	fn, args := stub.GetFunctionAndParameters()

	//调用对应函数
	if fn == "saveInfo" {
		return s.saveInfo(stub, args)

	} else if fn == "queryInfoById" {
		return s.queryInfoById(stub, args)

	}

	return shim.Error("Invalid Smart Contract function name.")
}









//peer chaincode query -C rank -n rank -c '{"Args":["queryInfoById","452145165233223223"]}'
