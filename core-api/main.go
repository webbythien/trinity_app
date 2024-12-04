package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hrshadhin/fiber-go-boilerplate/cmd/server"
	"github.com/hrshadhin/fiber-go-boilerplate/cmd/worker"
	_ "github.com/hrshadhin/fiber-go-boilerplate/docs" // load API Docs files (Swagger)
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/config"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/utils"
	"github.com/urfave/cli"
)

var (
	client *cli.App
)

func init() {
	client = cli.NewApp()
	client.Name = "Trinity Core API"
	client.Usage = "Trinity core api worker and handler"
	client.Version = "0.0.1"
}

// @title Fiber Go API
// @version 1.0
// @description Fiber go web framework based REST API boilerplate
// @termsOfService
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @host localhost:5005
func main() {
	// setup various configuration for app
	config.LoadAllConfigs(".env")

	client.Commands = []cli.Command{
		{
			Name:  "worker",
			Usage: "launch machinery worker",
			Action: func(c *cli.Context) error {
				log.Printf("start %s\n", c.Args().First())
				consume := fmt.Sprintf("core-api-consume:%s", utils.NewID())
				if err := worker.WorkerExecute(c.Args().First(), consume, 12); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
		{
			Name:  "server",
			Usage: "start api server",
			Action: func(c *cli.Context) error {
				log.Printf("start %s\n", c.Command.Name)
				server.Serve()
				return nil
			},
		},
	}

	// Run the CLI app
	client.Run(os.Args)
}
