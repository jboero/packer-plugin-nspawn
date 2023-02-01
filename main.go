package main

import (
	"fmt"
	"os"
	"packer-plugin-nspawn/builder/nspawn"
	nspawnData "packer-plugin-nspawn/datasource/nspawn"
	nspawnVersion "packer-plugin-nspawn/version"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("build", new(nspawn.Builder))
	pps.RegisterDatasource("images", new(nspawnData.Datasource))
	pps.SetVersion(nspawnVersion.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
