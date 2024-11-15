package main

import (
	"os"

	"cleancobra/app/presentation/cli/adapter"
	"cleancobra/app/presentation/cli/adapter/controller/todo"
)

func main() {
	g := todo.NewGlobalOption(adapter.NewTodoUseCases())
	os.Exit(g.Execute())
}
