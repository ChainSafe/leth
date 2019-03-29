package main

import (
	"fmt"
	"log"
	"flag"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"encoding/json"
	"io/ioutil"

	"github.com/ChainSafe/leth/core"
	"github.com/ChainSafe/leth/logger"
	//"github.com/ChainSafe/leth/test"
	"github.com/ChainSafe/leth/migrations"

	//"github.com/ethereum/go-ethereum/ethclient"
	//"github.com/ethereum/go-ethereum/accounts"
	//"github.com/ethereum/go-ethereum/accounts/keystore"
)

func main() {
	// flags
	help := flag.Bool("help", false, "print out command-line options")

	// init subcommand
	initCommand := flag.NewFlagSet("init", flag.ExitOnError)

	// bind subcommand 
	bindCommand := flag.NewFlagSet("bind", flag.ExitOnError)

	// compile subcommand and flags
	compileCommand := flag.NewFlagSet("compile", flag.ExitOnError)
	bindFlag := compileCommand.Bool("bind", true, "specify whether to create bindings for contracts while compiling")

	// migrate subcommand 
	migrateCommand := flag.NewFlagSet("migrate", flag.ExitOnError)
	network := migrateCommand.String("network", "default", "specify network to connect to (configured in config.json)")

	// deploy subcommand and flags
	deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)
	network = deployCommand.String("network", "default", "specify network to connect to (configured in config.json)")

	// test subcommand
	testCommand := flag.NewFlagSet("test", flag.ExitOnError)
	testContract := testCommand.String("test", "Test", "specify which function to call initially when testing")

	flag.Parse() 
	if *help {
		fmt.Println("\t\x1b[93mleth help\x1b[0m")
		fmt.Println("\tleth bind: create go bindings for all contracts in contracts/ directory and save in bindings/")
		fmt.Println("\tleth compile: compile all contracts in contracts/ directory and save results in build/. `compile` will automatically execute `bind`; to compile with out binding, use --bind=false")
		fmt.Println("\tleth migrate: run migration file and migrate to a network specified by --network")
		fmt.Println("\tleth deploy: deploy all contracts in contracts/ directory and save results of deployment in deployed/. specify network name with `--network NETWORK_NAME`. if no network is provided, leth will connect to the default network as specified in config.json")
		fmt.Println("\tleth test: run tests in test/ directory")
		os.Exit(0)
	} 

	// subcommands
	if len(os.Args) > 1 {
		switch os.Args[1]{
			case "init":
				initCommand.Parse(os.Args[2:])
			case "bind":
				bindCommand.Parse(os.Args[2:])
			case "compile":
				compileCommand.Parse(os.Args[2:])
			case "migrate":
				migrateCommand.Parse(os.Args[2:])
			case "deploy":
				deployCommand.Parse(os.Args[2:])
			case "test":
				testCommand.Parse(os.Args[2:])
			default:
				// continue
		}
	} else {
		os.Exit(0)
	}

	if initCommand.Parsed() {
		lethInit()
		os.Exit(0)
	}

	if bindCommand.Parsed() {
		bind()
		os.Exit(0)
	}

	if compileCommand.Parsed() {
		//contractArgs := compileCommand.Args()
		compile(*bindFlag)
		os.Exit(0)	
	} 

	if migrateCommand.Parsed() {
		migrate()
		os.Exit(0)
	}

	if deployCommand.Parsed() {
		deploy(*network)
		os.Exit(0)
	}

	if testCommand.Parsed() {
		testrun(*testContract)
		os.Exit(0)
	}
}

func lethInit() {
	files, err := core.SearchDirectory("./")
	if len(files) > 1 {
		logger.FatalError("cannot init in non-empty directory.")
	}
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}

	os.Mkdir("./contracts", os.ModePerm)
	os.Mkdir("./migrations", os.ModePerm)
	os.Mkdir("./keystore", os.ModePerm)
	os.Mkdir("./test", os.ModePerm)

	jsonStr, err := json.MarshalIndent(core.DefaultConfig, "", "\t")
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}

	ioutil.WriteFile("./config.json", jsonStr, os.ModePerm)

	mainStr := "package main\n\n" +
				"import (\n" +
					"\t// \"your-project/test\"\n" +
				")\n\n" +
				"func main() {\n" +
					"\ttest.Test()\n" +
				"}"
	mainFile :=  []byte(mainStr)
	ioutil.WriteFile("./main.go", mainFile, os.ModePerm)
}

func bind() {
	//fmt.Println(contracts)
	err := core.Bindings()
	if err != nil {
		logger.FatalError(fmt.Sprintf("could not create bindings: %s", err))
	} else {
		logger.Info("generation of bindings completed. saving bindings in bindings/ directory.")
	}
} 

func compile(bindFlag bool) ([]string) {
	contracts, err := core.Compile()
	if err != nil {
		logger.FatalError(fmt.Sprintf("compilation error: %s", err))
	} else {
		logger.Info("compilation completed. saving binaries in build/ directory.")
	}
	if bindFlag {
		bind()
	}
	return contracts
}

// set up migration to a network
// similar to deploy, except execute migrations/nigrate.go
func migrate() {
	migrations.Migrate()
}

// set up deployment to network
// compile, read config, dial network, set up accounts
func deploy(network string) {
	// compilation of contracts, if needed
	contracts := []string{}
	buildexists, err := core.Exists("build/")
	if !buildexists {
		logger.Info("build/ directory not found. compiling contracts...")
		compile(false) // don't need to generate bindings for deployment
	}

	files, err := core.SearchDirectory("./build")
	if err != nil {
		log.Fatal(err)
	} else if len(files) < 2 {
		logger.Info("build/ directory empty. compiling contracts...")
		compile(false)
		files, err = core.SearchDirectory("./build")
	} else {
		for _, file := range files {
			if(path.Ext(file) == ".bin") {
				contracts = append(contracts, file)
			}
		}
	}

	names := core.GetContractNames(contracts)

	ntwk := core.PrepNetwork(network)

	// dial client for network
	client, err := core.Dial(ntwk.Url)
	if err != nil {
		logger.FatalError("cannot dial client; likely incorrect url in config.json")
	}

	if ntwk.Name == "testrpc" || ntwk.Name == "ganache" || ntwk.Name == "ganache-cli" {
		accounts, err := core.GetAccounts(ntwk.Url)
		if err != nil {
			logger.FatalError(fmt.Sprintf("unable to get accounts from client url: %s", err))
		}
		//logger.Info(fmt.Sprintf("accounts: %s", accounts))
		core.PrintAccounts(accounts)

		if ntwk.From == "" {
			ntwk.From = accounts[0]
		}

		err = core.DeployTestRPC(ntwk, names)
		if err != nil {
			logger.FatalError("could not deploy contracts.")
		}
	} else {
		ks := core.NewKeyStore(ntwk.Keystore)
		ksaccounts := ks.Accounts()
		core.PrintKeystoreAccounts(ksaccounts)
		err = core.Deploy(client, ntwk, names, ks)
		if err != nil {
			logger.FatalError("could not deploy contracts.")
		}
	}
}

// note: all this function does is run the main function of the directory it's in.
// this isn't specific to testing - could be used for migrations as well.
// merge this with migrations; write documentation on how to use the main.go file
func testrun(contract string) {
	fp, _ := filepath.Abs("./main.go")
	cmd := exec.Command("go", "run", fp)
	stdout, err := cmd.CombinedOutput()
	out := string(stdout)
	logger.Info(fmt.Sprintf("executing %s...\n%s", fp, out))
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
}