package entry

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type (
	errMsg error
)

type EntryModel struct {
	Input textinput.Model
	Err   error
}

func NewEntryModel() EntryModel {
	ti := textinput.New()
	ti.Placeholder = "Fixed weird cache bug after 2hrs"
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 156
	return EntryModel{Input: ti, Err: nil}
}

func (m EntryModel) GetValue() string {
	return m.Input.Value()
}
