package tests

import "os/user"

var (
	// IsRoot is true if the current process is run as root
	IsRoot = false
	//
	// DefaultFalcoExecutable is the default path of the Falco executable
	DefaultFalcoExecutable = "/usr/bin/falco"
)

func init() {
	user, err := user.Current()
	if err != nil {
		panic("could not get the current user")
	}
	IsRoot = user.Username == "root"
}
