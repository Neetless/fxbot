package main

import (
	"log"
	"os"
	"sandspace/fxgot/bitflyer"
)

var logger = log.New(os.Stdout, "main: ", log.Lshortfile)

func main() {
	client := bitflyer.New()
	markets, err := client.GetMarkets()
	if err != nil {
		logger.Print(err)
		return
	}
	logger.Print(markets)

	executions, err := client.GetExecutions()
	if err != nil {
		logger.Print(err)
		return
	}
	logger.Print(executions)
}
