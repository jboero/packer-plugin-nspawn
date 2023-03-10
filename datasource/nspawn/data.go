//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package nspawn

import (
	"log"
	"os/exec"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
)

type Config struct {
	MockOption string `mapstructure:"mock"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	Machines string `mapstructure:"machines"`
	Images   string `mapstructure:"images"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}
	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	machines, err := exec.Command("machinectl", "-o", "json").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	images, err := exec.Command("machinectl", "list-images", "-o", "json").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	output := DatasourceOutput{
		Machines: string(machines),
		Images:   string(images),
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
