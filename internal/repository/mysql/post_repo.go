package mysql

import (
	"Test2/internal/domain"
	"context"
	"database/sql"
	"fmt"
)

func NewMysqlPostRepository(db *sql.DB) domain.PostRepository {
	return &mysqlPostRepo{db}
}

type mysqlPostRepo struct {
	db *sql.DB
}

func (m *mysqlPostRepo) Fetch(ctx context.Context, limit int64, offset int64) ([]domain.Post, error) {
	query := `SELECT id, title, description, content, thumbnail, status, update_date, created_at
			  FROM posts
			  WHERE status != ?
			  ORDER BY created_at DESC
			  LIMIT ? OFFSET ?`

	rows, err := m.db.QueryContext(ctx, query, domain.StatusDeleted, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]domain.Post, 0, int(limit))

	for rows.Next() {
		p := domain.Post{}
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Content, &p.Thumbnail, &p.Status, &p.UpdateDate, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (m *mysqlPostRepo) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	query := `SELECT id, title, description, content, thumbnail, status, update_date, created_at
				FROM posts
				WHERE id = ?
				AND status != ?`

	row := m.db.QueryRowContext(ctx, query, id, domain.StatusDeleted)

	p := &domain.Post{}
	err := row.Scan(&p.ID, &p.Title, &p.Description, &p.Content, &p.Thumbnail, &p.Status, &p.UpdateDate, &p.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, err
	}
	return p, nil
}

func (m *mysqlPostRepo) Store(ctx context.Context, p *domain.Post) error {
	query := `INSERT INTO posts (title, description, content, thumbnail, status, update_date, created_at)
				VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := m.db.ExecContext(ctx, query, p.Title, p.Description, p.Content, p.Thumbnail, p.Status, p.UpdateDate, p.CreatedAt)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	p.ID = id

	return nil
}

func (m *mysqlPostRepo) Update(ctx context.Context, p *domain.Post) error {
	query := `UPDATE posts SET 
				title = ?, 
				description = ?,
				content = ?,
				thumbnail = ?,
				status = ?, 
				update_date = ? 
				WHERE id = ?`

	_, err := m.db.ExecContext(ctx, query, p.Title, p.Description, p.Content, p.Thumbnail, p.Status, p.UpdateDate, p.ID)

	return err
}

func (m *mysqlPostRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE posts SET
				status = ?
				WHERE id = ?`

	_, err := m.db.ExecContext(ctx, query, domain.StatusDeleted, id)

	return err
}

func (m *mysqlPostRepo) Search(ctx context.Context, keyword string, limit int64, offset int64) ([]domain.Post, error) {
	query := `SELECT id, title, description, content, thumbnail, status, update_date, created_at
			  FROM posts
			  WHERE status != ?
			  AND (
			  	title LIKE ?
			  	OR description LIKE ?
			  	OR content LIKE ?
			  )
			  ORDER BY created_at DESC
			  LIMIT ? OFFSET ?`

	rows, err := m.db.QueryContext(ctx, query, domain.StatusDeleted, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]domain.Post, 0, int(limit))

	for rows.Next() {
		p := domain.Post{}
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Content, &p.Thumbnail, &p.Status, &p.UpdateDate, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}
