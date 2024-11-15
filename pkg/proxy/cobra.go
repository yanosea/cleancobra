package proxy

import (
	"io"

	"github.com/spf13/cobra"
)

type Cobra interface {
	ExactArgs(int) PositionalArgs
	MaximumNArgs(int) PositionalArgs
	NewCommand() Command
}

type cobraProxy struct{}

func NewCobra() Cobra {
	return &cobraProxy{}
}

func (*cobraProxy) ExactArgs(n int) PositionalArgs {
	return &positionalArgsProxy{PositionalArgs: cobra.ExactArgs(n)}
}

func (*cobraProxy) MaximumNArgs(n int) PositionalArgs {
	return &positionalArgsProxy{PositionalArgs: cobra.MaximumNArgs(n)}
}

func (*cobraProxy) NewCommand() Command {
	return &commandProxy{Command: &cobra.Command{}}
}

type PositionalArgs interface {
	GetPositionalArgs() cobra.PositionalArgs
}

type positionalArgsProxy struct {
	PositionalArgs cobra.PositionalArgs
}

func (p *positionalArgsProxy) GetPositionalArgs() cobra.PositionalArgs {
	return p.PositionalArgs
}

type Command interface {
	AddCommand(cmds ...Command)
	Execute() error
	GetCommand() *cobra.Command
	PersistentFlags() FlagSet
	RunE(cmd *cobra.Command, args []string) error
	SetArgs(positionalArgs PositionalArgs)
	SetErr(io io.Writer)
	SetHelpTemplate(s string)
	SetOut(io io.Writer)
	SetRunE(runE func(cmd *cobra.Command, args []string) error)
	SetShort(short string)
	SetSilenceErrors(silenceErrors bool)
	SetUse(use string)
}

type commandProxy struct {
	Command *cobra.Command
}

func (c *commandProxy) AddCommand(cmds ...Command) {
	for _, cmd := range cmds {
		c.Command.AddCommand(cmd.GetCommand())
	}
}

func (c *commandProxy) Execute() error {
	return c.Command.Execute()
}

func (c *commandProxy) GetCommand() *cobra.Command {
	return c.Command
}

func (c *commandProxy) PersistentFlags() FlagSet {
	return &flagSetProxy{FlagSet: c.Command.PersistentFlags()}
}

func (c *commandProxy) RunE(cmd *cobra.Command, args []string) error {
	return c.Command.RunE(cmd, args)
}

func (c *commandProxy) SetArgs(positionalArgs PositionalArgs) {
	c.Command.Args = positionalArgs.GetPositionalArgs()
}

func (c *commandProxy) SetErr(io io.Writer) {
	c.Command.SetErr(io)
}

func (c *commandProxy) SetHelpTemplate(s string) {
	c.Command.SetHelpTemplate(s)
}

func (c *commandProxy) SetOut(io io.Writer) {
	c.Command.SetOut(io)
}

func (c *commandProxy) SetRunE(runE func(cmd *cobra.Command, args []string) error) {
	c.Command.RunE = runE
}

func (c *commandProxy) SetShort(short string) {
	c.Command.Short = short
}

func (c *commandProxy) SetSilenceErrors(silenceErrors bool) {
	c.Command.SilenceErrors = silenceErrors
}

func (c *commandProxy) SetUse(use string) {
	c.Command.Use = use
}
