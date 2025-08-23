package store

type Ticket struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (s *Store) SaveTicket(ticketID int64, name string, entryID int64) error {
	upsertTicket := `INSERT INTO tickets (id, name)
	VALUES (?, ?)
	ON CONFLICT(id) DO UPDATE SET name=excluded.name;`
	_, err := s.Conn.Exec(upsertTicket, ticketID, name)
	if err != nil {
		return err
	}
	updateEntry := `UPDATE entries SET ticket_id=? WHERE id=?;`
	_, err = s.Conn.Exec(updateEntry, ticketID, entryID)
	return err
}

func (s *Store) GetTickets() ([]string, error) {
	rows, err := s.Conn.Query("SELECT name FROM tickets ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tickets []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tickets = append(tickets, name)
	}
	return tickets, nil
}
