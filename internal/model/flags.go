package model

type Flags struct {
	DebugMode   bool
	ShowVersion bool
	ShowHelp    bool

	Type   string
	Vision bool
	Port   int
	Host   string
}
