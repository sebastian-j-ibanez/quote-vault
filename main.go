package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func NewModel() (*model, error) {
	return &model{}, nil
}

type model struct {
	nameInput string
	listInput string
	event     string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	if m.event != "" {
		return fmt.Sprintf("You've selected: %s", m.event)
	}
	return "Quote Vault"
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		_, _ = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlBackslash:
			return m, tea.Quit
		}
	}

	return m, nil
}

var _ tea.Model = (*model)(nil)

func main() {
	state, err := NewModel()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error starting init command: %s\n", err))
		os.Exit(1)
	}

	if _, err := tea.NewProgram(state).Run(); err != nil {
		log.Fatal(err)
	}
}
