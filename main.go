package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	persistence "github.com/ReggieReo/todo-elm/persistance"
	"github.com/ReggieReo/todo-elm/todolist"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type uiState int

const (
	menu uiState = iota
	signIn
	signUp
	submitting
	authenticated
)

type authSuccessMsg struct{ username string }
type authErrMsg struct{ err error }

type model struct {
	width, height int
	state         uiState
	form          *huh.Form
	store         *persistence.Store
	spinner       spinner.Model
	board         tea.Model
	err           error
	username      string
	opInProgress  string
}

const titleStr = `
 _________  ________  ________  ________  ___  ___  ___
|\___   ___\\   __  \|\   ___ \|\   __  \|\  \|\  \|\  \
\|___ \  \_\ \  \|\  \ \  \_|\ \ \  \|\  \ \  \ \  \ \  \
     \ \  \ \ \  \\\  \ \  \ \\ \ \  \\\  \ \  \ \  \ \  \
      \ \  \ \ \  \\\  \ \  \_\\ \ \  \\\  \ \__\ \__\ \__\
       \ \__\ \ \_______\ \_______\ \_______\|__|\|__|\|__|
        \|__|  \|_______|\|_______|\|_______|   ___  ___  ___
                                               |\__\|\__\|\__\
                                               \|__|\|__|\|__|
`

// command for persistance
// createUserCmd creates a tea.Cmd that attempts to save the user via the persistence store.
func createUserCmd(store *persistence.Store, username, password string) tea.Cmd {
	return func() tea.Msg {
		err := store.CreateUser(username, password)
		if err != nil {
			return authErrMsg{err}
		}
		return authSuccessMsg{username: username}
	}
}

// authenticateUserCmd creates a tea.Cmd that attempts to sign in the user via the persistence store.
func authenticateUserCmd(store *persistence.Store, username, password string) tea.Cmd {
	return func() tea.Msg {
		uname, err := store.AuthenticateUser(username, password)
		if err != nil {
			return authErrMsg{err} // e.g., "invalid username or password"
		}
		return authSuccessMsg{username: uname}
	}
}

func createSignInForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("username").
				Title("Username").
				Placeholder("Enter your username").
				Validate(func(s string) error {
					if s == "" {
						return errors.New("username cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Key("password").
				Title("Password").
				Placeholder("Enter your password").
				EchoMode(huh.EchoModePassword).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("password cannot be empty")
					}
					return nil
				}),
		),
	)
}

func createSignUpForm() *huh.Form {
	var password string // Variable to store the first password for validation
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("username").
				Title("Choose Username").
				Placeholder("Enter a username").
				Validate(func(s string) error {
					if len(s) < 3 {
						return errors.New("username must be at least 3 characters")
					}
					if strings.ContainsAny(s, " ") {
						return errors.New("username cannot contain spaces")
					}
					return nil
				}),
			huh.NewInput().
				Key("password").
				Title("Choose Password").
				Placeholder("Enter a password").
				EchoMode(huh.EchoModePassword).
				Validate(func(s string) error {
					if len(s) < 6 {
						return errors.New("Password must be at least 6 characters.")
					}
					password = s // Store the password for the confirmation check
					return nil
				}),
			huh.NewInput().
				Key("confirmPassword").
				Title("Confirm Password").
				Placeholder("Enter the password again").
				EchoMode(huh.EchoModePassword).
				Validate(func(s string) error {
					if s != password {
						return errors.New("passwords do not match")
					}
					return nil
				}),
		),
	)
}
func createMenuForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("signOption").
				Title("Welcome to TODO!!!").
				Options(huh.NewOptions("Sign-in", "Sign-up")...),
		),
	)
	return f
}

func initialModel(store *persistence.Store) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		form:    createMenuForm(),
		state:   menu,
		store:   store,
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), m.spinner.Tick, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// handle authenticated status
	if m.state == authenticated {
		b, cmd := m.board.Update(msg)
		m.board = b
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "b":
			if m.state != menu {
				m.state = menu
				m.form = createMenuForm()
				m.err = nil
				m.opInProgress = ""
				cmds = append(cmds, m.form.Init())
				return m, tea.Batch(cmds...)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case authSuccessMsg:
		m.username = msg.username
		m.form = nil
		m.err = nil
		m.state = authenticated
		m.opInProgress = ""
		m.board = todolist.NewBoard()
		fakeResize := func() tea.Msg {
			return tea.WindowSizeMsg{Width: m.width, Height: m.height}
		}
		return m, fakeResize

	case authErrMsg:
		m.err = msg.err
		switch m.opInProgress {
		case "signin":
			m.state = signIn
			m.form = createSignInForm() // Re-create the sign-in form
		case "signup":
			m.state = signUp
			m.form = createSignUpForm()
		default:
			m.state = menu
			m.form = createMenuForm()
		}
		m.opInProgress = ""
		if m.form != nil {
			cmds = append(cmds, m.form.Init())
		}
		return m, tea.Batch(cmds...)

	case spinner.TickMsg:
		var spinCmd tea.Cmd
		// Only tick spinner if we are in the submitting state
		if m.state == submitting {
			m.spinner, spinCmd = m.spinner.Update(msg)
			cmds = append(cmds, spinCmd)
		}
		return m, tea.Batch(cmds...)
	}

	// form handling
	if m.form != nil && (m.state == menu || m.state == signIn || m.state == signUp) {
		var newForm tea.Model // huh.Form implements tea.Model
		newForm, formCmd := m.form.Update(msg)
		if f, ok := newForm.(*huh.Form); ok {
			m.form = f
			cmds = append(cmds, formCmd)

			if m.form.State == huh.StateCompleted {
				m.err = nil // Clear previous errors when attempting submission

				switch m.state {
				case menu:
					signOption := m.form.GetString("signOption")
					switch signOption {
					case "Sign-in":
						m.state = signIn
						m.form = createSignInForm()
						cmds = append(cmds, m.form.Init()) // Initialize the new form
					case "Sign-up":
						m.state = signUp
						m.form = createSignUpForm()
						cmds = append(cmds, m.form.Init()) // Initialize the new form
					}
				case signIn:
					username := m.form.GetString("username")
					password := m.form.GetString("password")
					m.state = submitting
					m.opInProgress = "signin"
					m.form = nil
					cmds = append(cmds, m.spinner.Tick, authenticateUserCmd(m.store, username, password))
				case signUp:
					username := m.form.GetString("username")
					password := m.form.GetString("password")
					m.state = submitting
					m.opInProgress = "signup"
					m.form = nil
					cmds = append(cmds, m.spinner.Tick, createUserCmd(m.store, username, password))
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var viewContent string
	footer := "\nPress 'q' to quit."

	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red

	if m.state != menu && m.state != submitting {
		footer += " Press 'b' to go back."
	}

	// Prepare error string if an error exists
	errorStr := ""
	if m.err != nil {
		errorStr = errorStyle.Render("Error: " + m.err.Error())
	}

	// Determine view content based on state
	switch m.state {
	case menu, signIn, signUp:
		formView := ""
		if m.form != nil {
			formView = m.form.View()
		}
		body := lipgloss.JoinVertical(lipgloss.Center,
			titleStr,
			errorStr,
			formView,
			footer,
		)
		viewContent = body

	case submitting:
		body := lipgloss.JoinVertical(lipgloss.Center,
			titleStr,
			fmt.Sprintf("%s Submitting...", m.spinner.View()), // Show spinner
			footer,
		)
		viewContent = body

	case authenticated:
		msg := fmt.Sprintf("Welcome, %s!", m.username)
		bContent := m.board.View()
		// return bContent
		viewContent = lipgloss.JoinVertical(lipgloss.Center, msg, footer, bContent)

	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center, // Horizontal alignment
		lipgloss.Center, // Vertical alignment
		viewContent,     // The content string to place
	)
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Warning: Could not get user home directory: %v. Using temp dir.", err)
		homeDir = os.TempDir()
	}
	dbBaseDir := filepath.Join(homeDir, ".todo-elm")

	store, err := persistence.NewStore(dbBaseDir)
	if err != nil {
		log.Fatalf("Failed to initialize persistence store: %v", err)
	}
	// Ensure the database is closed when the program exits
	defer func() {
		if err := store.Close(); err != nil {
			log.Printf("Error closing persistence store: %v", err)
		}
	}()

	p := tea.NewProgram(initialModel(store), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}
