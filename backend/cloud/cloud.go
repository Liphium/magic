package cloud

type ServerType = string

const (
	ServerType2CPU ServerType = "cpu-2" // 2 vCPU
	ServerType4CPU ServerType = "cpu-4" // 4 vCPU
	ServerType8CPU ServerType = "cpu-8" // 8 vCPU
)

// Initialize all cloud integrations
func Initialize() {
	initHetzner()
}
