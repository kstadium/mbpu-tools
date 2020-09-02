package main

import (
  "os"
  
  "github.com/hyperledger/fabric/mbpu-tools/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
