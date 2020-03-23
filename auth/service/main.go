package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	host := flag.String("auth-service-host", "", "auth service host")
	port := flag.Uint("auth-service-port", 0, "auth service port")
	internalHost := flag.String("auth-internal-service-host", "", "auth internal service host")
	internalPort := flag.Uint("auth-internal-service-port", 0, "auth internal service port")
	flag.Parse()
	s := AuthService{}
	addr := fmt.Sprintf("%s:%v", *host, *port)
	go s.Run(addr)
	is := AuthInternalService{}
	internamAddr := fmt.Sprintf("%s:%v", *internalHost, *internalPort)
	go is.Run(internamAddr)

	for {
		time.Sleep(100 * time.Millisecond)
	}
}
