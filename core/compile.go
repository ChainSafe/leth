package core

import (
	"fmt"
	"path"
	"path/filepath"
	"os/exec"
	"os"
	"strings"

	"github.com/ChainSafe/leth/logger"
)

func Compile() ([]string, error) {
	buildexists, err := Exists("build/")
	if buildexists {
		os.RemoveAll("./build")
	}
	os.Mkdir("./build", os.ModePerm)

	dir, _ := filepath.Abs("contracts/")

	files, err := SearchDirectory(dir)

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
		return nil, err
	}

	contracts := []string{}

	//fmt.Println(files)
	for _, file := range files {
		if(path.Ext(file) == ".sol") {
			err = compile(file)
			if err != nil {
				return nil, err
			}
			contracts = append(contracts, file)
		}
	}

	return contracts, nil
}

func compile(contract string) error { 
	logger.Info(fmt.Sprintf("compiling %s", path.Base(contract)))

	// generate bytecode
    app := "solc"
    arg0 := "--bin"
    arg1 := contract
    arg2 := "-o"
    arg3 := "build/"

    cmd := exec.Command(app, arg0, arg1, arg2, arg3)
    stdout, err := cmd.CombinedOutput()

    out := string(stdout)
    if strings.Contains(out, "Warning") {
    	logger.CompilerWarn(out)
	} else if strings.Contains(out, "Error") {
		logger.CompilerError(out)
	} else {
		fmt.Println(out)
	}

    if err != nil {
        return err
    }

    // generate abi
    app = "solc"
    arg0 = "--abi"
    arg1 = contract
    arg2 = "-o"
    arg3 = "build/"

    cmd = exec.Command(app, arg0, arg1, arg2, arg3)
    stdout, err = cmd.Output()

    if err != nil {
        return err
    }

    //print(string(stdout))	
    return nil
}