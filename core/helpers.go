package core

import (
	"os"
	"encoding/hex"
	"encoding/json"
	"path"
	"path/filepath"
	"io/ioutil"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/ChainSafe/leth/logger"
)

func ReadConfig() ([]byte, error) {
	path, _ := filepath.Abs("./config.json")
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}	
	return file, nil
}

func UnmarshalConfig(file []byte) (*Config, error) {
	conf := new(Config)
	err := json.Unmarshal(file, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func writeDeployment(network string, deployed map[string]string) error {
	deployedexists, err := Exists("deployed/")
	if err != nil {
		return err
	}
	if !deployedexists {
		os.Mkdir("./deployed", os.ModePerm)
	}

	jsonStr, err := json.MarshalIndent(deployed, "", "\t")
	if err != nil {
		return err
	}

	path, _ := filepath.Abs("./deployed/" + network + ".json")

	fileexists, err := Exists(path)
	if fileexists {
		os.Remove(path)
	}

	err = ioutil.WriteFile(path, jsonStr, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func getBytecode(contract string) ([]byte, error) {
	path, _ := filepath.Abs("./build/" + contract + ".bin")
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}	

	hexString := fmt.Sprintf("%s", file)
	//fmt.Println(hexString)

	hexBytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	return hexBytes, nil
}

func Exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func SearchDirectory(dir string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func SearchDirectoryForSol(dir string) ([]string, error) {
	files, err := SearchDirectory(dir)
	if err != nil {
		return []string{}, err
	}
	contracts := []string{}

	//fmt.Println(files)
	for _, file := range files {
		if(path.Ext(file) == ".sol") {
			contracts = append(contracts, file)
		}
	}

	return contracts, nil
}

func SearchDirectoryForAbi(dir string) ([]string, error) {
	files, err := SearchDirectory(dir)
	if err != nil {
		return []string{}, err
	}
	contracts := []string{}

	//fmt.Println(files)
	for _, file := range files {
		if(path.Ext(file) == ".abi") {
			contracts = append(contracts, file)
		}
	}

	return contracts, nil
}

func GetContractName(contract string) (string) {
	base := path.Base(contract)
	ext := path.Ext(contract)
	return base[0:len(base)-len(ext)]
}

func GetContractNames(contracts []string) ([]string) {
	names := []string{}
	for _, contract := range contracts {
		names = append(names, GetContractName(contract))
	}
	return names
}

func BinToSol(contracts []string) ([]string) {
	names := []string{}
	for _, contract := range contracts {
		name := GetContractName(contract)
		names = append(names, fmt.Sprintf("%s.sol", name))
	}
	return names
}

func NewKeyStore(path string) (*keystore.KeyStore) {
	newKeyStore := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	return newKeyStore
}

func PrintAccounts(accounts []string) {
	for i, account := range accounts {
		logger.Info(fmt.Sprintf("account %d: %s", i, account))
	}
}

func PrintKeystoreAccounts(accounts []accounts.Account) {
	for i, account := range accounts {
		logger.Info(fmt.Sprintf("account %d: %s", i, account.Address.Hex()))
	}
}