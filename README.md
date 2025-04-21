# Todo Kanban Board

A terminal-based Todo application with Kanban-style board using Go and the [Charm Bracelet](https://github.com/charmbracelet) libraries.

## Features

- Beautiful TUI (Terminal User Interface) with smooth animations and keyboard navigation
- User authentication with secure password hashing
- Kanban board with three columns: Todo, In Progress, and Done
- Create, edit, and delete tasks
- Move tasks between columns to track progress
- Persistent storage using BadgerDB
- Keyboard shortcuts for all operations

## Installation

### Prerequisites

- Go 1.24 or higher
### Install binary using go install
```
    go install github.com/ReggieReo/todo-elm@v1.0.0
```

### Building from Source

1. Clone the repository:

   ```
   git clone https://github.com/ReggieReo/todo-elm.git
   cd todo-elm
   ```

2. Build the application:

   ```
   go build
   ```

3. Run the application:
   ```
   ./todo-elm
   ```

## Usage

### Authentication

When you first run the application, you'll be presented with a menu to sign in or sign up:

- **Sign up**: Create a new user account with username and password
- **Sign in**: Log in with existing credentials

### Navigation

The Kanban board has three columns: Todo, In Progress, and Done. Use these keyboard shortcuts:

| Key            | Action                                |
| -------------- | ------------------------------------- |
| `↑` / `k`      | Move up within a column               |
| `↓` / `j`      | Move down within a column             |
| `→` / `l`      | Move focus to the right column        |
| `←` / `h`      | Move focus to the left column         |
| `enter`        | Move task to next column              |
| `n`            | Create a new task                     |
| `e`            | Edit the selected task                |
| `d`            | Delete the selected task              |
| `?`            | Toggle help menu                      |
| `q` / `ctrl+c` | Quit the application                  |
| `esc`          | Go back/exit current view             |
| `b`            | Logout and return to the sign-in menu |

### Task Management

- **Create tasks**: Press `n` to create a new task with a title and description
- **Edit tasks**: Press `e` to edit the selected task
- **Delete tasks**: Press `d` to delete the selected task
- **Move tasks**: Press `enter` to move a task to the next column (Todo → In Progress → Done)

## Architecture

The application is built with:

- [BubbleTea](https://github.com/charmbracelet/bubbletea): A Go framework for building terminal apps
- [Bubbles](https://github.com/charmbracelet/bubbles): UI components for BubbleTea
- [Lipgloss](https://github.com/charmbracelet/lipgloss): Styling for terminal applications
- [Huh](https://github.com/charmbracelet/huh): A form library for BubbleTea
- [BadgerDB](https://github.com/dgraph-io/badger): A fast key-value database for persistence

## Data Storage

The application uses BadgerDB to store:

- User credentials (with securely hashed passwords)
- Task data (organized by user and column)

Data is stored in a `badger` directory where the application is run.


## Credits

Created by [ReggieReo](https://github.com/ReggieReo)
Kanbad-board code is based on [Kancli](https://github.com/charmbracelet/kancli/tree/main)
