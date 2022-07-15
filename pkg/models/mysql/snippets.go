package mysql

import (
	"database/sql"

	"jeisaRaja.git/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id,title,content,created,expires FROM snippets where expires > UTC_TIMESTAMP() AND id=?`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT * FROM snippets WHERE UTC_TIMESTAMP()>expires ORDER BY created DESC LIMIT 10`
	snippets := []*models.Snippet{}
	rows, err := m.DB.Query(stmt)
	if err == sql.ErrNoRows || err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		s := &models.Snippet{}
		err := rows.Scan(s.ID, s.Title, s.Content, s.Created, s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err == rows.Err() || err != nil {
		return nil, err
	}

	return snippets, nil
}
