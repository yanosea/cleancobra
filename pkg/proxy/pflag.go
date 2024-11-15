package proxy

import (
	"github.com/spf13/pflag"
)

type FlagSet interface {
	StringVarP(p *string, name string, shorthand string, value string, usage string)
}

type flagSetProxy struct {
	FlagSet *pflag.FlagSet
}

func (f *flagSetProxy) StringVarP(p *string, name string, shorthand string, value string, usage string) {
	f.FlagSet.StringVarP(p, name, shorthand, value, usage)
}
