package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"path"

	"github.com/ChainSafe/leth/logger"
)

func Bindings() (error) {
	dir, _ := filepath.Abs("build/")

	contracts, err := SearchDirectoryForAbi(dir)
	if err != nil {
		return err
	}

	err = makeBindingDir()
	if err != nil {
		return err
	}

	for _, contract := range contracts {
		logger.Info(fmt.Sprintf("generating binding for %s", path.Base(contract)))
		err := bind(contract)
		if err != nil {
			logger.Error(fmt.Sprintf("could not generate binding for %s: %s", contract, err))
			return err
		}
	}

	return nil
}

func makeBindingDir() error {
	bindingexists, err := Exists("bindings/")
	if err != nil {
		return err
	}
	if bindingexists {
		os.RemoveAll("./bindings")
	}
	os.Mkdir("./bindings", os.ModePerm)
	return nil
}

func bind(contract string) (error) {
	name := GetContractName(contract)
	output, _ := filepath.Abs("./bindings/" + name + ".go")

    app := "abigen"
    arg0 := "--abi"
    arg1 := contract
    arg2 := "--pkg"
    arg3 := "bindings"
    arg4 := "--type"
    arg5 := name
    arg6 := "--out"
    arg7 := output

    cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
    stdout, err := cmd.CombinedOutput()
    if err != nil {
    	return err
    }

    out := string(stdout)
    if false { logger.Info(out) }
    return nil
}