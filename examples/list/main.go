package main

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanefung/bubble-datepicker"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type DayItem struct {
	label string
	time.Time
}

func (d DayItem) String() string {
	return d.label
}

func (d DayItem) Title() string {
	return d.label
}

func (d DayItem) Description() string {
	return d.Format(time.DateOnly)
}

func (d DayItem) FilterValue() string {
	return d.Format(time.DateOnly)
}

type model struct {
	holidays   list.Model
	datepicker datepicker.Model
}

func initializeModel() tea.Model {
	dates := []list.Item{
		DayItem{"Halloween", time.Date(2023, time.October, 31, 0, 0, 0, 0, time.UTC)},
		DayItem{"Thanksgiving", time.Date(2023, time.November, 23, 0, 0, 0, 0, time.UTC)},
		DayItem{"Christmas", time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC)},
		DayItem{"New Years", time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)},
	}

	l := list.New(dates, list.NewDefaultDelegate(), 0, 0)
	dp := datepicker.New(time.Now())

	item := l.SelectedItem().(DayItem) // sad
	dp.SetTime(item.Time)

	return model{
		holidays:   l,
		datepicker: dp,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		x, y := docStyle.GetFrameSize()
		m.holidays.SetSize(msg.Width-x, msg.Height-y)
	}

	var cmd tea.Cmd
	m.holidays, cmd = m.holidays.Update(msg)

	item := m.holidays.SelectedItem().(DayItem) // sad
	m.datepicker.SetTime(item.Time)

	return m, cmd
}

func (m model) View() string {
	content := lipgloss.JoinHorizontal(lipgloss.Left, m.holidays.View(), m.datepicker.View())
	return docStyle.Render(content)
}

func main() {
	p := tea.NewProgram(initializeModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, an error %s:", err)
	}
}
