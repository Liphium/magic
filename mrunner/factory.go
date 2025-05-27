package mrunner

type Factory struct {
	mDir string // Magic directory as a base directory
}

// Create a new factory from the magic directory.
func NewFactory(mDir string) Factory {
	return Factory{
		mDir: mDir,
	}
}
