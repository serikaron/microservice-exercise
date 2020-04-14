package main

import (
	"log"
	"mse/auth/internal"
	"mse/pkg"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	pkg.ParseItem([]pkg.FlagItem{pkg.AuthAddr, pkg.CertsPath})

	s := internal.AuthService{}
	s.Run(pkg.AuthAddr.Addr(), pkg.CertsPath.Pem(), pkg.CertsPath.Key())
}
