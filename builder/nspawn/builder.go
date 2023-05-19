//go:generate packer-sdc mapstructure-to-hcl2 -type Config
// A sample plugin to provide basic systemd-snpawn builder and datasource support.
// Note this may require root access if machinectl and /var/lib/machines are protected.

package nspawn

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
)

const BuilderId = "nspawn.builder"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Image               string `mapstructure:"image"`
	TmpImage            string `mapstructure:"tmp_image"`
}

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) (generatedVars []string, warnings []string, err error) {
	err = config.Decode(&b.config, &config.DecodeOpts{
		PluginType:  "packer.builder.nspawn",
		Interpolate: true,
	}, raws...)
	if err != nil {
		return nil, nil, err
	}

	// Clone the config image to a new image with the buildName to not alter the original.
	b.config.TmpImage = "packer-" + b.config.Image + "-" + fmt.Sprint(time.Now().Unix())
	log.Println("Cloning nspawn " + b.config.Image + " to " + b.config.TmpImage)
	output, err := exec.Command("machinectl", "clone", b.config.Image, b.config.TmpImage).CombinedOutput()
	if err != nil {
		log.Println("nspawn machinectl start output: " + string(output))
		log.Fatal(err)
		return nil, nil, err
	}

	// Return the placeholder for the generated data that will become available to provisioners and post-processors.
	// If the builder doesn't generate any data, just return an empty slice of string: []string{}
	buildGeneratedData := []string{b.config.TmpImage}
	return buildGeneratedData, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	steps := []multistep.Step{}

	steps = append(steps,
		&StepSayConfig{
			MockConfig: b.config.Image,
		},
		new(commonsteps.StepProvision),
	)

	log.Println("Starting nspawn " + b.config.TmpImage)
	output, err := exec.Command("machinectl", "start", b.config.TmpImage).CombinedOutput()
	if err != nil {
		log.Println("nspawn machinectl start output: " + string(output))
		log.Fatal(err)
		return nil, err.(error)
	}
	// Setup the state bag and initial state for the steps
	state := new(multistep.BasicStateBag)
	state.Put("hook", hook)
	state.Put("ui", ui)

	// Set the value of the generated data that will become available to provisioners.
	// To share the data with post-processors, use the StateData in the artifact.
	state.Put("generated_data", map[string]interface{}{
		"GeneratedMockData": b.config.TmpImage,
	})

	// Run!
	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)

	// If there was an error, return that
	if err, ok := state.GetOk("error"); ok {
		return nil, err.(error)
	}

	artifact := &Artifact{
		// Add the builder generated data to the artifact StateData so that post-processors
		// can access them.
		StateData: map[string]interface{}{"image": b.config.TmpImage},
	}

	log.Println("Stopping nspawn " + b.config.TmpImage)
	output, err = exec.Command("machinectl", "stop", b.config.TmpImage).CombinedOutput()
	if err != nil {
		log.Println("nspawn machinectl start output: " + string(output))
		log.Fatal(err)
		return nil, err.(error)
	}

	return artifact, nil
}
