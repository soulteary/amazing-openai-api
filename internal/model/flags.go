package model

type Flags struct {
	DebugMode   bool
	ShowVersion bool
	ShowHelp    bool

	Type string
	Port int
	Host string
}
