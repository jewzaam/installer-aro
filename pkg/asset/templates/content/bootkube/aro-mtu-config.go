package bootkube

import (
	"fmt"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
)

const (
	ignFilePath = "/etc/NetworkManager/dispatcher.d/30-eth0-mtu-3900"
	ignFileData = `#!/bin/bash
if [ "$1" == "eth0" ] && [ "$2" == "up" ]; then
    ip link set $1 mtu %d
fi`
)

var _ asset.WritableAsset = (*AROMTUConfig)(nil)

// AROMTUConfig overrides the eth0 device MTU on nodes
type AROMTUConfig struct {
	MTU int
}

// Dependencies returns all of the dependencies directly needed by the asset
func (c *AROMTUConfig) Dependencies() []asset.Asset {
	return nil
}

// Name returns the human-friendly name of the asset
func (c *AROMTUConfig) Name() string {
	return "AROMTUConfig"
}

// Generate generates the actual files by this asset
func (c *AROMTUConfig) Generate(parents asset.Parents) error {
	return nil
}

func (c *AROMTUConfig) Files() []*asset.File {
	return nil
}

// Load returns the asset from disk
func (c *AROMTUConfig) Load(f asset.FileFetcher) (bool, error) {
	return true, nil
}

// Generates a file for an ignition config
func (c *AROMTUConfig) IgnitionFile() igntypes.File {
	return ignition.FileFromString(ignFilePath, "root", 0555, fmt.Sprintf(ignFileData, c.MTU))
}

// Generates a MachineConfig for the given role
func (c *AROMTUConfig) MachineConfig(role string) (*mcfgv1.MachineConfig, error) {
	ignConfig := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Storage: igntypes.Storage{
			Files: []igntypes.File{
				c.IgnitionFile(),
			},
		},
	}

	rawExt, err := ignition.ConvertToRawExtension(ignConfig)
	if err != nil {
		return nil, err
	}

	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-mtu", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: rawExt,
		},
	}, nil
}
