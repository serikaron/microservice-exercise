package internal

import (
	"log"
	"mse/pkg"
)

func IdentifyWithPassword(name string, password string) (*pkg.Identity, error) {
	if name != password {
		log.Println(pkg.LoginErr)
		return nil, pkg.LoginErr
	}

	return &pkg.Identity{Name: name}, nil
}
