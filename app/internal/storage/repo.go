package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Common repository errors.
var (
	ErrNotFound        = errors.New("not found")
	ErrEmailTaken      = errors.New("email already registered")
	ErrInvalidPassword = errors.New("invalid credentials")
)

type User struct {
	ID           int64
	FullName     string
	Email        string
	Phone        string
	PasswordHash string
	Role         string // "customer" | "admin"
	Status       string // "active" | "suspended"
}

type Template struct {
	ID          string
	Name        string
	Category    string
	Description string
	Thumbnail   string
	PackageID   string
	IsActive    bool
	ConfigJSON  []byte
	SortOrder   int32
}

type PackageRow struct {
	ID           string
	Name         string
	Tagline      string
	Description  string
	PriceCents   int64
	DurationDays int32
	MaxPhotos    int32
	MaxVideos    int32
	MaxEvents    int32
	ComingSoon   bool
	Available    bool
}

// CreateUser inserts a user and returns its ID. Returns ErrEmailTaken on conflict.
func (db *DB) CreateUser(ctx context.Context, fullName, email, phone, passwordHash, role string) (int64, error) {
	var id int64
	err := db.Pool.QueryRow(ctx,
		`INSERT INTO users (full_name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		fullName, email, phone, passwordHash, role).Scan(&id)
	if err != nil {
		var pe *pgconn.PgError
		if errors.As(err, &pe) {
			return 0, ErrEmailTaken
		}
		return 0, fmt.Errorf("create user: %w", err)
	}
	return id, nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := db.Pool.QueryRow(ctx,
		`SELECT id, full_name, email, phone, password_hash, role, status
		 FROM users WHERE email = $1`, email)
	u := &User{}
	err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.Phone, &u.PasswordHash, &u.Role, &u.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

func (db *DB) GetUserByID(ctx context.Context, id int64) (*User, error) {
	row := db.Pool.QueryRow(ctx,
		`SELECT id, full_name, email, phone, password_hash, role, status
		 FROM users WHERE id = $1`, id)
	u := &User{}
	err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.Phone, &u.PasswordHash, &u.Role, &u.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}

func (db *DB) ListTemplates(ctx context.Context, packageID, category, query string) ([]Template, error) {
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, category, description, thumbnail, package_id, is_active, config, sort_order
		 FROM templates
		 WHERE is_active
		   AND ($1 = '' OR package_id = $1)
		   AND ($2 = '' OR category = $2)
		   AND ($3 = '' OR lower(name) LIKE '%'||lower($3)||'%' OR lower(description) LIKE '%'||lower($3)||'%')
		 ORDER BY sort_order, id`,
		packageID, category, query)
	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	defer rows.Close()
	var out []Template
	for rows.Next() {
		var t Template
		if err := rows.Scan(&t.ID, &t.Name, &t.Category, &t.Description, &t.Thumbnail,
			&t.PackageID, &t.IsActive, &t.ConfigJSON, &t.SortOrder); err != nil {
			return nil, fmt.Errorf("scan template: %w", err)
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (db *DB) GetTemplate(ctx context.Context, id string) (*Template, error) {
	row := db.Pool.QueryRow(ctx,
		`SELECT id, name, category, description, thumbnail, package_id, is_active, config, sort_order
		 FROM templates WHERE id = $1`, id)
	t := &Template{}
	err := row.Scan(&t.ID, &t.Name, &t.Category, &t.Description, &t.Thumbnail,
		&t.PackageID, &t.IsActive, &t.ConfigJSON, &t.SortOrder)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get template: %w", err)
	}
	return t, nil
}

func (db *DB) ListPackages(ctx context.Context) ([]PackageRow, error) {
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, tagline, description, price_cents, duration_days,
		        max_photos, max_videos, max_events, coming_soon, available
		 FROM packages ORDER BY sort_order`)
	if err != nil {
		return nil, fmt.Errorf("list packages: %w", err)
	}
	defer rows.Close()
	var out []PackageRow
	for rows.Next() {
		var p PackageRow
		if err := rows.Scan(&p.ID, &p.Name, &p.Tagline, &p.Description, &p.PriceCents,
			&p.DurationDays, &p.MaxPhotos, &p.MaxVideos, &p.MaxEvents, &p.ComingSoon, &p.Available); err != nil {
			return nil, fmt.Errorf("scan package: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
