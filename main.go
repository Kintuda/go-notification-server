package main

import (
	"log"

	cmd "github.com/Kintuda/notification-server/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
