
package sdkInit

import "github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

type InitInfo struct {
	ChannelID      string // 通道 id
	ChannelConfig  string // 通道配置文件
	OrgName        string // 组织名称
	OrgAdmin       string //已组织管理员身份创建
	OrdererName string // order
	Peer           string
	OrgResMgmt     *resmgmt.Client

	ChaincodeID     string
	ChaincodeGoPath string
	ChaincodePath   string
	UserName        string
}
