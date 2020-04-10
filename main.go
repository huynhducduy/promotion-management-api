package main

import (
	log "github.com/sirupsen/logrus"
	"promotion-management-api/cmd/app"
)

func main() {
	log.Fatal(app.Run())
}
