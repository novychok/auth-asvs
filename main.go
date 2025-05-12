package main

import "github.com/novychok/authasvs/cmd"

func main() {
	cmd.InitServeCommands()
	cmd.InitAuthApiCommands()
	cmd.Execute()
}
