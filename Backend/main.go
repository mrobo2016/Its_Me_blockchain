package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
)

// marshal/unmarshal to encode/decode json file
// Save Json into struct
type Config struct {
	Authorization      string `json:"authorization"`
	KasAccount         string `json:"kas_account"` // must be from KAS API Wallet (Default Account Pool)
	MasterContractList string `json:"master_contract_list"`
}

// Defined after main function.
var conf = readConfig()

func main() {
	// Gin is a web framework written in Go (Golang)
	router := gin.Default()
	// CORS is a middleware for gin
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))

	account := router.Group("/kas")
	{
		// GET-POST request
		account.GET("/account", wallet_getAccount)                         // ㅇ
		account.POST("/account", wallet_createAccount)                     // ㅇ
		account.POST("/store", wallet_deployStoreContract)                 // ㅇ JSON binding request
		account.POST("/store/newattendance", wallet_executeAddAttdendance) // ㅇ
		account.POST("/store/order", wallet_executeOrder)
		account.POST("/store/approve", wallet_executeApproveOrder)
		account.POST("/store/deny", wallet_executeDenyOrder)
		account.POST("/store/reward", wallet_executeSetReward)
		account.POST("/storelist", wallet_deployStoreListContract)    // ㅇ
		account.POST("/storelist/store", wallet_executeRegisterStore) // ㅇ
	}

	// klaytn node API
	// https://console.klaytnapi.com/ko/service/node
	klaytn := router.Group("/node")
	{
		klaytn.POST("/klaytn", node_postRPC)
		klaytn.GET("/blockNumber", node_getBlockNumber)
		klaytn.GET("/balance", node_getBalance)
		klaytn.GET("/logs", node_getLogs)
		klaytn.GET("/receipt", node_getReceipt) // ㅇ
		klaytn.GET("/stores", node_getStores)
		klaytn.GET("/menu", node_getMenu)
		klaytn.GET("/order", node_getOrder)
		klaytn.POST("/numAttendance", node_getNumAttendances)

		// klaytn.GET("/klaytn/reward", kasKlaytn_GetReward)
	}

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// reads Json file using ioutil
// uses unmarshal to decode file
func readConfig() Config {
	var conf Config
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &conf); err != nil {
		panic(err)
	}
	return conf
}

// set configs as headers for authentication
func setDefaultHeader(header http.Header) {
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", conf.Authorization)
	header.Add("x-chain-id", "1001") // select network
}

// GET kas/account
// https://refs.klaytnapi.com/en/wallet/latest#operation/RetrieveAccounts
func wallet_getAccount(c *gin.Context) {
	var getAcc FormGetAcc
	url := kasWalletUrl + "/v2/account"

	if c.ShouldBindQuery(&getAcc) != nil {
		c.String(-1, "invalid form")
		return
	}

	bodyData, _ := json.Marshal(&kasGetAcc{Address: getAcc.Id})

	ret, err := SendHTTPRequest("GET", url, bodyData)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	c.String(200, ret)
}

// POST kas/account
func wallet_createAccount(c *gin.Context) {
	url := kasWalletUrl + "/v2/account"
	c.Header("Access-Control-allow-origin", "*")

	ret, err := SendHTTPRequest("POST", url, nil)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	c.JSON(200, ret) // returns address etc
}

// POST node/klaytn
// https://refs.klaytnapi.com/en/node/latest#operation/CallJSON-RPC
func node_getBalance(c *gin.Context) {
	url := kasNodeUrl + "/v1/klaytn"
	c.Header("Access-Control-allow-origin", "*")

	var params FormAddress
	if err := c.ShouldBindQuery(&params); err != nil {
		c.String(-1, "invalid input data"+err.Error())
		return
	}

	// Using RPC - https://console.klaytnapi.com/ko/service/node
	bodyData, _ := json.Marshal(&kasKlaytnRPC{
		JsonRpc: "2.0",
		Method:  "klay_getBalance", // peb 단위 현재 잔액을 정수 형태로 반환합니다.
		Params:  []string{params.Address, "latest"},
		Id:      1,
	})

	ret, err := SendHTTPRequest("POST", url, bodyData)
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	c.JSON(200, ret)
}

// GET node/blockNumber
//
func node_getBlockNumber(c *gin.Context) {
	url := kasNodeUrl + "/v1/klaytn"
	c.Header("Access-Control-allow-origin", "*")

	bodyData, _ := json.Marshal(&kasKlaytnRPC{
		JsonRpc: "2.0",
		Method:  "klay_blockNumber",
		Params:  []struct{}{},
		Id:      1,
	})

	ret, err := SendHTTPRequest("POST", url, bodyData)
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	c.JSON(200, ret)
}

func node_postRPC(c *gin.Context) {
	url := kasNodeUrl + "/v1/klaytn"
	c.Header("Access-Control-allow-origin", "*")

	bodyData, _ := json.Marshal(&kasKlaytnRPC{
		JsonRpc: "2.0",
		Method:  "klay_blockNumber",
		Params:  []struct{}{},
		Id:      1,
	})

	ret, err := SendHTTPRequest("POST", url, bodyData)
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	c.String(200, ret)
}

//GET node/receipt
func node_getReceipt(c *gin.Context) {
	url := kasNodeUrl + "/v1/klaytn"
	c.Header("Access-Control-allow-origin", "*")

	var params ParamGetTxReceipt
	if err := c.ShouldBindQuery(&params); err != nil {
		c.String(-1, "invalid input data"+err.Error())
		return
	}

	bodyData, _ := json.Marshal(&kasKlaytnRPC{
		JsonRpc: "2.0",
		Method:  "klay_getTransactionReceipt",
		Params:  []string{params.TxHash},
		Id:      1,
	})

	ret, err := SendHTTPRequest("POST", url, bodyData)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	c.String(200, ret)
}

// POST kas/wallet_deployStoreListContract
// https://refs.klaytnapi.com/en/wallet/latest#operation/ContractDeployTransaction
// https://refs.klaytnapi.com/en/wallet/latest#operation/UFDValueTransferTransaction
func wallet_deployStoreListContract(c *gin.Context) {
	url := kasWalletUrl + "/v2/tx/fd/contract/deploy"

	bodyData, err := json.Marshal(&kasDeployTxFD{
		From:   conf.KasAccount,
		Value:  "0x0",
		Gas:    8000000,
		Input:  storeListContractDelpoyCode,
		Submit: true,
	})
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	ret, err := SendHTTPRequest("POST", url, bodyData)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	c.String(200, ret)
}

func ExecuteContract(to string, input string) (string, error) {
	txArgs := kasExecTxFD{
		From:     conf.KasAccount,
		Value:    "0x0",
		To:       to,
		GasLimit: 8000000,
		Submit:   true,
		Input:    input,
	}

	bodyData, _ := json.Marshal(&txArgs)
	ret, err := SendHTTPRequest("POST", kasWalletUrl+"/v2/tx/fd/contract/execute", bodyData)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func wallet_executeOrder(c *gin.Context) {
	c.Header("Access-Control-allow-origin", "*")

	var params ParamOrder
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(-1, "invalid input data"+err.Error())
		return
	}

	_abi, err := abi.JSON(bytes.NewBufferString(storeContractABI))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	inputParam, err := _abi.Pack("addOrder", params.Tid, params.Items, params.TotalPrice)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	ret, err := ExecuteContract(params.Contract, "0x"+Encode(inputParam))
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	c.String(200, ret)
}

func wallet_executeApproveOrder(c *gin.Context) {
	c.Header("Access-Control-allow-origin", "*")

	var params FormOrderId
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(-1, "invalid input data"+err.Error())
		return
	}

	_abi, err := abi.JSON(bytes.NewBufferString(storeContractABI))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	inputParam, err := _abi.Pack("approveOrder", params.OrderId, params.Reward)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	ret, err := ExecuteContract(params.Contract, "0x"+Encode(inputParam))
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	c.String(200, ret)
}

func wallet_executeDenyOrder(c *gin.Context) {
	c.Header("Access-Control-allow-origin", "*")

	var params FormOrderId
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(-1, "invalid input data"+err.Error())
		return
	}

	_abi, err := abi.JSON(bytes.NewBufferString(storeContractABI))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	inputParam, err := _abi.Pack("denyOrder", params.OrderId)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	ret, err := ExecuteContract(params.Contract, "0x"+Encode(inputParam))
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	c.String(200, ret)
}

func wallet_executeSetReward(c *gin.Context) {
	c.Header("Access-Control-allow-origin", "*")

	var params FormAddReward
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(-1, "invalid input data"+err.Error())
		return
	}

	if len(params.RewardHex)%2 == 1 {
		params.RewardHex = params.RewardHex[:2] + "0" + params.RewardHex[2:]
	}

	_abi, err := abi.JSON(bytes.NewBufferString(storeContractABI))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	inputParam, err := _abi.Pack("setRewardPolicy", params.MenuId, new(big.Int).SetBytes(hexutil.MustDecode(params.RewardHex)))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	ret, err := ExecuteContract(params.Contract, "0x"+Encode(inputParam))
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	c.String(200, ret)
}

// POST /kas/store/newattendance
//
func wallet_executeAddAttdendance(c *gin.Context) {
	c.Header("Access-Control-allow-origin", "*")

	var params ParamAddAttendance
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(-1, "invalid input data"+err.Error())
		return
	}

	_abi, err := abi.JSON(bytes.NewBufferString(storeContractABI))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	inputParam, err := _abi.Pack("addAttendance", params.ClassDate, params.PresentPrice, params.LatePrice)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	ret, err := ExecuteContract(params.ContractAddr, "0x"+Encode(inputParam))
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	c.String(200, ret)
}

// POST /kas/store
// https://refs.klaytnapi.com/en/wallet/v2#operation/FDContractDeployTransaction
// Json Params needed ->  params.go/ ParamDeployStore
func wallet_deployStoreContract(c *gin.Context) {
	url := kasWalletUrl + "/v2/tx/fd/contract/deploy"
	c.Header("Access-Control-allow-origin", "*")

	var params ParamDeployStore
	if err := c.ShouldBindJSON(&params); err != nil { // gin.Context의 Json body를 params와 묶는다.
		// Request 측에서 보낸 JSON 형태 정보를 params와 Binding 한다
		c.String(-1, "invalid input data"+err.Error()) // EOF 에러가 날 수도 있다.
		return
	}

	fmt.Println("LLmain-store1")

	txInput, err := GenStoreContractDeployCode(params.ClassName, params.Owner)
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	fmt.Println("LLmain-store2")

	bodyData, err := json.Marshal(&kasDeployTxFD{
		From:   conf.KasAccount,
		Value:  "0x0",
		Gas:    8000000,
		Input:  txInput,
		Submit: true,
	})
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	fmt.Println("LLmain-store3")

	ret, err := SendHTTPRequest("POST", url, bodyData)
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	c.String(200, ret)
}

// POST
func wallet_executeRegisterStore(c *gin.Context) {
	c.Header("Access-Control-allow-origin", "*")

	var params ParamRegisterClass
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(-1, "invalid input data"+err.Error())
		return
	}

	_abi, err := abi.JSON(bytes.NewBufferString(storeListContractABI))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	fmt.Println(params)
	fmt.Println(common.HexToAddress(params.ContractAddr))

	inputParam, err := _abi.Pack("addClass", params.ProfessorName, params.Name, params.SemesterYear,
		common.HexToAddress(params.Owner), common.HexToAddress(params.ContractAddr))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	ret, err := ExecuteContract(conf.MasterContractList, "0x"+Encode(inputParam))
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	c.String(200, ret)
}

func getMenu(addr string) ([]StoreMenu, error) {
	_abi, err := abi.JSON(bytes.NewBufferString(storeContractABI))
	if err != nil {
		return nil, err
	}

	// 1. get the number of menu
	inputParam, err := _abi.Pack("getNumMenus")
	if err != nil {
		return nil, err
	}

	ret, err := callContract(addr, "0x"+Encode(inputParam))
	if err != nil {
		return nil, err
	}

	var rpcRet RPCReturnString
	if err := json.Unmarshal([]byte(ret), &rpcRet); err != nil {
		return nil, err
	}

	numMenu, err := strconv.ParseUint(rpcRet.Result[len(rpcRet.Result)-8:], 16, 64)
	if err != nil {
		return nil, err
	}

	// 2. get all menu
	var allMenu []StoreMenu
	// allMenu := make(map[uint64]StoreMenu)
	for i := uint64(0); i < numMenu; i++ {
		var menu StoreMenu

		inputParam, err := _abi.Pack("getMenu", uint32(i))
		if err != nil {
			return nil, err
		}

		ret, err = callContract(addr, "0x"+Encode(inputParam))
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(ret), &rpcRet); err != nil {
			return nil, err
		}

		if err := _abi.Unpack(&menu, "getMenu", hexutil.MustDecode(rpcRet.Result)); err != nil {
			return nil, err
		}

		menu.Id = uint32(i)
		menu.RewardHex = hexutil.Encode(menu.Reward.Bytes())
		allMenu = append(allMenu, menu)
	}

	return allMenu, nil
}

func node_getMenu(c *gin.Context) {
	c.Header("Access-Control-allow-origin", "*")

	var txArgs FormAddress
	if err := c.ShouldBindQuery(&txArgs); err != nil {
		c.String(-1, "invalid input data"+err.Error())
		return
	}

	menus, err := getMenu(txArgs.Address)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	c.JSON(200, menus)
}

func node_getStores(c *gin.Context) {
	c.Header("Access-Control-allow-origin", "*")

	_abi, err := abi.JSON(bytes.NewBufferString(storeListContractABI))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	inputParam, err := _abi.Pack("getStores")
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	ret, err := callContract(conf.MasterContractList, "0x"+Encode(inputParam))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	var rpcRet RPCReturnString
	if err := json.Unmarshal([]byte(ret), &rpcRet); err != nil {
		c.String(-1, err.Error())
		return
	}

	var arr []StoreInfoArray
	if err := _abi.Unpack(&arr, "getStores", hexutil.MustDecode(rpcRet.Result)); err != nil {
		fmt.Println(rpcRet.Result)
		fmt.Println(hexutil.MustDecode(rpcRet.Result))
		c.String(-1, err.Error())
		return
	}

	c.JSON(200, arr)
}

func node_getLogs(c *gin.Context) {
	c.Header("Access-Control-allow-origin", "*")

	var txArgs FormGetLogs
	if err := c.ShouldBindQuery(&txArgs); err != nil {
		c.String(-1, "invalid input data"+err.Error())
		return
	}

	receipts, err := getLogs(txArgs.FromBlock, txArgs.ToBlock, txArgs.Contract)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	c.JSON(200, receipts)
}

func node_getOrder(c *gin.Context) {
	c.Header("Access-Control-allow-origin", "*")

	var txArgs FormOrderId
	if err := c.ShouldBindQuery(&txArgs); err != nil {
		c.String(-1, "invalid input data"+err.Error())
		return
	}

	_abi, err := abi.JSON(bytes.NewBufferString(storeContractABI))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	inputData, err := _abi.Pack("getOrder", txArgs.OrderId)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	ret, err := callContract(txArgs.Contract, hexutil.Encode(inputData))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	var rpcRet RPCReturnString
	if err := json.Unmarshal([]byte(ret), &rpcRet); err != nil {
		c.String(-1, err.Error())
		return
	}

	var order Order
	if err := _abi.Unpack(&order, "getOrder", hexutil.MustDecode(rpcRet.Result)); err != nil {
		c.String(-1, err.Error())
		return
	}

	c.JSON(200, order)
}

// To Make into module
func node_getNumAttendances(c *gin.Context) {
	c.Header("Access-Control-allow-origin", "*")

	var params FormAddress
	if err := c.ShouldBindQuery(&params); err != nil {
		c.String(-1, "invalid input data "+err.Error())
		return
	}

	_abi, err := abi.JSON(bytes.NewBufferString(storeContractABI))
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	inputParam, err := _abi.Pack("getNumAttendances")
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	fmt.Println(params.Address)

	ret, err := callContract(params.Address, "0x"+Encode(inputParam))
	if err != nil {
		c.String(-1, err.Error())
		return
	}
	fmt.Println(ret)

	var rpcRet RPCReturnString
	if err := json.Unmarshal([]byte(ret), &rpcRet); err != nil {
		c.String(-1, err.Error())
		return
	}

	fmt.Println(rpcRet)
	numAttendance, err := strconv.ParseUint(rpcRet.Result[len(rpcRet.Result)-8:], 16, 64)
	if err != nil {
		c.String(-1, err.Error())
		return
	}

	c.JSON(200, numAttendance)
}
