package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type snippetModel struct {
	DB *sql.DB
}

type SnippetModel interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (Snippet, error)
	Latest() ([]Snippet, error)
}

func NewSnippetModel(db *sql.DB) SnippetModel {
	return &snippetModel{DB: db}
}

func (m *snippetModel) Insert(title string, content string, expires int) (int, error) {
	// TODO: there is a bug here, 4 columns but 5 values
	stmt := `INSERT INTO snippets (title, content, created, expires) 
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	res, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *snippetModel) Get(id int) (Snippet, error) {
	query := `SELECT id, title, content, created, expires 
	FROM snippets WHERE id = ? AND expires > UTC_TIMESTAMP();`

	var s Snippet
	err := m.DB.QueryRow(query, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		}
		return Snippet{}, err
	}

	return s, nil
}

func (m *snippetModel) Latest() ([]Snippet, error) {
	query := `SELECT id, title, content, created, expires 
	FROM snippets WHERE expires > UTC_TIMESTAMP()
	ORDER BY created DESC
	LIMIT 10;`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet
	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
