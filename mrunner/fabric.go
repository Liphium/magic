package mrunner

type Fabric struct {
	mDir string // Magic directory as a base directory
}

// Create a new fabric from the magic directory.
func NewFabric(mDir string) Fabric {
	return Fabric{
		mDir: mDir,
	}
}
