//go:generate mockgen -source=cobra.go -destination=cobra_mock.go -package=proxy

package proxy

import (
	"github.com/spf13/cobra"
)

// Cobra provides a proxy interface for cobra package functionality
type Cobra interface {
	NewCommand() Command
	ExactArgs(n int) cobra.PositionalArgs
	NoArgs() cobra.PositionalArgs
	ArbitraryArgs() cobra.PositionalArgs
}

// Command provides a proxy interface for cobra.Command
type Command interface {
	SetUse(use string)
	SetShort(short string)
	SetLong(long string)
	SetArgs(args cobra.PositionalArgs)
	SetRunE(runE func(cmd *cobra.Command, args []string) error)
	SetSilenceErrors(silenceErrors bool)
	SetSilenceUsage(silenceUsage bool)
	AddCommand(cmds ...Command)
	Execute() error
	Flags() FlagSet
	PersistentFlags() FlagSet
	GenBashCompletion(w interface{}) error
	GenZshCompletion(w interface{}) error
	GenFishCompletion(w interface{}, includeDesc bool) error
	GenPowerShellCompletion(w interface{}) error
}

// FlagSet provides a proxy interface for cobra flag functionality
type FlagSet interface {
	StringP(name, shorthand string, value string, usage string) *string
	String(name string, value string, usage string) *string
	BoolP(name, shorthand string, value bool, usage string) *bool
	Bool(name string, value bool, usage string) *bool
	IntP(name, shorthand string, value int, usage string) *int
	Int(name string, value int, usage string) *int
}

// CobraImpl implements the Cobra interface using the cobra package
type CobraImpl struct{}

// CommandImpl implements the Command interface wrapping cobra.Command
type CommandImpl struct {
	cmd *cobra.Command
}

// FlagSetImpl implements the FlagSet interface wrapping cobra flag functionality
type FlagSetImpl struct {
	flags interface {
		StringP(name, shorthand string, value string, usage string) *string
		String(name string, value string, usage string) *string
		BoolP(name, shorthand string, value bool, usage string) *bool
		Bool(name string, value bool, usage string) *bool
		IntP(name, shorthand string, value int, usage string) *int
		Int(name string, value int, usage string) *int
	}
}

// NewCobra creates a new Cobra implementation
func NewCobra() Cobra {
	return &CobraImpl{}
}

func (c *CobraImpl) NewCommand() Command {
	return &CommandImpl{cmd: &cobra.Command{}}
}

func (c *CobraImpl) ExactArgs(n int) cobra.PositionalArgs {
	return cobra.ExactArgs(n)
}

func (c *CobraImpl) NoArgs() cobra.PositionalArgs {
	return cobra.NoArgs
}

func (c *CobraImpl) ArbitraryArgs() cobra.PositionalArgs {
	return cobra.ArbitraryArgs
}

func (cmd *CommandImpl) SetUse(use string) {
	cmd.cmd.Use = use
}

func (cmd *CommandImpl) SetShort(short string) {
	cmd.cmd.Short = short
}

func (cmd *CommandImpl) SetLong(long string) {
	cmd.cmd.Long = long
}

func (cmd *CommandImpl) SetArgs(args cobra.PositionalArgs) {
	cmd.cmd.Args = args
}

func (cmd *CommandImpl) SetRunE(runE func(c *cobra.Command, args []string) error) {
	cmd.cmd.RunE = runE
}

func (cmd *CommandImpl) SetSilenceErrors(silenceErrors bool) {
	cmd.cmd.SilenceErrors = silenceErrors
}

func (cmd *CommandImpl) SetSilenceUsage(silenceUsage bool) {
	cmd.cmd.SilenceUsage = silenceUsage
}

func (cmd *CommandImpl) AddCommand(cmds ...Command) {
	cobraCmds := make([]*cobra.Command, len(cmds))
	for i, c := range cmds {
		if impl, ok := c.(*CommandImpl); ok {
			cobraCmds[i] = impl.cmd
		}
	}
	cmd.cmd.AddCommand(cobraCmds...)
}

func (cmd *CommandImpl) Execute() error {
	return cmd.cmd.Execute()
}

func (cmd *CommandImpl) Flags() FlagSet {
	return &FlagSetImpl{flags: cmd.cmd.Flags()}
}

func (cmd *CommandImpl) PersistentFlags() FlagSet {
	return &FlagSetImpl{flags: cmd.cmd.PersistentFlags()}
}

func (cmd *CommandImpl) GenBashCompletion(w interface{}) error {
	if writer, ok := w.(interface{ Write([]byte) (int, error) }); ok {
		return cmd.cmd.GenBashCompletion(writer)
	}
	return nil
}

func (cmd *CommandImpl) GenZshCompletion(w interface{}) error {
	if writer, ok := w.(interface{ Write([]byte) (int, error) }); ok {
		return cmd.cmd.GenZshCompletion(writer)
	}
	return nil
}

func (cmd *CommandImpl) GenFishCompletion(w interface{}, includeDesc bool) error {
	if writer, ok := w.(interface{ Write([]byte) (int, error) }); ok {
		return cmd.cmd.GenFishCompletion(writer, includeDesc)
	}
	return nil
}

func (cmd *CommandImpl) GenPowerShellCompletion(w interface{}) error {
	if writer, ok := w.(interface{ Write([]byte) (int, error) }); ok {
		return cmd.cmd.GenPowerShellCompletion(writer)
	}
	return nil
}

func (f *FlagSetImpl) StringP(name, shorthand string, value string, usage string) *string {
	return f.flags.StringP(name, shorthand, value, usage)
}

func (f *FlagSetImpl) String(name string, value string, usage string) *string {
	return f.flags.String(name, value, usage)
}

func (f *FlagSetImpl) BoolP(name, shorthand string, value bool, usage string) *bool {
	return f.flags.BoolP(name, shorthand, value, usage)
}

func (f *FlagSetImpl) Bool(name string, value bool, usage string) *bool {
	return f.flags.Bool(name, value, usage)
}

func (f *FlagSetImpl) IntP(name, shorthand string, value int, usage string) *int {
	return f.flags.IntP(name, shorthand, value, usage)
}

func (f *FlagSetImpl) Int(name string, value int, usage string) *int {
	return f.flags.Int(name, value, usage)
}