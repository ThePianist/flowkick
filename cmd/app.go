package cmd

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/ThePianist/flowkick/cmd/entry"
	"github.com/ThePianist/flowkick/cmd/scope"
	"github.com/ThePianist/flowkick/store"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	EntryInputView AppState = iota
	ScopeInputView
	TypeSelectView
	TicketInputView
)

func InitialAppModel(store *store.Store) AppModel {
	return AppModel{
		currentView: EntryInputView,
		entryModel:  entry.NewEntryModel(),
		data:        Data{},
		store:       store,
	}
}

func (m AppModel) Init() tea.Cmd {
	switch m.currentView {
	case EntryInputView:
		return m.entryModel.Init()
	default:
		return nil
	}
}

func (m *AppModel) handleEntrySubmit() (tea.Model, tea.Cmd) {
	// read input
	inputText := m.entryModel.GetValue()

	// persist entry first (domain logic)
	id, err := entry.ProcessEntryInput(inputText, m.store)
	if err != nil {
		log.Printf("save failed: %v", err)
		// stay in the same state so user can retry
		return *m, nil
	}

	m.data.Entry = store.Entry{
		ID:        id,
		Note:      inputText,
		CreatedAt: time.Now(),
	}

	// now prepare the next UI state
	m.scopeModel = scope.NewScopeModel(inputText, m.data.Type, m.store)
	m.currentView = ScopeInputView

	return *m, m.scopeModel.Init()
}

func (m *AppModel) handleScopeSubmit() (tea.Model, tea.Cmd) {
	// read input
	inputText := m.scopeModel.GetValue()

	// persist scope (domain logic)
	scopeId, err := scope.ProcessScopeInput(inputText, m.store)
	if err != nil {
		log.Printf("save failed: %v", err)
		// stay in the same state so user can retry
		return *m, nil
	}

	if m.data.Entry.ID != 0 {
		m.store.SaveEntry(store.Entry{
			ID:        m.data.Entry.ID,
			Note:      m.data.Entry.Note,
			ScopeID:   sql.NullInt64{Int64: scopeId, Valid: true},
			CreatedAt: time.Now(),
		})
	}

	// now prepare the next UI state
	m.currentView = TypeSelectView
	m.typeSelectionModel = NewTypeSelectionModel(m.store)

	return *m, m.typeSelectionModel.Init()
}

func (m *AppModel) handleTypeSelectionEnter() (tea.Model, tea.Cmd) {
	m.data.Type = m.typeSelectionModel.choice
	log.Print(m.data.Type)
	if m.dataEntryID != 0 {

		var typeID int64
		if parsedTypeID, err := strconv.ParseInt(m.data.Type, 10, 64); err == nil {
			typeID = parsedTypeID
		} else {
			typeID = 0 // or handle error as needed
			log.Printf("Failed to convert Type to int64: %v", err)
		}

		m.store.SaveEntry(store.Entry{
			ID:        m.data.Entry.ID,
			Note:      m.data.Entry.Note,
			TypeID:    sql.NullInt64{Int64: typeID, Valid: typeID != 0},
			CreatedAt: time.Now(),
		})
	}
	m.currentView = TicketInputView
	m.issueSearchModel = NewIssueSearchModel(m.data.Entry.Note, m.data.Type, m.store)
	return *m, m.issueSearchModel.Init()
}

func updateModel[M any](model M, msg tea.Msg, updateFunc func(M, tea.Msg) (tea.Model, tea.Cmd)) (M, tea.Cmd) {
	updated, cmd := updateFunc(model, msg)
	return updated.(M), cmd
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if cmd, quit := m.handleQuit(msg); quit {
		return m, cmd
	}

	switch m.currentView {
	case EntryInputView:
		var cmd tea.Cmd
		m.entryModel, cmd = updateModel(m.entryModel, msg, func(model entry.EntryModel, msg tea.Msg) (tea.Model, tea.Cmd) {
			return model.Update(msg)
		})

		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			return m.handleEntrySubmit()
		}

		return m, cmd

	case ScopeInputView:
		var cmd tea.Cmd
		m.scopeModel, cmd = updateModel(m.scopeModel, msg, func(model scope.ScopeModel, msg tea.Msg) (tea.Model, tea.Cmd) {
			return model.Update(msg)
		})
		m.data.Scope = store.Scope{
			Name: m.scopeModel.GetValue(),
		}

		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			return m.handleScopeSubmit()
		}

		return m, cmd

	case TypeSelectView:
		var cmd tea.Cmd
		m.typeSelectionModel, cmd = updateModel(m.typeSelectionModel, msg, func(model TypeSelectionModel, msg tea.Msg) (tea.Model, tea.Cmd) {
			return model.Update(msg)
		})

		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			return m.handleTypeSelectionEnter()
		}

		return m, cmd

	case TicketInputView:
		var cmd tea.Cmd
		m.issueSearchModel, cmd = updateModel(m.issueSearchModel, msg, func(model IssueSearchModel, msg tea.Msg) (tea.Model, tea.Cmd) {
			return model.Update(msg)
		})
		m.data.Issue = m.issueSearchModel.textInput.Value()

		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			return m.saveAndExit()
		}

		return m, cmd

	default:
		return m, nil
	}
}

func (m AppModel) View() string {
	switch m.currentView {
	case EntryInputView:
		return m.entryModel.View()
	case ScopeInputView:
		return m.scopeModel.View()
	case TypeSelectView:
		return m.typeSelectionModel.View()
	case TicketInputView:
		return m.issueSearchModel.View()
	default:
		return "Unknown state"
	}
}

func (m AppModel) saveAndExit() (tea.Model, tea.Cmd) {
	log.Printf("Saving data and exiting: %+v", m.data)
	if m.dataEntryID != 0 && m.data.Issue != "" {
		m.store.SaveTicket(0, m.data.Issue, m.data.Entry.ID)
		m.store.SaveEntry(store.Entry{
			ID:        m.data.Entry.ID,
			Note:      m.data.Entry.Note,
			TicketID:  sql.NullInt64{Int64: 0, Valid: false}, // You may want to fetch the actual ticket ID
			CreatedAt: time.Now(),
		})
	}
	log.Println("Success! Data logged.")
	clearTerminal()
	return m, tea.Quit
}

func (m AppModel) handleQuit(msg tea.Msg) (tea.Cmd, bool) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.Type == tea.KeyCtrlC || (key.Type == tea.KeyRunes && key.String() == "q") {
			// Exit without saving mid-flow
			log.Println("Exiting before completion, data not saved.")
			return tea.Quit, true
		}
	}
	return nil, false
}

func clearTerminal() {
	print("\033[H\033[2J")
}
