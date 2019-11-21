package config

type Config struct {
	Id             string
	Pid            int
	PidPath        string
	ImagePath      string
	ContainerPath  string
	ApplicationCmd string
	MountType      string
	MountSource    string
	CloneFlags     int32
}


type State struct {
	Pid string
	root string
}