# Design Document

## Overview

「gct」は、クリーンアーキテクチャの原則に従って設計されたTODOアプリケーションです。アプリケーションは以下の主要な特徴を持ちます：

- **Clean Architecture**: ドメイン、アプリケーション、インフラストラクチャ、プレゼンテーション層の明確な分離
- **Dual Interface**: CLI（Cobra）とTUI（Bubbletea）の両方をサポート
- **JSON Storage**: 軽量で人間が読めるJSON形式でのデータ永続化
- **Environment Configuration**: 環境変数による柔軟な設定
- **High Testability**: 100%テストカバレッジを目指した設計

## Architecture

### Layer Structure

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────────────┐    ┌─────────────────────────────┐ │
│  │        CLI          │    │            TUI              │ │
│  │   (Cobra-based)     │    │    (Bubbletea + ELM)        │ │
│  └─────────────────────┘    └─────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                         │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌────────┐ │
│  │ AddTodo     │ │ DeleteTodo  │ │ ListTodo    │ │ Toggle │ │
│  │ UseCase     │ │ UseCase     │ │ UseCase     │ │ UseCase│ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                            │
│  ┌─────────────────────┐    ┌─────────────────────────────┐ │
│  │      Todo           │    │    TodoRepository           │ │
│  │    (Entity)         │    │    (Interface)              │ │
│  └─────────────────────┘    └─────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Infrastructure Layer                         │
│  ┌─────────────────────┐    ┌─────────────────────────────┐ │
│  │   JSONRepository    │    │     Configuration           │ │
│  │  (Implementation)   │    │      (moved to config)      │ │
│  └─────────────────────┘    └─────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│              Configuration & Composition Root               │
│  ┌─────────────────────┐    ┌─────────────────────────────┐ │
│  │      Config         │    │      Container              │ │
│  │   (app/config)      │    │   (app/container)           │ │
│  └─────────────────────┘    └─────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Project Structure

```
app/
├── config/              # 設定管理（横断的関心事）
├── container/           # DIコンテナ/コンポジションルート
├── domain/              # ドメイン層
├── application/         # アプリケーション層
├── infrastructure/      # インフラ層（永続化、外部API等）
└── presentation/        # プレゼンテーション層
```

### Dependency Flow

- **Presentation Layer** → **Application Layer** → **Domain Layer**
- **Infrastructure Layer** → **Domain Layer** (implements interfaces)
- **Configuration Layer** → 全ての層で利用可能
- **Container Layer** → 全ての依存関係を組み立て

## Components and Interfaces

### Domain Layer

#### Todo Entity

```go
type Todo struct {
    ID          int       `json:"id"`
    Description string    `json:"description"`
    Done        bool      `json:"done"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### TodoRepository Interface

```go
type TodoRepository interface {
    FindAll() ([]Todo, error)
    Save(todos ...Todo) ([]Todo, error)  // 新規なら新規保存、既存なら更新保存
    DeleteById(id int) error
}
```

### Application Layer

#### Use Cases

- **AddTodoUseCase**: 新しいTODOを追加
- **DeleteTodoUseCase**: 指定されたIDのTODOを削除
- **ListTodoUseCase**: 全てのTODOを取得
- **ToggleTodoUseCase**: 指定されたIDのTODOの完了状態を切り替え

各UseCaseは以下のパターンに従います：

```go
type AddTodoUseCase struct {
    repo TodoRepository
}

func (uc *AddTodoUseCase) Run(description string) (*Todo, error)
```

### Infrastructure Layer

#### JSONRepository

- ファイルベースのJSON永続化を実装
- 環境変数による設定ファイルパスの管理
- ファイルロックによる並行アクセス制御
- エラーハンドリングとファイル復旧機能

### Configuration Layer (app/config)

#### Configuration Management

設定は横断的関心事として独立したパッケージで管理：

```go
type Config struct {
    DataFile string `envconfig:"GCT_DATA_FILE"`
}

// 設定の読み込みと検証を担当
func Load() (*Config, error)
func (c *Config) Validate() error
```

### Container Layer (app/container)

#### Dependency Injection Container

全ての依存関係の組み立てを担当：

```go
type Container struct {
    config     *config.Config
    proxies    *Proxies
    repository domain.TodoRepository
    useCases   *UseCases
}

type Proxies struct {
    OS       proxy.OS
    Filepath proxy.Filepath
    JSON     proxy.JSON
    Time     proxy.Time
    // ... other proxies
}

type UseCases struct {
    AddTodo    *application.AddTodoUseCase
    DeleteTodo *application.DeleteTodoUseCase
    ListTodo   *application.ListTodoUseCase
    ToggleTodo *application.ToggleTodoUseCase
}

func NewContainer() (*Container, error)
func (c *Container) GetUseCases() *UseCases
func (c *Container) GetRepository() domain.TodoRepository
```

### Presentation Layer

#### CLI (Cobra-based)

```
gct/
├── main.go                 # エントリーポイント（シンプルな実行部分のみ）
├── commands/
│   ├── command.go         # コマンド初期化ロジック（InitializeCommand関数）
│   ├── command_test.go    # コマンド初期化のテスト
│   ├── root.go            # ルートコマンド（デフォルトでlist実行）
│   └── gct/
│       ├── add.go         # add サブコマンド
│       ├── delete.go      # delete サブコマンド
│       ├── list.go        # list サブコマンド
│       ├── toggle.go      # toggle サブコマンド
│       └── completion/
│           ├── completion.go  # completion 親コマンド
│           ├── bash.go        # bash 補完サブコマンド
│           ├── fish.go        # fish 補完サブコマンド
│           ├── powershell.go  # powershell 補完サブコマンド
│           └── zsh.go         # zsh 補完サブコマンド
├── formatter/
│   ├── json.go           # JSON出力フォーマッター
│   ├── table.go          # テーブル出力フォーマッター
│   └── plain.go          # プレーンテキスト出力フォーマッター
└── presenter/
    ├── presenter.go      # TodoPresenter実装
    └── presenter_test.go # TodoPresenterテスト
```

##### CLI Command Usage Examples

```bash
# ルートコマンド（デフォルトでlist実行）
go run ./app/presentation/cli/gct/main.go

# TODOを追加
go run ./app/presentation/cli/gct/main.go add "Buy groceries"

# 全TODOをリスト表示
go run ./app/presentation/cli/gct/main.go list

# JSON形式でリスト表示
go run ./app/presentation/cli/gct/main.go list --format json

# プレーンテキスト形式でリスト表示
go run ./app/presentation/cli/gct/main.go list --format plain

# TODOの完了状態を切り替え（ID指定）
go run ./app/presentation/cli/gct/main.go toggle 1

# TODOを削除（ID指定）
go run ./app/presentation/cli/gct/main.go delete 1



# ヘルプ表示
go run ./app/presentation/cli/gct/main.go --help
go run ./app/presentation/cli/gct/main.go add --help

# シェル補完生成
go run ./app/presentation/cli/gct/main.go completion bash
go run ./app/presentation/cli/gct/main.go completion zsh
```

##### Command Structure

- **Root Command**: 引数なしで実行した場合は `list` コマンドを実行
- **Add Command**: `gct add <description>` - 新しいTODOを追加
- **List Command**: `gct list [--format json]` - TODOリストを表示
- **Toggle Command**: `gct toggle <id>` - 指定IDのTODOの完了状態を切り替え
- **Delete Command**: `gct delete <id>` - 指定IDのTODOを削除

##### Cobra Proxy Command Implementation Pattern

各コマンドはcobraプロキシを使用して以下のパターンで実装:

```go
func NewAddCommand(
    cobra proxy.Cobra,
    useCase *application.AddTodoUseCase,
    presenter *presenter.TodoPresenter,
) proxy.Command {
    cmd := cobra.NewCommand()
    cmd.SetUse("add <description>")
    cmd.SetShort("Add a new todo")
    cmd.SetArgs(cobra.ExactArgs(1))
    cmd.SetSilenceErrors(true)
    cmd.SetRunE(
        func(_ *c.Command, args []string) error {
            return runAdd(useCase, presenter, args[0])
        },
    )
    return cmd
}

func runAdd(useCase *application.AddTodoUseCase, presenter *presenter.TodoPresenter, description string) error {
    todo, err := useCase.Run(description)
    if err != nil {
        return err
    }
    presenter.ShowAddSuccess(todo)
    return nil
}
```

#### TUI (Bubbletea + ELM Architecture)

```
gct-tui/
├── main.go               # エントリーポイント（シンプルな実行部分のみ）
├── program/
│   └── program.go        # プログラム初期化ロジック（CLI の commands/command.go 相当）
├── model/
│   ├── state.go          # 状態管理の核心部分
│   ├── mode.go           # モード管理
│   ├── navigation.go     # カーソル・ナビゲーション
│   ├── input.go          # 入力処理
│   ├── item.go           # TODOアイテムモデル
│   └── messages.go       # メッセージ定義
├── update/
│   ├── handler.go        # メインハンドラー
│   ├── keyboard.go       # キーボード入力処理
│   ├── operations.go     # 操作処理（追加・削除・切り替え・更新）
│   └── item.go           # アイテム更新ロジック
└── view/
    ├── layout.go         # メインレイアウト
    ├── header.go         # ヘッダー表示
    ├── footer.go         # フッター表示
    ├── list.go           # リスト表示
    └── item.go           # アイテム表示コンポーネント
```

### Proxy Layer

テスタビリティのために標準パッケージとサードパーティパッケージをラップ：

```
pkg/proxy/
├── os.go                # os パッケージのプロキシ
├── filepath.go          # filepath パッケージのプロキシ
├── json.go              # encoding/json パッケージのプロキシ
├── time.go              # time パッケージのプロキシ
├── io.go                # io パッケージのプロキシ
├── fmt.go               # fmt パッケージのプロキシ
├── strings.go           # strings パッケージのプロキシ
├── strconv.go           # strconv パッケージのプロキシ
├── cobra.go             # github.com/spf13/cobra のプロキシ
├── bubbletea.go         # github.com/charmbracelet/bubbletea のプロキシ
├── bubbles.go           # github.com/charmbracelet/bubbles のプロキシ
├── lipgloss.go          # github.com/charmbracelet/lipgloss のプロキシ
├── color.go             # github.com/fatih/color のプロキシ
└── envconfig.go         # github.com/kelseyhightower/envconfig のプロキシ
```

## Data Models

### Todo JSON Structure

```json
{
    "todos": [
        {
            "id": 1,
            "description": "Buy groceries",
            "done": false,
            "created_at": "2025-01-15T10:00:00Z",
            "updated_at": "2025-01-15T10:00:00Z"
        }
    ],
    "next_id": 2
}
```

### Configuration Model

```go
type Config struct {
    DataFile string `envconfig:"GCT_DATA_FILE"`
}

// Default paths:
// 1. $GCT_DATA_FILE (if set)
// 2. $XDG_DATA_HOME/gct/todos.json (if XDG_DATA_HOME is set)
// 3. ~/.local/share/gct/todos.json (fallback)
```

## Error Handling

### Error Types

```go
type TodoError struct {
    Type    ErrorType
    Message string
    Cause   error
}

type ErrorType int

const (
    ErrorTypeNotFound ErrorType = iota
    ErrorTypeInvalidInput
    ErrorTypeFileSystem
    ErrorTypeJSON
    ErrorTypeConfiguration
)
```

### Error Handling Strategy

- **Domain Layer**: ビジネスロジックエラーを定義
- **Application Layer**: UseCaseレベルでのエラーハンドリング
- **Infrastructure Layer**: ファイルシステムとJSONエラーの処理
- **Presentation Layer**: ユーザーフレンドリーなエラーメッセージ

## Testing Strategy

### Unit Testing

- **Domain Layer**: エンティティとビジネスロジックのテスト
- **Application Layer**: UseCaseの単体テスト（モックリポジトリ使用）
- **Infrastructure Layer**: JSONRepository の統合テスト
- **Presentation Layer**: コマンドハンドラーとビューのテスト

### Test Coverage Goals

- app/ ディレクトリ配下で100%のテストカバレッジを達成
- proxy/ ディレクトリ配下のテストは不要
- テーブル駆動テストパターンの採用
- モックとスタブを使用した依存関係の分離

### Testing Tools

- 標準の `testing` パッケージ
- `testify` ライブラリ（アサーション）
- `gomock` と `mockgen` （モック生成）
- カスタムテストヘルパー関数

### Mock Generation Strategy

- **gomock/mockgen** を使用してインターフェースのモックを自動生成
- モックは同じパッケージ内に `*_mock.go` の命名で配置
- モック生成コマンド例:

```bash
//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=domain
//go:generate mockgen -source=os.go -destination=os_mock.go -package=proxy
```

- 各パッケージ内でのモック管理により、テストの可読性と保守性を向上
- CI/CDでのモック生成の自動化

### Testing Conventions

- **テスト形式**: 全てのテストはテーブル駆動テスト（Table-Driven Tests）で実装
- **テスト名命名規則**:
    - **正常系**: `positive testing` で始まる
    - **異常系**: `negative testing (...) failed` の形式
- **テスト構造例**:

```go
func TestAddTodoUseCase_Run(t *testing.T) {
    tests := []struct {
        name        string
        description string
        setupMock   func(*MockTodoRepository)
        want        *Todo
        wantErr     bool
    }{
        {
            name:        "positive testing",
            description: "Buy groceries",
            setupMock:   func(m *MockTodoRepository) { /* setup */ },
            want:        &Todo{ID: 1, Description: "Buy groceries"},
            wantErr:     false,
        },
        {
            name:        "negative testing (empty description failed)",
            description: "",
            setupMock:   func(m *MockTodoRepository) { /* setup */ },
            want:        nil,
            wantErr:     true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Test Structure

テストファイルは対象ファイルと同じディレクトリに配置:

```
app/
├── config/
│   ├── config.go
│   └── config_test.go
├── container/
│   ├── container.go
│   └── container_test.go
├── domain/
│   ├── todo.go
│   ├── todo_test.go
│   ├── repository.go
│   └── repository_test.go
├── application/
│   ├── add_todo_usecase.go
│   ├── add_todo_usecase_test.go
│   ├── delete_todo_usecase.go
│   ├── delete_todo_usecase_test.go
│   ├── list_todo_usecase.go
│   ├── list_todo_usecase_test.go
│   ├── toggle_todo_usecase.go
│   └── toggle_todo_usecase_test.go
├── infrastructure/
│   ├── json_repository.go
│   └── json_repository_test.go
└── presentation/
    ├── cli/
    │   └── gct/
    │       ├── main.go
    │       ├── main_test.go
    │       ├── commands/
    │       │   ├── command.go
    │       │   ├── command_test.go
    │       │   ├── root.go
    │       │   ├── root_test.go
    │       │   └── gct/
    │       │       ├── add.go
    │       │       ├── add_test.go
    │       │       ├── delete.go
    │       │       ├── delete_test.go
    │       │       ├── list.go
    │       │       ├── list_test.go
    │       │       ├── toggle.go
    │       │       ├── toggle_test.go
    │       │       └── completion/
    │       │           ├── completion.go
    │       │           ├── completion_test.go
    │       │           ├── bash.go
    │       │           ├── bash_test.go
    │       │           ├── fish.go
    │       │           ├── fish_test.go
    │       │           ├── powershell.go
    │       │           ├── powershell_test.go
    │       │           ├── zsh.go
    │       │           └── zsh_test.go
    │       ├── formatter/
    │       │   ├── json.go
    │       │   ├── json_test.go
    │       │   ├── table.go
    │       │   ├── table_test.go
    │       │   ├── plain.go
    │       │   └── plain_test.go
    │       └── presenter/
    │           ├── presenter.go
    │           └── presenter_test.go
    └── tui/
        └── gct-tui/
            ├── main.go
            ├── main_test.go
            ├── program/
            │   ├── program.go
            │   └── program_test.go
            ├── model/
            │   ├── state.go
            │   ├── state_test.go
            │   ├── mode.go
            │   ├── mode_test.go
            │   ├── navigation.go
            │   ├── navigation_test.go
            │   ├── input.go
            │   ├── input_test.go
            │   ├── item.go
            │   ├── item_test.go
            │   ├── messages.go
            │   └── messages_test.go
            ├── update/
            │   ├── handler.go
            │   ├── handler_test.go
            │   ├── keyboard.go
            │   ├── keyboard_test.go
            │   ├── operations.go
            │   ├── operations_test.go
            │   ├── item.go
            │   └── item_test.go
            └── view/
                ├── layout.go
                ├── layout_test.go
                ├── header.go
                ├── header_test.go
                ├── footer.go
                ├── footer_test.go
                ├── list.go
                ├── list_test.go
                ├── item.go
                └── item_test.go
```

## TUI ELM Architecture Integration

### Model

```go
type Model struct {
    todos       []domain.Todo
    cursor      int
    selected    map[int]struct{}
    mode        Mode
    input       textinput.Model
    useCase     *application.UseCases
}
```

### Update

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case TodoAddedMsg:
        return m.handleTodoAdded(msg)
    // ... other message handlers
    }
}
```

### View

```go
func (m Model) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.renderHeader(),
        m.renderTodoList(),
        m.renderFooter(),
    )
}
```

## Performance Considerations

### File I/O Optimization

- JSONファイルの読み込みは起動時とデータ変更時のみ
- バッチ操作での複数変更の最適化
- ファイルサイズが大きくなった場合の対策

### Memory Management

- 大量のTODOアイテムに対するメモリ効率的な処理
- TUIでの仮想スクロール（必要に応じて）

### Concurrency

- ファイルアクセスの排他制御
- TUIでの非同期操作の適切な処理

## Security Considerations

### File System Security

- ファイルパスのサニタイゼーション
- 適切なファイル権限の設定
- ディレクトリトラバーサル攻撃の防止

### Input Validation

- ユーザー入力の検証とサニタイゼーション
- JSONデータの検証
- 環境変数の検証

