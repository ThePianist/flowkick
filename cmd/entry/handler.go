package entry

import (
	"fmt"
	"time"

	"github.com/ThePianist/flowkick/store"
)

type EntrySaver interface {
	SaveEntry(store.Entry) error
}

func ProcessEntryInput(value string, saver EntrySaver) (int64, error) {
	e := store.Entry{
		ID:        time.Now().UnixNano(),
		Note:      value,
		CreatedAt: time.Now(),
	}
	if err := saver.SaveEntry(e); err != nil {
		return 0, fmt.Errorf("failed to save entry %q: %w", value, err)
	}
	return e.ID, nil
}
