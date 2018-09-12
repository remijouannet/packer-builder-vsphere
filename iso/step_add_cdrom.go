package iso

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
)

type CDRomConfig struct {
	ISOPaths []string `mapstructure:"iso_paths"`
}

type StepAddCDRom struct {
	RemotePath string
	Config     *CDRomConfig
}

func (s *StepAddCDRom) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("vm").(*driver.VirtualMachine)

	ui.Say("Adding CD-ROM drives...")
	if err := vm.AddSATAController(); err != nil {
		state.Put("error", fmt.Errorf("error adding SATA controller: %v", err))
		return multistep.ActionHalt
	}

	if path, ok := state.GetOk(s.RemotePath); ok {
		if err := vm.AddCdrom(path.(string)); err != nil {
			state.Put("error", fmt.Errorf("error adding a cdrom: %v", err))
			return multistep.ActionHalt
		}
	} else {
		for _, path := range s.Config.ISOPaths {
			if err := vm.AddCdrom(path); err != nil {
				state.Put("error", fmt.Errorf("error adding a cdrom: %v", err))
				return multistep.ActionHalt
			}
		}
	}

	return multistep.ActionContinue
}

func (s *StepAddCDRom) Cleanup(state multistep.StateBag) {}
