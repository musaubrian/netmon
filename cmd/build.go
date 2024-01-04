package main

import "github.com/musaubrian/netmon/gno"

func main() {
	g := gno.New()
	g.BootstrapBuild("build", "netmon", ".")
	g.CopyResources("logs")
	g.CopyResources("templates")
	g.CopyResources("web")
	g.CopyResources(".env")
	g.CopyResources("config.yml")
	g.AddCommand("./build/netmon")
	g.Build()
	g.RunCommands()
}
