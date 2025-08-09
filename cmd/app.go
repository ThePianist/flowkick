package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type AppState int

const (
	TextInputState AppState = iota
	ProjectSearchState
	TypeSelectionState
	IssueSearchState
)

type Data struct {
	Entry   string
	Project string
	Type    string
	Issue   string
}

type AppModel struct {
	state              AppState
	textInputModel     Model
	typeSelectionModel TypeSelectionModel
	issueSearchModel   IssueSearchModel
	projectSearchModel ProjectSearchModel
	data               Data
}

func InitialAppModel() AppModel {
	return AppModel{
		state:          TextInputState,
		textInputModel: InitialModel(),
		data:           Data{},
	}
}

func (m AppModel) Init() tea.Cmd {
	switch m.state {
	case TextInputState:
		return m.textInputModel.Init()
	case ProjectSearchState:
		return m.projectSearchModel.Init()
	case TypeSelectionState:
		return m.typeSelectionModel.Init()
	case IssueSearchState:
		return m.issueSearchModel.Init()
	default:
		return nil
	}
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

func (m *AppModel) handleTextInputEnter() (tea.Model, tea.Cmd) {
	m.data.Entry = m.textInputModel.AddEntry.Value()
	log.Print(m.data.Entry)
	m.state = ProjectSearchState
	m.projectSearchModel = NewProjectSearchModel(m.data.Entry, m.data.Type)
	return *m, m.projectSearchModel.Init()
}

func (m *AppModel) handleProjectSelectionEnter() (tea.Model, tea.Cmd) {
	m.data.Project = m.projectSearchModel.textInput.Value()
	log.Print(m.data.Project)
	m.state = TypeSelectionState
	m.typeSelectionModel = NewTypeSelectionModel()
	return *m, m.typeSelectionModel.Init()
}

func (m *AppModel) handleTypeSelectionEnter() (tea.Model, tea.Cmd) {
	m.data.Type = m.typeSelectionModel.choice
	log.Print(m.data.Type)
	m.state = IssueSearchState
	m.issueSearchModel = NewIssueSearchModel(m.data.Entry, m.data.Type)
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

	switch m.state {
	case TextInputState:
		var cmd tea.Cmd
		m.textInputModel, cmd = updateModel(m.textInputModel, msg, func(model Model, msg tea.Msg) (tea.Model, tea.Cmd) {
			return model.Update(msg)
		})

		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			return m.handleTextInputEnter()
		}

		return m, cmd

	case ProjectSearchState:
		var cmd tea.Cmd
		m.projectSearchModel, cmd = updateModel(m.projectSearchModel, msg, func(model ProjectSearchModel, msg tea.Msg) (tea.Model, tea.Cmd) {
			return model.Update(msg)
		})
		m.data.Project = m.projectSearchModel.textInput.Value()

		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			return m.handleProjectSelectionEnter()
		}

		return m, cmd

	case TypeSelectionState:
		var cmd tea.Cmd
		m.typeSelectionModel, cmd = updateModel(m.typeSelectionModel, msg, func(model TypeSelectionModel, msg tea.Msg) (tea.Model, tea.Cmd) {
			return model.Update(msg)
		})

		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			return m.handleTypeSelectionEnter()
		}

		return m, cmd

	case IssueSearchState:
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
	switch m.state {
	case TextInputState:
		return m.textInputModel.View()
	case ProjectSearchState:
		return m.projectSearchModel.View()
	case TypeSelectionState:
		return m.typeSelectionModel.View()
	case IssueSearchState:
		return m.issueSearchModel.View()
	default:
		return "Unknown state"
	}
}

func (m AppModel) saveAndExit() (tea.Model, tea.Cmd) {
	log.Printf("Saving data and exiting: %+v", m.data)
	// TODO: persist m.data to SQLite here
	log.Println("Success! Data logged.")
	clearTerminal()
	return m, tea.Quit
}

func clearTerminal() {
	print("\033[H\033[2J")
}
