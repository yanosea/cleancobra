package model

import todoApp "github.com/yanosea/gct/app/application/todo"

type TodosLoadedMsg struct {
	Todos []*todoApp.ListTodoUsecaseOutputDto
}

type ErrorMsg struct {
	Error string
}

type SuccessMsg struct {
	Message string
}

type WindowSizeMsg struct {
	Width  int
	Height int
}

