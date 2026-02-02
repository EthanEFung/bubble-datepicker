// Package datepicker provides a bubble tea component for viewing and selecting
// a date from a monthly view.
package datepicker

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Message
type InvalidDateNavigationMsg struct {
	From time.Time
	To   time.Time
}

// Focus is a value passed to `model.SetFocus` to indicate what component
// controls should be available.
type Focus int

const (
	// FocusNone is a value passed to `model.SetFocus` to ignore all date altering key msgs
	FocusNone Focus = iota
	// FocusHeaderMonth is a value passed to `model.SetFocus` to accept key msgs that change the month
	FocusHeaderMonth
	// FocusHeaderYear is a value passed to `model.SetFocus` to accept key msgs that change the year
	FocusHeaderYear
	// FocusCalendar is a value passed to `model.SetFocus` to accept key msgs that change the week or date
	FocusCalendar
)

//go:generate stringer -type=Focus

// KeyMap is the key bindings for different actions within the datepicker.
type KeyMap struct {
	Up        key.Binding
	Right     key.Binding
	Down      key.Binding
	Left      key.Binding
	FocusPrev key.Binding
	FocusNext key.Binding
	Quit      key.Binding
}

// DefaultKeyMap returns a KeyMap struct with default values
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:        key.NewBinding(key.WithKeys("up", "k")),
		Right:     key.NewBinding(key.WithKeys("right", "l")),
		Down:      key.NewBinding(key.WithKeys("down", "j")),
		Left:      key.NewBinding(key.WithKeys("left", "h")),
		FocusPrev: key.NewBinding(key.WithKeys("shift+tab")),
		FocusNext: key.NewBinding(key.WithKeys("tab")),
		Quit:      key.NewBinding(key.WithKeys("ctrl+c", "q")),
	}
}

// Styles is a struct of lipgloss styles to apply to various elements of the datepicker
type Styles struct {
	Header lipgloss.Style
	Date   lipgloss.Style

	HeaderText   lipgloss.Style
	Text         lipgloss.Style
	SelectedText lipgloss.Style
	FocusedText  lipgloss.Style
	DisabledText lipgloss.Style
}

// DefaultStyles returns a default `Styles` struct
func DefaultStyles() Styles {
	// TODO: refactor for adaptive colors
	r := lipgloss.DefaultRenderer()
	return Styles{
		Header:       r.NewStyle().Padding(1, 0, 0),
		Date:         r.NewStyle().Padding(0, 1, 1),
		HeaderText:   r.NewStyle().Bold(true),
		Text:         r.NewStyle().Foreground(lipgloss.Color("247")),
		SelectedText: r.NewStyle().Bold(true),
		FocusedText:  r.NewStyle().Foreground(lipgloss.Color("212")).Bold(true),
		DisabledText: r.NewStyle().Foreground(lipgloss.Color("240")).Faint(true),
	}
}

// Model is a struct that contains the state of the datepicker component and satisfies
// the `tea.Model` interface
type Model struct {
	// Time is the `time.Time` struct that represents the selected date month and year
	Time time.Time

	// KeyMap encodes the keybindings recognized by the model
	KeyMap KeyMap

	// Styles represent the Styles struct used to render the datepicker
	Styles Styles

	// Focused indicates the component which the end user is focused on
	Focused Focus

	// Selected indicates whether a date is Selected in the datepicker
	Selected bool

	// StartDate represents the start date of the selected date range
	StartDate time.Time

	// EndDate represents the end date of the selected date range
	EndDate time.Time
}

// New returns the Model of the datepicker
func New(time time.Time) Model {
	return Model{
		Time:   time,
		KeyMap: DefaultKeyMap(),
		Styles: DefaultStyles(),

		Focused:  FocusCalendar,
		Selected: false,
	}
}

func NewWithRange(time time.Time, start, end time.Time) Model {
	return Model{
		Time:   time,
		KeyMap: DefaultKeyMap(),
		Styles: DefaultStyles(),

		Focused:  FocusCalendar,
		Selected: false,

		StartDate: start,
		EndDate:   end,
	}
}

// Init satisfies the `tea.Model` interface. This sends a nil cmd
func (m Model) Init() tea.Cmd {
	return nil
}

// Update changes the state of the datepicker. Update satisfies the `tea.Model` interface
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		var cmd tea.Cmd

		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.KeyMap.Up):
			cmd = m.updateUp()

		case key.Matches(msg, m.KeyMap.Right):
			cmd = m.updateRight()

		case key.Matches(msg, m.KeyMap.Down):
			cmd = m.updateDown()

		case key.Matches(msg, m.KeyMap.Left):
			cmd = m.updateLeft()

		case key.Matches(msg, m.KeyMap.FocusPrev):
			switch m.Focused {
			case FocusHeaderYear:
				m.SetFocus(FocusHeaderMonth)
			case FocusCalendar:
				m.SetFocus(FocusHeaderYear)
			}

		case key.Matches(msg, m.KeyMap.FocusNext):
			switch m.Focused {
			case FocusHeaderMonth:
				m.SetFocus(FocusHeaderYear)
			case FocusHeaderYear:
				m.SetFocus(FocusCalendar)
			}
		}

		return m, cmd
	}
	return m, nil
}

func (m *Model) updateUp() tea.Cmd {
	switch m.Focused {
	case FocusHeaderYear:
		return m.LastYear()
	case FocusHeaderMonth:
		return m.LastMonth()
	case FocusCalendar:
		return m.LastWeek()
	case FocusNone:
		// do nothing
	}
	return nil
}

func (m *Model) updateRight() tea.Cmd {
	switch m.Focused {
	case FocusHeaderYear:
		// do nothing
	case FocusHeaderMonth:
		m.SetFocus(FocusHeaderYear)
	case FocusCalendar:
		return m.Tomorrow()
	case FocusNone:
		// do nothing
	}
	return nil
}

func (m *Model) updateDown() tea.Cmd {
	switch m.Focused {
	case FocusHeaderYear:
		return m.NextYear()
	case FocusHeaderMonth:
		return m.NextMonth()
	case FocusCalendar:
		return m.NextWeek()
	case FocusNone:
		// do nothing
	}
	return nil
}

func (m *Model) updateLeft() tea.Cmd {
	switch m.Focused {
	case FocusHeaderYear:
		m.SetFocus(FocusHeaderMonth)
	case FocusHeaderMonth:
		// do nothing
	case FocusCalendar:
		return m.Yesterday()
	case FocusNone:
		// do nothing
	}
	return nil
}

// View renders a month view as a multiline string in the bubbletea application.
// View satisfies the `tea.Model` interface.
func (m Model) View() string {

	b := strings.Builder{}
	month := m.Time.Month()
	year := m.Time.Year()

	tMonth, tYear := month.String(), strconv.Itoa(year)

	if m.Focused == FocusHeaderMonth {
		tMonth = m.Styles.FocusedText.Render(tMonth)
	} else {
		tMonth = m.Styles.HeaderText.Render(tMonth)
	}

	if m.Focused == FocusHeaderYear {
		tYear = m.Styles.FocusedText.Render(tYear)
	} else {
		tYear = m.Styles.HeaderText.Render(tYear)
	}

	title := m.Styles.Header.Render(fmt.Sprintf("%s %s\n", tMonth, tYear))

	// get all the dates of the current month
	firstDayOfTheMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	lastSundayOfLastMonth := firstDayOfTheMonth.AddDate(0, 0, -1)
	for lastSundayOfLastMonth.Weekday() != time.Sunday {
		lastSundayOfLastMonth = lastSundayOfLastMonth.AddDate(0, 0, -1)
	}

	lastDayOfTheMonth := firstDayOfTheMonth.AddDate(0, 1, -1)

	firstSundayOfNextMonth := lastDayOfTheMonth.AddDate(0, 0, 1)
	for firstSundayOfNextMonth.Weekday() != time.Sunday {
		firstSundayOfNextMonth = firstSundayOfNextMonth.AddDate(0, 0, 1)
	}

	day := lastSundayOfLastMonth
	if firstDayOfTheMonth.Weekday() == time.Sunday {
		day = firstDayOfTheMonth
	}

	weekHeaders := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	for i, h := range weekHeaders {
		weekHeaders[i] = m.Styles.Date.Copy().Inherit(m.Styles.HeaderText).Render(h)
	}

	cal := [][]string{weekHeaders}
	j := 1

	for day.Before(firstSundayOfNextMonth) {
		if j >= len(cal) {
			cal = append(cal, []string{})
		}
		out := "  "
		if day.Month() == month {
			out = fmt.Sprintf("%02d", day.Day())
		}

		style := m.Styles.Date
		textStyle := m.Styles.Text

		if !m.Selected {
			// skip modifications to the date
		} else if day.Day() == m.Time.Day() && day.Month() == m.Time.Month() && m.Focused == FocusCalendar {
			textStyle = m.Styles.FocusedText
		} else if day.Day() == m.Time.Day() && day.Month() == m.Time.Month() {
			textStyle = m.Styles.SelectedText
		}

		// check if date is within range, and cross out text if not.
		if !m.StartDate.IsZero() || !m.EndDate.IsZero() {
			dayDate := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)

			if !m.StartDate.IsZero() {
				startDate := time.Date(m.StartDate.Year(), m.StartDate.Month(), m.StartDate.Day(), 0, 0, 0, 0, time.UTC)
				// StartDate inclusive: disable only if strictly before.
				if dayDate.Before(startDate) {
					textStyle = m.Styles.DisabledText
				}
			}

			if !m.EndDate.IsZero() {
				endDate := time.Date(m.EndDate.Year(), m.EndDate.Month(), m.EndDate.Day(), 0, 0, 0, 0, time.UTC)
				// EndDate inclusive: disable only if strictly after.
				if dayDate.After(endDate) {
					textStyle = m.Styles.DisabledText
				}
			}
		}

		out = style.Copy().Inherit(textStyle.Copy()).Render(out)
		cal[j] = append(cal[j], out)

		if day.Weekday() == time.Saturday {
			j++
		}
		day = day.AddDate(0, 0, 1)
	}

	rows := []string{title}
	for _, row := range cal {
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Center, row...))
	}
	b.WriteString(lipgloss.JoinVertical(lipgloss.Center, rows...))

	return b.String()
}

// SetsFocus focuses one of the datepicker components. This can also be used to blur
// the datepicker by passing the Focus `FocusNone`.
func (m *Model) SetFocus(f Focus) {
	m.Focused = f
}

// Blur sets the datepicker focus to `FocusNone`
func (m *Model) Blur() {
	m.Focused = FocusNone
}

// SetTime sets the model's `Time` struct and is used as reference to the selected date
func (m *Model) SetTime(t time.Time) {
	m.Time = t
}

// LastWeek sets the model's `Time` struct back 7 days
func (m *Model) LastWeek() tea.Cmd {
	return m.setTimeWithinRange(m.Time.AddDate(0, 0, -7))
}

// NextWeek sets the model's `Time` struct forward 7 days
func (m *Model) NextWeek() tea.Cmd {
	return m.setTimeWithinRange(m.Time.AddDate(0, 0, 7))
}

// Yesterday sets the model's `Time` struct back 1 day
func (m *Model) Yesterday() tea.Cmd {
	return m.setTimeWithinRange(m.Time.AddDate(0, 0, -1))
}

// Tomorrow sets the model's `Time` struct forward 1 day
func (m *Model) Tomorrow() tea.Cmd {
	return m.setTimeWithinRange(m.Time.AddDate(0, 0, 1))
}

// LastMonth sets the model's `Time` struct back 1 month
func (m *Model) LastMonth() tea.Cmd {
	return m.setTimeWithinRange(m.Time.AddDate(0, -1, 0))
}

// NextMonth sets the model's `Time` struct forward 1 month
func (m *Model) NextMonth() tea.Cmd {
	return m.setTimeWithinRange(m.Time.AddDate(0, 1, 0))
}

// LastYear sets the model's `Time` struct back 1 year
func (m *Model) LastYear() tea.Cmd {
	return m.setTimeWithinRange(m.Time.AddDate(-1, 0, 0))
}

// NextYear sets the model's `Time` struct forward 1 year
func (m *Model) NextYear() tea.Cmd {
	return m.setTimeWithinRange(m.Time.AddDate(1, 0, 0))
}

// SelectDate changes the model's Selected to true
func (m *Model) SelectDate() {
	m.Selected = true
}

// UnselectDate changes the model's Selected to false
func (m *Model) UnselectDate() {
	m.Selected = false
}

// setTimeWithinRange attempts to set the time to the given value, but will
// prevent the change if it falls outside the configured StartDate / EndDate
// range. If StartDate or EndDate are zero, that side of the range is treated
// as unbounded.
func (m *Model) setTimeWithinRange(next time.Time) tea.Cmd {
	// Normalize comparison to date-only (discard hour/min/sec/nano) to match
	// how the calendar is rendered and range is visually applied.
	nextDate := time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, time.UTC)

	if !m.StartDate.IsZero() {
		startDate := time.Date(m.StartDate.Year(), m.StartDate.Month(), m.StartDate.Day(), 0, 0, 0, 0, time.UTC)
		if nextDate.Before(startDate) {
			// Out of range on the lower bound; ignore change.
			return func() tea.Msg {
				return InvalidDateNavigationMsg{From: m.StartDate, To: m.EndDate}
			}
		}
	}

	if !m.EndDate.IsZero() {
		endDate := time.Date(m.EndDate.Year(), m.EndDate.Month(), m.EndDate.Day(), 0, 0, 0, 0, time.UTC)
		if nextDate.After(endDate) {
			// Out of range on the upper bound; ignore change.
			return func() tea.Msg {
				return InvalidDateNavigationMsg{From: m.StartDate, To: m.EndDate}
			}
		}
	}

	m.Time = next
	return nil
}
