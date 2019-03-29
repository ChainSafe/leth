package core

import (
	"os"
	"fmt"

	"github.com/ChainSafe/leth/logger"
)

func PrepNetwork(network string) Network {
	// read config file
	file, err := ReadConfig()
	if err != nil {
		logger.FatalError("no config.json file found.")
		os.Exit(1)
	}

	config, err := UnmarshalConfig(file)
	if err != nil {
		logger.FatalError(fmt.Sprintf("could not read config.json: %s", err))
	}

	ntwk := config.Networks[network]
	ntwk.Name = network

	return ntwk
}

func Migrate(network string, contract string) error {
	ntwk := PrepNetwork(network)

	client, err := Dial(ntwk.Url)
	if err != nil {
		logger.FatalError("cannot dial client; likely incorrect url in config.json")
	}

	contracts := []string{contract}

	if ntwk.Name == "testrpc" || ntwk.Name == "ganache" || ntwk.Name == "ganache-cli" {
		accounts, err := GetAccounts(ntwk.Url)
		if err != nil {
			logger.FatalError(fmt.Sprintf("unable to get accounts from client url: %s", err))
		}
		//logger.Info(fmt.Sprintf("accounts: %s", accounts))
		PrintAccounts(accounts)

		if ntwk.From == "" {
			ntwk.From = accounts[0]
		}

		err = DeployTestRPC(ntwk, contracts)
		if err != nil {
			logger.FatalError("could not deploy contracts.")
		}
	} else {
		ks := NewKeyStore(ntwk.Keystore)
		ksaccounts := ks.Accounts()
		PrintKeystoreAccounts(ksaccounts)
		err = Deploy(client, ntwk, contracts, ks)
		if err != nil {
			logger.FatalError("could not deploy contracts.")
		}
	}

	return nil
}