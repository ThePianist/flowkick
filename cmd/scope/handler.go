package scope

import (
	"fmt"

	"github.com/ThePianist/flowkick/store"
)

type ScopeSaver interface {
	SaveScope(store.Scope) (int64, error)
}

func ProcessScopeInput(value string, saver ScopeSaver) (int64, error) {
	e := store.Scope{
		Name: value,
	}

	id, err := saver.SaveScope(e)
	if err != nil {
		return 0, fmt.Errorf("failed to save scope %q: %w", value, err)
	}

	return id, nil
}
