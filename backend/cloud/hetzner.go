package cloud

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Liphium/magic/backend/util"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type HetznerIntegration struct {
	client *hcloud.Client
	env    string
}

var Hetzner *HetznerIntegration

// All the server types mapped to their hetzner equivalent (for now shared amd cores)
var hetznerTypeMappings = map[ServerType]string{
	ServerType2CPU: "cpx11",
	ServerType4CPU: "cpx31",
	ServerType8CPU: "cpx41",
}

// Mappings to the actual type on the API
var hetznerMappings = map[ServerType]*hcloud.ServerType{}

// Initialize the Hetzner integration
func initHetzner() {
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("MAGIC_HETZNER")))
	Hetzner = &HetznerIntegration{
		client: client,
		env:    os.Getenv("MAGIC_HETZNER_ENV"),
	}

	// Get all of the hetzner server types
	types, err := client.ServerType.All(context.Background())
	if err != nil {
		log.Fatalln("Couldn't get server types for Hetzner:", err)
	}

	// Parse the server types into hetznerMappings
	for _, hst := range types {
		for st, name := range hetznerTypeMappings {
			if name == hst.Name {
				hetznerMappings[st] = hst
			}
		}
	}

	// Make sure all mappings have been set
	for st, name := range hetznerTypeMappings {
		if _, ok := hetznerMappings[st]; !ok {
			log.Fatalln("Type", st, "with Hetzner name ", name, "is missing!")
		}
	}
}

// Cloud-Init script for Forge
const hetznerForgeCloudInit = `
runcmd:
  - ./home/spellcast forge %s 
`

// Deploy a new instance of spellcast on a cloud server (specifically for Forge)
func (h *HetznerIntegration) DeployForge(name string, serverType ServerType, account string, spellcastToken string) error {

	// Get the latest spellcast snapshot
	images, err := h.client.Image.AllWithOpts(context.Background(), hcloud.ImageListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: "spellcast-forge-latest,env=" + h.env,
		},
	})
	if err != nil {
		return err
	}

	// Make sure there is only one image with the latest version of spellcast
	if len(images) != 1 {
		return fmt.Errorf("too many or too less latest images for spellcast: %d", len(images))
	}

	// Deploy a new cloud server with this image
	result, _, err := h.client.Server.Create(context.Background(), hcloud.ServerCreateOpts{
		Name:       name,
		ServerType: hetznerMappings[serverType],
		Image:      images[0],
		Labels: map[string]string{
			"env":     h.env,
			"account": account,
			"forge":   "true",
		},
		StartAfterCreate: hcloud.Ptr(true),
		UserData:         fmt.Sprintf(hetznerForgeCloudInit, spellcastToken),
	})
	if err != nil {
		return fmt.Errorf("couldn't start server: %e", err)
	}

	// Print the password in case of testing mode
	if util.IsTesting() {
		log.Println("Started server", result.Server.Name, "with password:", result.RootPassword)
	}

	return nil
}
