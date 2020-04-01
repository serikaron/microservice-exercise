package main

import (
	"log"
	"mse/pkg"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	pkg.ParseItem([]pkg.FlagItem{pkg.AuthAddr})

	s := AuthService{}
	s.Run(pkg.AuthAddr.Addr())
}
