package mysql

import (
	"Test2/internal/domain"
	"context"
	"database/sql"
	"fmt"
)

func NewMysqlCateRepository(db *sql.DB) domain.CategoryRepository {
	return &mysqlCateRepo{db}
}

type mysqlCateRepo struct {
	db *sql.DB
}

func (m *mysqlCateRepo) Fetch(ctx context.Context, limit int64, offset int64) ([]domain.Category, error) {
	query := `SELECT id, title, description, thumbnail, status, updated_at, created_at
				FROM categories
				WHERE status != ?
				ORDER BY created_at DESC
				LIMIT ? OFFSET ?`
	rows, err := m.db.QueryContext(ctx, query, domain.CategoryStatusInactive, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]domain.Category, 0, int(limit))

	for rows.Next() {
		c := domain.Category{}
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Thumbnail, &c.Status, &c.UpdatedAt, &c.CreatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, c)
	}
	return result, nil
}

func (m *mysqlCateRepo) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	query := `SELECT id, title, description, thumbnail, status, updated_at, created_at
				FROM categories
				WHERE id = ?
				AND status != ?`

	row := m.db.QueryRowContext(ctx, query, id, domain.CategoryStatusInactive)

	c := &domain.Category{}
	err := row.Scan(&c.ID, &c.Title, &c.Description, &c.Thumbnail, &c.Status, &c.UpdatedAt, &c.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}
	return c, nil
}

func (m *mysqlCateRepo) Store(ctx context.Context, c *domain.Category) error {
	query := `INSERT INTO categories (title , description, thumbnail, status, updated_at, created_at)
				VALUES (?, ?, ?, ?, ?, ?)`

	res, err := m.db.ExecContext(ctx, query, c.Title, c.Description, c.Thumbnail, c.Status, c.UpdatedAt, c.CreatedAt)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	c.ID = id

	return nil
}

func (m *mysqlCateRepo) Update(ctx context.Context, c *domain.Category) error {
	query := `UPDATE categories SET
				title = ?,
				description = ?,
				thumbnail = ?,
				status = ?,
				updated_at = ?
				WHERE id = ?`

	_, err := m.db.ExecContext(ctx, query, c.Title, c.Description, c.Thumbnail, c.Status, c.UpdatedAt, c.ID)

	return err
}

func (m *mysqlCateRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE categories SET
				status = ?
				WHERE id = ?`

	_, err := m.db.ExecContext(ctx, query, domain.CategoryStatusInactive, id)

	return err
}
