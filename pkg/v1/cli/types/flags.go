package types

type flag struct {
	Name      string
	Shorthand string
	Usage     string
}

type BoolFlag struct {
	flag
	Value bool
}

type StringFlag struct {
	flag
	Value string
}

func NewBoolFlag(name, shorthand, usage string, value bool) BoolFlag {
	return BoolFlag{
		flag: flag{
			Name:      name,
			Shorthand: shorthand,
			Usage:     usage,
		},
		Value: value,
	}
}

func NewStringFlag(name, shorthand, usage, value string) StringFlag {
	return StringFlag{
		flag: flag{
			Name:      name,
			Shorthand: shorthand,
			Usage:     usage,
		},
		Value: value,
	}
}
