package repository

import (
	"database/sql"
	"fmt"
	"os"

	"todocli/internal/model"
	"todocli/internal/service"

	_ "github.com/lib/pq"
)

type DBExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type PostgresRepo struct {
	db DBExecutor
}

// Конструктор без транзакции
func NewPostgresRepo(db DBExecutor) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Конструктор с подключением
func NewPostgresRepository() (*PostgresRepo, error) {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:password@localhost:5432/tododb?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return NewPostgresRepo(db), nil
}

func (r *PostgresRepo) Save(t *model.Task) error {
	query := `INSERT INTO tasks (title, description, status, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(query, t.Title, t.Description, t.Status.String(), t.CreatedAt).Scan(&t.ID)
}

func (r *PostgresRepo) Update(t *model.Task) error {
	query := `UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4`
	_, err := r.db.Exec(query, t.Title, t.Description, t.Status.String(), t.ID)
	return err
}

func (r *PostgresRepo) GetAll() []*model.Task {
	query := `SELECT id, title, description, status, created_at FROM tasks`
	rows, err := r.db.Query(query)
	if err != nil {
		fmt.Println("Query error:", err)
		return nil
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		var t model.Task
		var statusStr string
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &statusStr, &t.CreatedAt); err == nil {
			t.Status, _ = model.ParseStatus(statusStr)
			tasks = append(tasks, &t)
		}
	}
	return tasks
}

func (r *PostgresRepo) FindByID(id int) *model.Task {
	query := `SELECT id, title, description, status, created_at FROM tasks WHERE id=$1`
	row := r.db.QueryRow(query, id)

	var t model.Task
	var statusStr string
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &statusStr, &t.CreatedAt); err != nil {
		return nil
	}
	t.Status, _ = model.ParseStatus(statusStr)
	return &t
}

func (r *PostgresRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM tasks WHERE id=$1`, id)
	return err
}

// Транзакции
func (r *PostgresRepo) WithTx(fn func(service.TaskRepository) error) error {
	realDB, ok := r.db.(*sql.DB)
	if !ok {
		return fmt.Errorf("WithTx requires *sql.DB")
	}

	tx, err := realDB.Begin()
	if err != nil {
		return err
	}

	txRepo := &PostgresRepo{db: tx}
	if err := fn(txRepo); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
