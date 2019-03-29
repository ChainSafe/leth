package migrations

import (
	"fmt"

	"github.com/ChainSafe/leth/core"
)

func Migrate() {
	err := core.Migrate("testnet", "Example")
	if err != nil {
		fmt.Println("could not deploy Example.sol to testnet")
	}

	// err = core.Migrate("testnet", "ExampleLib")
	// if err != nil {
	// 	fmt.Println("could not deploy ExampleLib.sol to testnet")
	// }
}