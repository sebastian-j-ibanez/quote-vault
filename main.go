package main

import (
	"fmt"
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
		return "You've selected: " + m.event
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
			fmt.Println("Exiting...")
			return m, tea.Quit
		case tea.KeyRunes:
			switch msg.String() {
			case "down", "j":
				fmt.Println("J has been pressed...")
			}

		}

	}

	return m, nil
}

var _ tea.Model = (*model)(nil)

func main() {

	// Init Db
	db, err := InitDb()
	if err != nil {
		// Potentially pass a message to the update function to display error.
		os.Exit(1)
	}

	//quotes, err := GetAllQuotes(db)
	quote, err := GetQuote(db, 1)

	// Check for errors from GetQuote.
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(quote.ID, quote.Body, quote.Author, quote.Date)

	//state, err := NewModel()
	//if err != nil {
	//	fmt.Println(fmt.Printf("Error starting init command: %s\n", err))
	//	os.Exit(1)
	//}

	//if _, err := tea.NewProgram(state, tea.WithAltScreen()).Run(); err != nil {
	//	log.Fatal(err)
	//}
}
