package main

import (
	"log"

	pds "github.com/mikew79/port-domain-service/application/portdomainserver"
)

func main() {

	var server pds.PortDomainServer
	err := server.InitPortDomainService()
	if err != nil {
		log.Fatalf("Failed To initialise Port Domian Server: %s", err)
		return
	}
	server.Go()

}
