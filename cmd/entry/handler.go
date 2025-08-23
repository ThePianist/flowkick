package entry

import (
	"log"
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
		return 0, err
	}
	log.Print(value)
	return e.ID, nil
}
