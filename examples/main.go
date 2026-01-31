// Command main runs all viem-go examples.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	examples "github.com/ChefBingbong/viem-go/examples/src"
)

func main() {
	// Parse command line flags
	exampleName := flag.String("example", "all", "Name of example to run (client, read-contract, all)")
	flag.Parse()

	fmt.Println("╔═══════════════════════════════════════╗")
	fmt.Println("║        viem-go Examples Runner        ║")
	fmt.Println("╚═══════════════════════════════════════╝")
	// Run specified example or all
	switch *exampleName {
	case "all":
		runAll()
	case "client":
		if err := examples.RunClientExample(); err != nil {
			fmt.Println(err)
		}
	case "read-contract":
		if err := examples.RunReadContractExample(); err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown example: %s\n", *exampleName)
		os.Exit(1)
	}

	fmt.Println("========================================")
	fmt.Println("All examples completed successfully!")
	fmt.Println("========================================")
}

func runAll() {
	examples.RunClientExample()
	time.Sleep(time.Second * 5)
	examples.RunReadContractExample()
}
