package cmd

import (
	"github.com/ThePianist/flowkick/cmd/entry"
	"github.com/ThePianist/flowkick/cmd/scope"
	"github.com/ThePianist/flowkick/store"
)

type AppState int

type Data struct {
	Entry store.Entry
	Scope store.Scope
	Type  string
	Issue string
}

type AppModel struct {
	currentView        AppState
	entryModel         entry.EntryModel
	typeSelectionModel TypeSelectionModel
	issueSearchModel   IssueSearchModel
	scopeModel         scope.ScopeModel
	data               Data
	store              *store.Store
	dataEntryID        int64 // Store entry ID for linking
}
