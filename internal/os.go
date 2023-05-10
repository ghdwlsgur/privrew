package internal

import (
	"strings"
)

type OS struct {
	Name string
	Arch string
}

func (os *OS) GetName() string {
	return strings.Title(os.Name)
}

func (os *OS) GetArch() string {
	return strings.Split(os.Arch, "_")[0]
}

func (os *OS) GetExtension() string {
	return "tar.gz"
}
