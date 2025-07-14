<div align="right">

![golangci-lint](https://github.com/yanosea/gct/actions/workflows/golangci-lint.yml/badge.svg)

</div>

<div align="center">

# ✅ gct

![Language:Go](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)

</div>

## ℹ️ About

`gct` (Go Clean-Architecture Todo) is a clean architecture sample TODO application built with Go.  
This tool provides both CLI and TUI interfaces for efficient todo management with JSON-based storage.

## 📟 CLI

You can manage your todos via the command line interface with simple commands.

### 💻 Usage

```
Available Subcommands:
  add         Add a new todo
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a todo
  help        Help about any command
  list        List all todos
  toggle      Toggle todo status

Flags:
  -h, --help  help for todo
```

### 💻 Examples

```sh
# add a new todo
gct add "Buy groceries"
# list all todos (default command)
gct
# or
gct list
# toggle todo completion status
gct toggle 1
# delete a todo
gct delete 1
# output in JSON format
gct add "Meeting at 3pm" --format json
gct list --format json
```

### 🔧 Installation

#### 🐭 Using go

```sh
go install github.com/yanosea/gct/app/presentation/cli/gct@latest
```

### ✨ Update

#### 🐭 Using go

Reinstall `gct`!

```sh
go install github.com/yanosea/gct/app/presentation/cli/gct@latest
```

### 🧹 Uninstallation

#### 🐭 Using go

```sh
rm $GOPATH/bin/gct
# maybe you have to execute with sudo
rm -fr $GOPATH/pkg/mod/github.com/yanosea/gct
```

## 🖥️ TUI

`gct` also provides a text user interface for interactive todo management.

### 💻 Usage

```sh
# launch TUI mode
gct-tui
```

### ✨ Features

- Interactive todo management
- Real-time todo list updates
- Keyboard navigation
- Clean, minimal interface built with Bubbletea

### 🔧 Installation

#### 🐭 Using go

```sh
go install github.com/yanosea/gct/app/presentation/tui/gct-tui@latest
```

### ✨ Update

#### 🐭 Using go

Reinstall `gct-tui`!

```sh
go install github.com/yanosea/gct/app/presentation/tui/gct-tui@latest
```

### 🧹 Uninstallation

#### 🐭 Using go

```sh
rm $GOPATH/bin/gct-tui
# maybe you have to execute with sudo
rm -fr $GOPATH/pkg/mod/github.com/yanosea/gct-tui
```

## 🌍 Environment Variables

### 📁 Todo data storage location

Default: `$XDG_DATA_HOME/gct/todos.json` or `$HOME/.local/share/gct/todos.json`

```sh
export GCT_DATA_FILE=/path/to/your/todos.json
```

### 🗑️ Remove data files

If you've set custom environment variables, please replace the default paths accordingly.

```sh
# remove todo data file (default location)
rm $HOME/.local/share/gct/todos.json
# remove the entire gct data directory
rm -rf $HOME/.local/share/gct
```

## 🏗️ Architecture

This project follows Clean Architecture principles:

- **Domain Layer**: Todo models and repository interfaces
- **Application Layer**: Use cases for todo operations
- **Infrastructure Layer**: JSON-based repository implementation
- **Presentation Layer**: CLI and TUI interfaces

## 🖊️ Author

[🏹 yanosea](https://github.com/yanosea)

## 🤝 Contributing

Feel free to point me in the right direction🙏
