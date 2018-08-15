package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	sdkClient "github.com/suddutt1/fabricgosdkclientcore"
)

//InitClient is intializing the client
func InitClient(configFilePath string) (bool, *sdkClient.FabricSDKClient) {
	client := new(sdkClient.FabricSDKClient)
	if client.Init(configFilePath) {
		return true, client
	}
	return false, nil
}

//Initialize initialize the client
func Initialize(configFilePath string, isNewAdmin bool) (bool, *sdkClient.FabricSDKClient) {
	isSuccess, client := InitClient(configFilePath)
	if !isSuccess {
		fmt.Println("Client initialization failure")
		return false, nil
	}
	fmt.Println("Client initialization success")
	return client.ErollOrgAdmin(isNewAdmin, "admin"), client

}

//EnrollOrgUser a new user in the system
func EnrollOrgUser(configFilePath, userID, org string) bool {
	return false
}

//InstallAdminCerts Will install admin cert in IBP after checking
//if that is not alreadt installed. Also restart the peers
//and sync channels
func InstallAdminCerts(ibpCredentialPath, specFile string) bool {
	return false
}

//DeployeCC install chanin code
func DeployeCC(configFilePath string, specMap map[string]interface{}) bool {
	isSucess, client := Initialize(configFilePath, false)
	if !isSucess {
		return false
	}
	ccID := getString(specMap["ccID"])
	version := getString(specMap["version"])
	goPath := getString(specMap["goPath"])
	ccPath := getString(specMap["ccSrcRootPath"])
	return client.InstallChainCode(ccID, version, goPath, ccPath, nil)

}

//InstantiateChainCode Instantiates chain code . They must to be deployed before hand
func InstantiateChainCode(configFilePath string, specMap map[string]interface{}) bool {
	isSucess, client := Initialize(configFilePath, false)
	if !isSucess {
		return false
	}
	ccID := getString(specMap["ccID"])
	version := getString(specMap["version"])
	ccPath := getString(specMap["ccSrcRootPath"])
	channel := getString(specMap["channel"])
	policy := getString(specMap["ccPolicy"])
	initParams := getByteSlice(specMap["initParams"])
	isSucess, err := client.InstantiateCC(channel, ccID, ccPath, version, initParams, policy, nil)
	if err != nil {
		fmt.Println("Instantiation failure")
		return false
	}
	return isSucess

}

//UpgradeChainCode should upgrade the chain code
func UpgradeChainCode(configFilePath string, specMap map[string]interface{}) bool {
	isSucess, client := Initialize(configFilePath, false)
	if !isSucess {
		return false
	}
	ccID := getString(specMap["ccID"])
	version := getString(specMap["version"])
	ccPath := getString(specMap["ccSrcRootPath"])
	channel := getString(specMap["channel"])
	policy := getString(specMap["ccPolicy"])
	initParams := getByteSlice(specMap["initParams"])
	isSucess, err := client.UpdateCC(channel, ccID, ccPath, version, initParams, policy, nil)
	if err != nil {
		fmt.Println("Instantiation failure")
		return false
	}
	return isSucess
}
func QueryChaninCode(configFilePath string, specMap map[string]interface{}) bool {
	return false

}
func InvokeChainCode(configFilePath string, specMap map[string]interface{}) bool {
	return false

}

//GenerateSpecFile prints the spec file for varios activities
func GenerateSpecFile(spec string) bool {
	switch spec {
	case "cc-deploy":
	case "cc-instantiate":
	case "cc-upgrade":

		specInstallUpgrade := `
		{
			"description":"Install/Upgrade Chain Code",
			"ccID":"<ccID>",
			"version":"<ccVersion>",
			"channel":"<channelName",
			"goPath":"",
			"ccSrcRootPath":"",
			"initParams":["param1","param2"],
			"ccPolicy":"<cc policy>"
		}
		`
		fmt.Println("Chain code install/upgrade spec \n", specInstallUpgrade)
	case "add-admin-cert":
		specAdminCert := `
		{
			"descrption":"Admin cert install specification",
			"peerID":"peer1-org1",
			"certFilePath":"",
			"keyFilePath":"",,
			"adminCertName":"<unique-cert-name>"

		}
		`
		fmt.Println("Admin cert installation spec \n", specAdminCert)
	case "cc-query":
		specCCQuery := `
		{

		}
		`
		fmt.Println("Chain code query spec\n", specCCQuery)
	case "cc-invoke":
		specCCInvoke := `
		{

		}
		`
		fmt.Println("Chain code invoke spec\n", specCCInvoke)
	}

	return true
}
func getSpecificationMap(specFile string) (bool, map[string]interface{}) {
	specBytes, err := ioutil.ReadFile(specFile)
	if err != nil {
		return false, nil
	}
	specMap := make(map[string]interface{})
	if err := json.Unmarshal(specBytes, &specMap); err != nil {
		return false, nil
	}
	return true, specMap
}
func isJSON(bytes []byte) (interface{}, bool) {
	var genericIntfc interface{}
	if err := json.Unmarshal(bytes, &genericIntfc); err == nil {
		return genericIntfc, true
	}
	return nil, false
}
func getString(strIntfc interface{}) string {
	if strIntfc != nil {
		if str, ok := strIntfc.(string); ok {
			return str
		}
	}
	return ""
}
func getByteSlice(strIntfc interface{}) [][]byte {
	if strIntfc != nil {

		if intfcSlice, ok := strIntfc.([]interface{}); ok {
			retList := make([][]byte, 0)
			for _, intfc := range intfcSlice {
				retList = append(retList, []byte(getString(intfc)))
			}
			fmt.Printf("\nBefore return %+v", retList)
			return retList
		}
	}
	return make([][]byte, 0)
}

/*func buildArgsList(strIntfc interface{}) [][]byte {
	outputBytesSlice := make([][]byte, 0)
	strSlice := getStringSlice(strIntfc)
	for _, str := range strSlice {
		outputBytesSlice = append(outputBytesSlice, []byte(str))
	}
	return outputBytesSlice
}*/

func isSpecRequired(command string) bool {
	isRequired := false
	switch command {
	case "cc-deploy":
	case "cc-instantiate":
	case "cc-upgrade":
	case "cc-query":
	case "cc-invoke":
		isRequired = true
	}
	return isRequired
}
func main() {
	//ibpclient  <command> [param1] [param2] [param3] --config="" --specFile=""

	var specMap map[string]interface{}
	var isValidSpec bool

	var specFile string
	configFile := ""
	flag.StringVar(&configFile, "config", "", "Please provide the client config yaml/json file")
	flag.StringVar(&specFile, "spec", "", "Please provide the spec file")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 || len(configFile) == 0 {
		usage()
		os.Exit(1)
	}
	command := args[0]
	isSuccess := false
	fmt.Println("Using config file ", configFile)
	if isSpecRequired(command) {
		isValidSpec, specMap = getSpecificationMap(specFile)
		if !isValidSpec {
			fmt.Println("Invalid json provided as spec")
			os.Exit(2)
		}
	}

	switch command {
	case "init":
		isSuccess, _ = Initialize(configFile, true)
	case "cc-deploy":
		isSuccess = DeployeCC(configFile, specMap)
	case "cc-instantiate":
		isSuccess = InstantiateChainCode(configFile, specMap)
	case "cc-upgrade":
		isSuccess = UpgradeChainCode(configFile, specMap)

	case "cc-query":
		isSuccess = QueryChaninCode(configFile, specMap)
	case "cc-invoke":
		isSuccess = InvokeChainCode(configFile, specMap)
	case "add-admin-cert":
		isSuccess = InstallAdminCerts(configFile, specFile)
	case "spec-gen":
		if len(args) > 1 {
			isSuccess = GenerateSpecFile(args[1])
		} else {
			fmt.Println("Invalid generate command")
		}
	default:
		fmt.Println("Usage :---")
		flag.Usage()
	}

	fmt.Println("Done....")
	if !isSuccess {
		os.Exit(2)
	}

}
func usage() {
	fmt.Println("ibpclient --config=\"config file path\" [--spec=\"spec json file path\"]  <command> [param1] [param2] [param3] ")
	fmt.Println("<command is in one of the following format> ")
	fmt.Println("--config=<config file path> init ")
	fmt.Println("\t intialize the tool")
	fmt.Println("--config=<config file path> --spec=<spec file path> cc-deploy")
	fmt.Println("\t install the chain code")
	fmt.Println("--config=<config file path> --spec=<spec file path> cc-instantiate")
	fmt.Println("\t instantiates the chain code")
	fmt.Println("--config=<config file path> --spec=<spec file path> cc-upgrade ")
	fmt.Println("\t upgrade the chain code")
	fmt.Println("--config=<config file path> --spec=<spec file path> cc-query ")
	fmt.Println("\t invoke a chain code query")
	fmt.Println("--config=<config file path> --spec=<spec file path> cc-invoke ")
	fmt.Println("\t performs a transaction ")
	fmt.Println("add-admin-cert")
	fmt.Println("\t add a new admin cert in IBP organization <Not implemented>")
	fmt.Println("spec-gen <command>")
	fmt.Println("\t generate the specification json")

}
