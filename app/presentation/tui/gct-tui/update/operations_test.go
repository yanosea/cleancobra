package update

import (
	"errors"
	"testing"
	"time"

	"github.com/yanosea/gct/app/application"
	"github.com/yanosea/gct/app/domain"
	"github.com/yanosea/gct/app/presentation/tui/gct-tui/model"
	"github.com/yanosea/gct/pkg/proxy"
	"go.uber.org/mock/gomock"
)

func createTestStateModelForOps(t *testing.T) (*model.StateModel, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	
	// Create mock use cases
	mockAddUseCase := &application.AddTodoUseCase{}
	mockListUseCase := &application.ListTodoUseCase{}
	mockToggleUseCase := &application.ToggleTodoUseCase{}
	mockDeleteUseCase := &application.DeleteTodoUseCase{}
	
	// Create mock bubbles
	mockBubbles := proxy.NewMockBubbles(ctrl)
	mockTextInput := proxy.NewMockTextInput(ctrl)
	
	mockBubbles.EXPECT().NewTextInput().Return(mockTextInput)
	mockTextInput.EXPECT().SetPlaceholder("Enter todo description...")
	mockTextInput.EXPECT().SetCharLimit(500)
	mockTextInput.EXPECT().SetWidth(50)
	
	// Allow common input operations that might be called during mode changes
	mockTextInput.EXPECT().Blur().AnyTimes()
	mockTextInput.EXPECT().SetValue(gomock.Any()).AnyTimes()
	mockTextInput.EXPECT().Value().Return("").AnyTimes()
	mockTextInput.EXPECT().Focus().Return(nil).AnyTimes()
	
	stateModel := model.NewStateModel(
		mockAddUseCase,
		mockListUseCase,
		mockToggleUseCase,
		mockDeleteUseCase,
		mockBubbles,
	)
	
	return stateModel, ctrl
}

func TestNewOperationsHandler(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewOperationsHandler()
			
			if handler == nil {
				t.Error("NewOperationsHandler() returned nil")
			}
		})
	}
}

func TestOperationsHandler_HandleTodosLoaded(t *testing.T) {
	tests := []struct {
		name  string
		todos []*model.ItemModel
	}{
		{
			name:  "positive testing (empty todos)",
			todos: []*model.ItemModel{},
		},
		{
			name: "positive testing (with todos)",
			todos: []*model.ItemModel{
				model.NewItemModel(&domain.Todo{
					ID:          1,
					Description: "Test todo",
					Done:        false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModelForOps(t)
			defer ctrl.Finish()
			
			handler := NewOperationsHandler()
			msg := model.TodosLoadedMsg{Todos: tt.todos}
			
			result := handler.HandleTodosLoaded(stateModel, msg)
			
			if result.Model == nil {
				t.Error("HandleTodosLoaded() returned nil model")
			}
			
			if len(result.Model.Todos()) != len(tt.todos) {
				t.Errorf("HandleTodosLoaded() todos count = %v, want %v", len(result.Model.Todos()), len(tt.todos))
			}
			
			if len(tt.todos) > 0 && result.Model.Cursor() != 0 {
				t.Errorf("HandleTodosLoaded() cursor = %v, want %v", result.Model.Cursor(), 0)
			}
		})
	}
}

func TestOperationsHandler_HandleTodoAdded(t *testing.T) {
	tests := []struct {
		name string
		todo *domain.Todo
	}{
		{
			name: "positive testing",
			todo: &domain.Todo{
				ID:          1,
				Description: "New todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModelForOps(t)
			defer ctrl.Finish()
			
			handler := NewOperationsHandler()
			msg := model.TodoAddedMsg{Todo: tt.todo}
			
			result := handler.HandleTodoAdded(stateModel, msg)
			
			if result.Model == nil {
				t.Error("HandleTodoAdded() returned nil model")
			}
			
			if result.Cmd == nil {
				t.Error("HandleTodoAdded() should return a command to reload todos")
			}
		})
	}
}

func TestOperationsHandler_HandleTodoToggled(t *testing.T) {
	tests := []struct {
		name string
		todo *domain.Todo
	}{
		{
			name: "positive testing",
			todo: &domain.Todo{
				ID:          1,
				Description: "Toggled todo",
				Done:        true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModelForOps(t)
			defer ctrl.Finish()
			
			// Set up existing todo
			existingTodo := &domain.Todo{
				ID:          1,
				Description: "Toggled todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			itemModel := model.NewItemModel(existingTodo)
			stateModel.SetTodos([]*model.ItemModel{itemModel})
			
			handler := NewOperationsHandler()
			msg := model.TodoToggledMsg{Todo: tt.todo}
			
			result := handler.HandleTodoToggled(stateModel, msg)
			
			if result.Model == nil {
				t.Error("HandleTodoToggled() returned nil model")
			}
			
			// Check that the todo was updated
			if result.Model.Todos()[0].Todo().Done != tt.todo.Done {
				t.Errorf("HandleTodoToggled() todo done = %v, want %v", result.Model.Todos()[0].Todo().Done, tt.todo.Done)
			}
		})
	}
}

func TestOperationsHandler_HandleTodoDeleted(t *testing.T) {
	tests := []struct {
		name string
		id   int
	}{
		{
			name: "positive testing",
			id:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModelForOps(t)
			defer ctrl.Finish()
			
			handler := NewOperationsHandler()
			msg := model.TodoDeletedMsg{ID: tt.id}
			
			result := handler.HandleTodoDeleted(stateModel, msg)
			
			if result.Model == nil {
				t.Error("HandleTodoDeleted() returned nil model")
			}
			
			if result.Cmd == nil {
				t.Error("HandleTodoDeleted() should return a command to reload todos")
			}
		})
	}
}

func TestOperationsHandler_HandleTodoUpdated(t *testing.T) {
	tests := []struct {
		name string
		todo *domain.Todo
	}{
		{
			name: "positive testing",
			todo: &domain.Todo{
				ID:          1,
				Description: "Updated todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModelForOps(t)
			defer ctrl.Finish()
			
			// Set up existing todo
			existingTodo := &domain.Todo{
				ID:          1,
				Description: "Old description",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			itemModel := model.NewItemModel(existingTodo)
			stateModel.SetTodos([]*model.ItemModel{itemModel})
			
			handler := NewOperationsHandler()
			msg := model.TodoUpdatedMsg{Todo: tt.todo}
			
			result := handler.HandleTodoUpdated(stateModel, msg)
			
			if result.Model == nil {
				t.Error("HandleTodoUpdated() returned nil model")
			}
			
			if result.Model.Mode() != model.ModeNormal {
				t.Errorf("HandleTodoUpdated() mode = %v, want %v", result.Model.Mode(), model.ModeNormal)
			}
			
			// Check that the todo was updated
			if result.Model.Todos()[0].Todo().Description != tt.todo.Description {
				t.Errorf("HandleTodoUpdated() todo description = %v, want %v", result.Model.Todos()[0].Todo().Description, tt.todo.Description)
			}
		})
	}
}

func TestOperationsHandler_HandleError(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "positive testing",
			err:  errors.New("test error"),
		},
		{
			name: "positive testing (nil error)",
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModelForOps(t)
			defer ctrl.Finish()
			
			stateModel.SetMode(model.ModeInput) // Set to non-normal mode
			
			handler := NewOperationsHandler()
			msg := model.ErrorMsg{Error: tt.err}
			
			result := handler.HandleError(stateModel, msg)
			
			if result.Model == nil {
				t.Error("HandleError() returned nil model")
			}
			
			if result.Model.Mode() != model.ModeNormal {
				t.Errorf("HandleError() mode = %v, want %v", result.Model.Mode(), model.ModeNormal)
			}
			
			if tt.err != nil && result.Model.ErrorMessage() == "" {
				t.Error("HandleError() should set error message")
			}
		})
	}
}

func TestOperationsHandler_HandleItemAdd(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "positive testing",
			description: "New item",
		},
		{
			name:        "positive testing (empty description)",
			description: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModelForOps(t)
			defer ctrl.Finish()
			
			handler := NewOperationsHandler()
			msg := ItemAddMsg{Description: tt.description}
			
			result := handler.HandleItemAdd(stateModel, msg)
			
			if result.Model == nil {
				t.Error("HandleItemAdd() returned nil model")
			}
			
			if result.Cmd == nil {
				t.Error("HandleItemAdd() should return a command")
			}
		})
	}
}

func TestOperationsHandler_HandleItemDelete(t *testing.T) {
	tests := []struct {
		name string
		id   int
	}{
		{
			name: "positive testing",
			id:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModelForOps(t)
			defer ctrl.Finish()
			
			// Set up existing todo
			existingTodo := &domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			itemModel := model.NewItemModel(existingTodo)
			stateModel.SetTodos([]*model.ItemModel{itemModel})
			
			handler := NewOperationsHandler()
			msg := ItemDeleteMsg{ID: tt.id}
			
			result := handler.HandleItemDelete(stateModel, msg)
			
			if result.Model == nil {
				t.Error("HandleItemDelete() returned nil model")
			}
			
			if result.Model.ConfirmationMessage() == "" {
				t.Error("HandleItemDelete() should set confirmation message")
			}
		})
	}
}

func TestOperationsHandler_HandleItemToggleAsync(t *testing.T) {
	tests := []struct {
		name string
		id   int
	}{
		{
			name: "positive testing",
			id:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModelForOps(t)
			defer ctrl.Finish()
			
			// Set up existing todo
			existingTodo := &domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			itemModel := model.NewItemModel(existingTodo)
			stateModel.SetTodos([]*model.ItemModel{itemModel})
			
			handler := NewOperationsHandler()
			msg := ItemToggleAsyncMsg{ID: tt.id}
			
			result := handler.HandleItemToggleAsync(stateModel, msg)
			
			if result.Model == nil {
				t.Error("HandleItemToggleAsync() returned nil model")
			}
			
			if result.Cmd == nil {
				t.Error("HandleItemToggleAsync() should return a command")
			}
		})
	}
}

func TestOperationsHandler_HandleItemEditState(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		editing bool
	}{
		{
			name:    "positive testing (start editing)",
			id:      1,
			editing: true,
		},
		{
			name:    "positive testing (stop editing)",
			id:      1,
			editing: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModelForOps(t)
			defer ctrl.Finish()
			
			// Set up existing todo
			existingTodo := &domain.Todo{
				ID:          1,
				Description: "Test todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			itemModel := model.NewItemModel(existingTodo)
			stateModel.SetTodos([]*model.ItemModel{itemModel})
			
			handler := NewOperationsHandler()
			msg := ItemEditStateMsg{ID: tt.id, Editing: tt.editing}
			
			result := handler.HandleItemEditState(stateModel, msg)
			
			if result.Model == nil {
				t.Error("HandleItemEditState() returned nil model")
			}
			
			expectedMode := model.ModeNormal
			if tt.editing {
				expectedMode = model.ModeEdit
			}
			
			if result.Model.Mode() != expectedMode {
				t.Errorf("HandleItemEditState() mode = %v, want %v", result.Model.Mode(), expectedMode)
			}
		})
	}
}

func TestOperationsHandler_HandleItemUpdatedAsync(t *testing.T) {
	tests := []struct {
		name string
		todo *domain.Todo
	}{
		{
			name: "positive testing",
			todo: &domain.Todo{
				ID:          1,
				Description: "Updated todo",
				Done:        false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateModel, ctrl := createTestStateModelForOps(t)
			defer ctrl.Finish()
			
			handler := NewOperationsHandler()
			msg := ItemUpdatedAsyncMsg{Todo: tt.todo}
			
			result := handler.HandleItemUpdatedAsync(stateModel, msg)
			
			if result.Model == nil {
				t.Error("HandleItemUpdatedAsync() returned nil model")
			}
			
			if result.Cmd != nil {
				t.Error("HandleItemUpdatedAsync() should not return a command")
			}
		})
	}
}

func TestCreateModeTransitionCommand(t *testing.T) {
	tests := []struct {
		name string
		mode model.Mode
	}{
		{
			name: "positive testing (normal mode)",
			mode: model.ModeNormal,
		},
		{
			name: "positive testing (input mode)",
			mode: model.ModeInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CreateModeTransitionCommand(tt.mode)
			
			if cmd == nil {
				t.Error("CreateModeTransitionCommand() returned nil")
			}
			
			msg := cmd()
			modeMsg, ok := msg.(ModeTransitionMsg)
			if !ok {
				t.Error("CreateModeTransitionCommand() should return ModeTransitionMsg")
			}
			
			if modeMsg.Mode != tt.mode {
				t.Errorf("CreateModeTransitionCommand() mode = %v, want %v", modeMsg.Mode, tt.mode)
			}
		})
	}
}

func TestCreateNavigationCommand(t *testing.T) {
	tests := []struct {
		name   string
		action NavigationAction
		value  int
	}{
		{
			name:   "positive testing (up)",
			action: NavigationUp,
			value:  0,
		},
		{
			name:   "positive testing (down)",
			action: NavigationDown,
			value:  0,
		},
		{
			name:   "positive testing (to position)",
			action: NavigationTo,
			value:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CreateNavigationCommand(tt.action, tt.value)
			
			if cmd == nil {
				t.Error("CreateNavigationCommand() returned nil")
			}
			
			msg := cmd()
			navMsg, ok := msg.(NavigationMsg)
			if !ok {
				t.Error("CreateNavigationCommand() should return NavigationMsg")
			}
			
			if navMsg.Action != tt.action {
				t.Errorf("CreateNavigationCommand() action = %v, want %v", navMsg.Action, tt.action)
			}
			
			if navMsg.Value != tt.value {
				t.Errorf("CreateNavigationCommand() value = %v, want %v", navMsg.Value, tt.value)
			}
		})
	}
}

func TestCreateAsyncOperationCommand(t *testing.T) {
	tests := []struct {
		name      string
		operation AsyncOperation
		data      interface{}
	}{
		{
			name:      "positive testing (add)",
			operation: AsyncOperationAdd,
			data:      "test data",
		},
		{
			name:      "positive testing (toggle)",
			operation: AsyncOperationToggle,
			data:      123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CreateAsyncOperationCommand(tt.operation, tt.data)
			
			if cmd == nil {
				t.Error("CreateAsyncOperationCommand() returned nil")
			}
			
			msg := cmd()
			asyncMsg, ok := msg.(AsyncOperationMsg)
			if !ok {
				t.Error("CreateAsyncOperationCommand() should return AsyncOperationMsg")
			}
			
			if asyncMsg.Operation != tt.operation {
				t.Errorf("CreateAsyncOperationCommand() operation = %v, want %v", asyncMsg.Operation, tt.operation)
			}
			
			if asyncMsg.Data != tt.data {
				t.Errorf("CreateAsyncOperationCommand() data = %v, want %v", asyncMsg.Data, tt.data)
			}
		})
	}
}