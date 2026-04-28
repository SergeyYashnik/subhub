package postgres

import (
	"Test_EM/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
	Db *sqlx.DB
}

func New(dsn string) (*Repository, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to postgres: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging postgres: %v", err)
	}

	return &Repository{Db: db}, nil
}

func (r *Repository) Create(ctx context.Context, sub models.Subscription) error {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES (:service_name, :price, :user_id, :start_date, :end_date)
	`

	_, err := r.Db.NamedExecContext(ctx, query, sub)
	if err != nil {
		return fmt.Errorf("repository.Create: %w", err)
	}

	return nil
}

func (r *Repository) GetAll(ctx context.Context) ([]models.Subscription, error) {
	var subs []models.Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`

	err := r.Db.SelectContext(ctx, &subs, query)
	if err != nil {
		return nil, fmt.Errorf("repo.GetAll: %w", err)
	}

	return subs, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (models.Subscription, error) {
	var sub models.Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`

	err := r.Db.GetContext(ctx, &sub, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Subscription{}, fmt.Errorf("subscription not found")
		}
		return models.Subscription{}, fmt.Errorf("repo.GetByID: %w", err)
	}

	return sub, nil
}

func (r *Repository) Update(ctx context.Context, sub models.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET service_name = :service_name, 
		    price = :price, 
		    start_date = :start_date, 
		    end_date = :end_date 
		WHERE id = :id`

	result, err := r.Db.NamedExecContext(ctx, query, sub)
	if err != nil {
		return fmt.Errorf("repo.Update: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no subscription found to update")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("nothing to delete")
	}

	return nil
}

func (r *Repository) GetSum(ctx context.Context, userID string, serviceName string, from, to time.Time) (int, error) {
	var total int

	query := `
		SELECT COALESCE(SUM(price), 0) 
		FROM subscriptions 
		WHERE user_id = $1 
		  AND start_date >= $2 
		  AND start_date <= $3`

	if serviceName != "" {
		query += " AND service_name = $4"
		err := r.Db.GetContext(ctx, &total, query, userID, from, to, serviceName)
		return total, err
	}

	err := r.Db.GetContext(ctx, &total, query, userID, from, to)
	return total, err
}
