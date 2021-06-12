package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

// 初始化账本
func (t *ServiceSetup) InitLedger() error {

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "initLedger", Args: [][]byte{}}
	_, err := t.ChannelClient.Execute(req)
	if err != nil {
		return err
	}
	return nil
}
// 调用链码向账本添加 Data, 返回一个　TX id
func (t *ServiceSetup) SaveData(c Company) (string, error) {

	company, _ := json.Marshal(c)

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "saveData", Args: [][]byte{[]byte(company)}}


	// the proposal responses from peer(s)

	response, err := t.ChannelClient.Execute(req)
	if err != nil {
		fmt.Println("err2---->")
		return "", err
	}


	return string(response.TransactionID), nil
}

// 修改 Data 信息
func (t *ServiceSetup) ModifyData(DataID, nextOwner string) (string, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "DataExchange", Args: [][]byte{[]byte(DataID), []byte(nextOwner)}}
	respone, err := t.ChannelClient.Execute(req)

	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

// 通过 DataID 查询
func (t *ServiceSetup) FindDataByID(ID string) ([]byte, error) {

	req := channel.Request{
		ChaincodeID: t.ChaincodeId,
		Fcn: "queryInfoById",
		Args: [][]byte{[]byte(ID)},
	}

	respone, err := t.ChannelClient.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 通过 txId 查询交易
func (t *ServiceSetup) QueryTransactionByTxID(txId string) ([]byte, error) {
	req, err := t.LedgerClient.QueryTransaction(fab.TransactionID(txId))
	if err != nil {
		return []byte{0x00}, err
	}
	return req.TransactionEnvelope.Payload, nil
}

// 通过 txId 查询区块
func (t *ServiceSetup) QueryBlockByTxID(txId string) (*common.Block, error) {
	req, err := t.LedgerClient.QueryBlockByTxID(fab.TransactionID(txId))
	if err != nil {
		return nil, err
	}
	return req, nil
}

// 通过 区块高度 查询区块
func (t *ServiceSetup) QueryBlockByNum(num uint64) (*common.Block, error) {
	req, err := t.LedgerClient.QueryBlock(num)
	if err != nil {
		return nil, err
	}
	return req, nil
}


