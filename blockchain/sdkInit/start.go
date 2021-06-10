package sdkInit

import (
	"fmt"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

const (
	ChaincodeVersion = "1.0"
	lvl              = logging.INFO
)

func SetupLogLevel() {
	logging.SetLevel("fabsdk", lvl)
	logging.SetLevel("fabsdk/common", lvl)
	logging.SetLevel("fabsdk/fab", lvl)
	logging.SetLevel("fabsdk/client", lvl)
}

// 实例化 SDK
func SetupSDK(ConfigFile string, initialized bool) (*fabsdk.FabricSDK, error) {

	if initialized {
		return nil, fmt.Errorf("Fabric SDK 已被实例化")
	}

	sdk, err := fabsdk.New(config.FromFile(ConfigFile))
	if err != nil {
		return nil, fmt.Errorf("实例化Fabric SDK失败: %v", err)
	}

	fmt.Println("Fabric SDK 初始化成功")
	return sdk, nil
}

//--------------------------------------------------------------------
// 创建通道
//--------------------------------------------------------------------
func CreateChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {
	// channel is exists?
	// 1. 生成客户端上下文环境， 什么身份--> 组织管理员（哪个组织）
	clientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	if clientContext == nil {
		return fmt.Errorf("根据指定的组织管理员创建户端Context失败")
	}

	// 2. 根据上下文环境，创建 resMgmtClient, 用来通道的创建，链码的安装、实例化和升级等
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败: %v", err)
	}

	info.OrgResMgmt = resMgmtClient

	// 3. mspClient 与证书有关的客户端
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(info.OrgName))
	if err != nil {
		return fmt.Errorf("根据指定的 OrgName 创建 Org MSP 客户端实例失败: %v", err)
	}

	adminIdentity, err := mspClient.GetSigningIdentity(info.OrgAdmin)
	if err != nil {
		return fmt.Errorf("获取指定id的签名身份失败: %v", err)
	}

	// 生成创建通道请求
	channelReq := resmgmt.SaveChannelRequest{
		ChannelID:         info.ChannelID,
		ChannelConfigPath: info.ChannelConfig,
		SigningIdentities: []msp.SigningIdentity{adminIdentity},
	}
	// RC创建通道
	_, err = resMgmtClient.SaveChannel(channelReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(info.OrdererName),
	)
	if err != nil {
		return errors.Errorf("创建应用通道失败: %v", err)
	}
	fmt.Printf("成功创建通道\n")
	return nil
}

//--------------------------------------------------------------------
// 加入通道
//--------------------------------------------------------------------
func JoinChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {

	err := info.OrgResMgmt.JoinChannel(
		info.ChannelID,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(info.OrdererName),
	)
	if err != nil {
		return fmt.Errorf("Peers加入通道失败: %v", err)
	}

	fmt.Println("peers 已成功加入通道.")

	return nil
}
//---------------------------------------------------------------------------------------
// 安装链码
//--------------------------------------------------------------------------------------
func InstallCC(sdk *fabsdk.FabricSDK, info *InitInfo) error {

	// 创建一个链码包
	ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
	if err != nil {
		return fmt.Errorf("创建链码包失败: %v", err)
	}

	// 生成安装链码请求
	installCCReq := resmgmt.InstallCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: ChaincodeVersion,
		Package: ccPkg,
	}

	// 安装链码
	_, err = info.OrgResMgmt.InstallCC(
		installCCReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)

	if err != nil {
		return errors.Errorf("安装链码失败: %v", err)
	}

	fmt.Printf("安装链码 %v 成功\n", info.ChaincodeID)
	return nil
}

//---------------------------------------------------------------------------------------
// 查询已安装链码
//--------------------------------------------------------------------------------------
func QueryInstalledCC(sdk *fabsdk.FabricSDK, info *InitInfo) {

	resp2, err := info.OrgResMgmt.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(info.Peer))
	if err != nil {
		fmt.Println("查询已安装的链码失败: ", err)
	}

	fmt.Println("已安装链码包括: ", resp2.GetChaincodes())
}

//---------------------------------------------------------------------------------------
// 实例化链码
//--------------------------------------------------------------------------------------
func InstantiateCC(sdk *fabsdk.FabricSDK, info *InitInfo) error {

	ccPolicy, _ := cauthdsl.FromString("AND ('Org1MSP.member')")
	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: ChaincodeVersion,
		Args:    [][]byte{[]byte("init")},
		Policy:  ccPolicy,

	}

	_, err := info.OrgResMgmt.InstantiateCC(
		info.ChannelID,
		instantiateCCReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)

	if err != nil {
		return errors.Errorf("实例化链码失败: %v\n", err)
	}

	fmt.Printf("实例化 %v 成功\n", info.ChaincodeID)
	return nil
}
