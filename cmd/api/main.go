package main

import (
	"github.com/mmpg/api"
	"github.com/mmpg/api/jutge"
)

func main() {
	api.Run(jutge.ValidateCredentials)
}
