package main

import (
	"flag"
	"log"
	"os"
	"syscall"

	"github.com/steevehook/expenses-rest-api/app"
)

func main() {
	configPath := flag.String(
		"config",
		"./config/app-config.yaml",
		"Path to the application config file",
	)
	flag.Parse()

	application, err := app.Init(*configPath)
	if err != nil {
		log.Fatal("could not init application: ", err)
	}

	go func() {
		if err := application.Start(); err != nil {
			log.Fatal("could not start application: ", err)
		}
	}()

	app.ListenToSignals([]os.Signal{os.Interrupt, syscall.SIGTERM}, application)
}
