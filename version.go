package core

import "os"

const Version = "v1.1.1"

func ShowVersion() string {
	os.Setenv("version.goplume.version", Version)
	return Version
}
