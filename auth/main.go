package main

import (
	"log"
	"mse/pkg"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	pkg.ParseItem([]pkg.FlagItem{pkg.AuthAddr, pkg.CertsPath})

	s := AuthService{}
	s.Run(pkg.AuthAddr.Addr(), pkg.CertsPath.Pem(), pkg.CertsPath.Key())
}
