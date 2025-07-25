package command

import (
	"fmt"
	"os"

	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	"github.com/yanosea/gct/pkg/proxy"
)

type Runner struct {
	bubbletea proxy.Bubbletea
	usecases  *model.Usecases
}

func NewRootRunner(bubbletea proxy.Bubbletea, usecases *model.Usecases) *Runner {
	return &Runner{
		bubbletea: bubbletea,
		usecases:  usecases,
	}
}

func (r *Runner) Run() int {
	m := model.NewModel(r.usecases)
	program := r.bubbletea.NewProgram(m)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		return 1
	}

	return 0
}
