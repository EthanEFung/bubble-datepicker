/*
An interesting problem to solve in this example is acheiving
two way data binding. On one hand there the input should
update the view of the datepicker, and in another context
the datepicker should change the date that is displayed in
the text input. This of course means that the datepicker itself
needs to be able to support a state where it is displaying a
date, but indicates to the user that the date reflected is not
selected.
*/
package main

import (
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	datepicker "github.com/ethanefung/bubble-datepicker"
)

type focus int

const (
    FocusNone focus = iota
    FocusInput
    FocusDatePicker
)

type Model struct {
    focus focus
    input textinput.Model
    datepicker datepicker.Model
}

var inputStyles = lipgloss.NewStyle().Padding(1, 1, 0)

func initializeModel() tea.Model {
    dp := datepicker.New(time.Now())

    input := textinput.New()
    input.Placeholder = "YYYY-MM-DD (enter date)"
    input.Focus()
    input.Width = 20

    return Model{
        focus: FocusInput,
        input: input,
        datepicker: dp,
    }
}

func (m Model) Init() tea.Cmd {
    return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // TODO figure out how we want to size things
        // we'll probably want both bubbles to be vertically stacked
        // and to take as much room as the can
        return m, nil
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "tab":
            if m.focus == FocusInput {
                m.focus = FocusDatePicker
                m.input.Blur()
                m.input.SetValue(m.datepicker.Time.Format(time.DateOnly))

                m.datepicker.SelectDate()
                m.datepicker.SetFocus(datepicker.FocusHeaderMonth)
                m.datepicker = m.datepicker
                return m, nil

            }
        case "shift+tab":
            if m.focus == FocusDatePicker && m.datepicker.Focused == datepicker.FocusHeaderMonth {
                m.focus = FocusInput
                m.datepicker.Blur()

                m.input.Focus()
                return m, nil
            }
        }
    }

    switch m.focus {
    case FocusInput:
        m.input, cmd = m.UpdateInput(msg)
    case FocusDatePicker:
        m.datepicker, cmd = m.UpdateDatepicker(msg)
    case FocusNone:
        // do nothing
    }

    return  m, cmd 
}


func (m Model) View() string {
    return lipgloss.JoinVertical(lipgloss.Left, inputStyles.Render(m.input.View()), m.datepicker.View())
}

func (m *Model) UpdateInput(msg tea.Msg) (textinput.Model, tea.Cmd) {
    var cmd tea.Cmd

    m.input, cmd = m.input.Update(msg)

    val := m.input.Value()
    t, err := time.Parse(time.DateOnly, strings.TrimSpace(val))
    if err == nil {
        m.datepicker.SetTime(t)
        m.datepicker.SelectDate()
        m.datepicker.Blur()
    } 
    if err != nil && m.datepicker.Selected {
        m.datepicker.UnselectDate()
    }

    return m.input, cmd
}

func (m *Model) UpdateDatepicker(msg tea.Msg) (datepicker.Model, tea.Cmd) {
    var cmd tea.Cmd

    prev := m.datepicker.Time

    m.datepicker, cmd = m.datepicker.Update(msg)

    if prev != m.datepicker.Time {
        m.input.SetValue(m.datepicker.Time.Format(time.DateOnly))
    }

    return m.datepicker, cmd
}

func main() {
    p := tea.NewProgram(initializeModel())
    if _, err := p.Run(); err != nil {
        log.Fatalf("alas, an error %s", err)
    }
}
