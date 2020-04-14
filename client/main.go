package main

import (
	"log"
	"mse/client/internal"
	"mse/pkg"
	"os"
)

func init() {
	pkg.ParseItem([]pkg.FlagItem{
		pkg.AuthAddr,
		pkg.CertsPath,
		pkg.ChatAddr,
	})
}

func main() {
	jwt, err := internal.Login(pkg.AuthAddr.Addr(), pkg.CertsPath.Pem())
	if err != nil {
		log.Fatal(err)
	}

	typeHins := "\n[1]chater\n[2]monitor\nselect 1 or 2 or quit: "
	input := internal.GetInput(typeHins)
	for {
		switch input {
		case "quit":
			os.Exit(0)
		case "1":
			internal.StartChat(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem(), jwt)
		case "2":
			internal.StartMonitor(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem(), jwt)
		}
		input = internal.GetInput(typeHins)
	}
}
