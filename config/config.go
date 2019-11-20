package config

type Config struct {
	Id string
	ImagePath string
	ContainerPath string
	ApplicationCmd string
	MountType string
	MountSource string
}


type State struct {
	Pid string
	root string
}