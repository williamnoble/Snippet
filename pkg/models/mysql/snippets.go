package mysql

import (
	"database/sql"
	"errors"
	"snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

// Insert a single snippet to the database
func (m *SnippetModel) Insert(title string, content string, expires string) (int, error) {
	stmt := `INSERT INTO snippets(title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

//Get a single Snippet by ID from the database
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// Select where expiry date is greater (i.e. later) than today.
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	snippet := &models.Snippet{}

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

//Latest retrieves the Latest snippets form a database (limit to the most recent 10 rows)
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP ORDER BY created LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var snippets []*models.Snippet

	for rows.Next() {
		s := &models.Snippet{}
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
