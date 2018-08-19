package main

import (
	"log"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/microhq/federation-srv/federation"
	"github.com/microhq/federation-srv/handler"

	"github.com/micro/go-config"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/file"
	proto "github.com/microhq/federation-srv/proto/federation"
)

var (
	defaultFile = "federation.json"
)

func main() {

	service := micro.NewService(
		micro.Name("go.micro.srv.federation"),
		micro.Flags(
			cli.StringFlag{
				Name:   "config_source",
				EnvVar: "CONFIG_SOURCE",
				Usage:  "Source to read the config from e.g file, platform",
			},
		),
	)

	// initialise service
	service.Init(
		micro.Action(func(c *cli.Context) {
			var src string

			parts := strings.Split(c.String("config_source"), ":")

			if len(parts) > 0 {
				src = parts[0]
			}

			var source source.Source

			switch src {
			case "file":
				fileName := defaultFile

				if len(parts) > 1 {
					fileName = parts[1]
				}

				log.Println("Using file source:", fileName)
				source = file.NewSource(file.WithPath(fileName))
			default:
				fileName := defaultFile

				if len(parts) > 1 {
					fileName = parts[1]
				}

				log.Println("Using file source:", fileName)
				source = file.NewSource(file.WithPath(fileName))
			}

			config.Load(source)
			federation.Init(config.DefaultConfig, service)
		}),
		micro.BeforeStart(federation.Run),
	)

	proto.RegisterFederationHandler(service.Server(), new(handler.Federation))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
