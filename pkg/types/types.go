package types


import _ "github.com/BurntSushi/toml"

type JobsConfig struct {
	Command string `toml:"command" json:"command" yaml:"command" env-required:"true"`
	Needs []string `toml:"needs" yaml:"needs" json:"needs" env-required:"true"`
}

type Jobs struct {
	Title string `toml:"title" yaml:"title" json:"title" env-required:"true"`
	Version string `toml:"version" json:"version" yaml:"version" env-required:"true"`
	Jobs map[string]JobsConfig `toml:"jobs" json:"jobs" yaml:"jobs" env-required:"true"`
}


type ExecReq struct {
	RunnerId string
	ProjectPath string
	Image string
	cmd string
}