package main

import (
	elasticbeat "github.com/radoondas/elasticbeat/beat"

	"github.com/elastic/beats/libbeat/beat"
)

// You can overwrite these, e.g.: go build -ldflags "-X main.Version 1.0.0-beta3"
var Version = "1.0.0-dev"
var Name = "elasticbeat"

func main() {
	beat.Run(Name, Version, elasticbeat.New())
}
