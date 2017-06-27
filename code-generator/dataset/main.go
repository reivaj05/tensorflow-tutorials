package main

import (
	"fmt"
	"os"

	"github.com/reivaj05/apigateway/generator"

	"github.com/reivaj05/GoCLI"
	"github.com/reivaj05/GoConfig"
	"github.com/reivaj05/apigateway/server"
)

// TODO: Validate all things needed(GOPATH is set, etc.)

const appName = "apigateway"

func main() {
	setup()
	startApp()
}

func setup() {
	startConfig()
	startLogger()
}

func startConfig() {
	if err := GoConfig.Init(createConfigOptions()); err != nil {
		finishExecution("Error while starting config", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func createConfigOptions() *GoConfig.ConfigOptions {
	return &GoConfig.ConfigOptions{
		ConfigType: "json",
		ConfigFile: "config",
		ConfigPath: ".",
	}
}

func startLogger() {

}

func startApp() {
	if err := GoCLI.StartCLI(createCLIOptions()); err != nil {
		finishExecution("Error while starting application", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func createCLIOptions() *GoCLI.Options {
	return &GoCLI.Options{
		AppName:       appName,
		AppUsage:      "TODO: Set app usage",
		Commands:      createCommands(),
		DefaultAction: server.Serve,
	}
}

func createCommands() []*GoCLI.Command {
	return []*GoCLI.Command{
		&GoCLI.Command{
			Name:   "start",
			Usage:  "TODO: Set start usage",
			Action: server.Serve,
		},
		&GoCLI.Command{
			Name:   "create-service",
			Usage:  "TODO: Set create-service usage",
			Action: generator.Generate,
		},
	}
}

func finishExecution(msg string, fields map[string]interface{}) {
	fmt.Println(msg, fields)
	os.Exit(1)
}
