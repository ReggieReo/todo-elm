package todolist

import (
	"log"

	persistence "github.com/ReggieReo/todo-elm/persistance"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Board struct {
	help     help.Model
	loaded   bool
	focused  status
	cols     []column
	quitting bool
	username string
	store    *persistence.Store
}

var board *Board

func NewBoard(username string, store *persistence.Store) *Board {
	help := help.New()
	help.ShowAll = true
	board = &Board{
		help:     help,
		focused:  todo,
		username: username,
		store:    store,
	}
	board.initLists()
	return board
}

func (m *Board) Init() tea.Cmd {
	return nil
}

func (m *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		m.help.Width = msg.Width - margin
		for i := 0; i < len(m.cols); i++ {
			var res tea.Model
			res, cmd = m.cols[i].Update(msg)
			m.cols[i] = res.(column)
			cmds = append(cmds, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmds...)
	case *Form:
		cmd := m.cols[m.focused].Set(msg.index, msg.CreateTask())
		if err := m.saveTasks(); err != nil {
			log.Printf("Error saving tasks: %v", err)
		}
		return m, cmd
	case moveMsg:
		cmd := m.cols[m.focused.getNext()].Set(APPEND, msg.Task)
		if err := m.saveTasks(); err != nil {
			log.Printf("Error saving tasks: %v", err)
		}
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			if err := m.saveTasks(); err != nil {
				log.Printf("Error saving tasks: %v", err)
			}
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.Left):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getPrev()
			m.cols[m.focused].Focus()
		case key.Matches(msg, keys.Right):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getNext()
			m.cols[m.focused].Focus()
		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
	}
	res, cmd := m.cols[m.focused].Update(msg)
	if _, ok := res.(column); ok {
		m.cols[m.focused] = res.(column)
		if deleteMsg, ok := msg.(deleteMsg); ok && deleteMsg.status == m.focused {
			if err := m.saveTasks(); err != nil {
				log.Printf("Error saving tasks after deletion: %v", err)
			}
		}
	} else {
		return res, cmd
	}
	return m, cmd
}

// Changing to pointer receiver to get back to this model after adding a new task via the form... Otherwise I would need to pass this model along to the form and it becomes highly coupled to the other models.
func (m *Board) View() string {
	if m.quitting {
		return ""
	}
	if !m.loaded {
		return "loading..."
	}
	boardView := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.cols[todo].View(),
		m.cols[inProgress].View(),
		m.cols[done].View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left, boardView, m.help.View(keys))
}
