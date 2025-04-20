package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type uiState int

const (
	menu uiState = iota
	signIn
	signUp
)

type model struct {
	width, height int
	state         uiState
	form          *huh.Form
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

func initForm() *huh.Form {
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

func initialModel() model {
	return model{form: initForm(), state: menu}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "b":
			if m.state != menu {
				m.state = menu
				m.form = initForm()
				cmds = append(cmds, m.form.Init())
				return m, tea.Batch(cmds...)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.state == menu {
		// Let the form handle its own keys and events.
		newForm, cmd := m.form.Update(msg)
		if f, ok := newForm.(*huh.Form); ok {
			m.form = f
		}
		cmds = append(cmds, cmd)

		// Check if the form is completed.
		if m.form.State == huh.StateCompleted {
			// Get the selected value.
			signOption := m.form.GetString("signOption")
			switch signOption {
			case "Sign-in":
				m.state = signIn
			case "Sign-up":
				m.state = signUp
			default:
				// Should not happen with the current options, but good practice
				m.state = menu
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// if form is done, just show the choice
	footer := "\nPress q to quit.\n"
	if m.state == signIn {
		return lipgloss.JoinVertical(lipgloss.Center, "You selete sign In", footer+"Press b to go back to main menu.")
	}

	if m.state == signUp {
		return lipgloss.JoinVertical(lipgloss.Center, "You select sign Up", footer+"Press b to go back to main menu.")
	}

	// assemble the pieces
	body := lipgloss.JoinVertical(lipgloss.Center,
		titleStr,
		m.form.View(),
		footer,
	)

	// center it in the available window
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		body,
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
