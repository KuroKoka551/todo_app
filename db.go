package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type database struct {
	db *sql.DB
}

var errNotFound = errors.New("Задача не найдена.")

func newDB(cfg config) *database {
	dbRaw, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName))
	if err != nil {
		log.Fatal(err)
	}
	db := &database{db: dbRaw}
	if err = db.Init(); err != nil {
		log.Fatal(err)
	}
	return db
}

func (d *database) Close() error {
	return d.db.Close()
}

func (d *database) Init() error {
	err := d.db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}
	_, err = d.db.Exec(
		"CREATE TABLE IF NOT EXISTS tasks " +
			"(id SERIAL PRIMARY KEY, title TEXT, description TEXT," +
			"due_date TIMESTAMP WITH TIME ZONE, created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()," +
			"updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW())")
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	return nil
}

func (d *database) GetTasks() ([]dbTask, error) {
	rows, err := d.db.Query("SELECT id, title, description, due_date, created_at, updated_at FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	defer rows.Close()

	todos := make([]dbTask, 0, 10)
	for rows.Next() {
		var t dbTask
		if err := rows.Scan(
			&t.ID, &t.Title,
			&t.Description, &t.DueDate,
			&t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan: %v", err)
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func (d *database) AddTask(t userTask) error {
	_, err := d.db.Exec(
		"INSERT INTO tasks (title, description, due_date) VALUES ($1, $2, $3)",
		t.Title, t.Description, t.DueDate)
	if err != nil {
		return fmt.Errorf("failed to insert: %v", err)
	}
	return nil
}

func (d *database) GetTask(id int) (dbTask, error) {
	var t dbTask
	err := d.db.QueryRow("SELECT id, title, description, due_date, created_at, updated_at FROM tasks WHERE id = $1", id).Scan(
		&t.ID, &t.Title, &t.Description, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return t, errNotFound
	} else if err != nil {
		return t, fmt.Errorf("failed to query: %v", err)
	}
	return t, nil
}

func (d *database) UpdateTask(id int, t userTask) (dbTask, error) {
	var dbT dbTask
	err := d.db.QueryRow(
		"UPDATE tasks SET title = $1, description = $2, due_date = $3, updated_at = NOW() WHERE id = $4 RETURNING id, title, description, due_date, created_at, updated_at",
		t.Title, t.Description, t.DueDate, id).Scan(&dbT.ID, &dbT.Title, &dbT.Description, &dbT.DueDate, &dbT.CreatedAt, &dbT.UpdatedAt)
	if err == sql.ErrNoRows {
		return dbT, errNotFound
	} else if err != nil {
		return dbT, fmt.Errorf("failed to query: %v", err)
	}
	return dbT, nil
}

func (d *database) DeleteTask(id int) error {
	res, err := d.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if n == 0 {
		return errNotFound
	}
	return nil
}
