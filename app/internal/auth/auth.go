package auth

import (
"context"
"crypto/rand"
"crypto/sha256"
"crypto/subtle"
"database/sql"
"encoding/hex"
"errors"
"fmt"
"net/http"
"time"

"github.com/alexedwards/scs/v2"
_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver
"golang.org/x/crypto/bcrypt"
)

const (
RoleCustomer = "customer"
RoleAdmin    = "admin"

userIDKey    = "user_id"
userRoleKey  = "role"
userNameKey  = "user_name"
csrfTokenKey = "csrf_token"
)

// Manager wraps the SCS session manager and exposes auth helpers.
type Manager struct {
Session *scs.SessionManager
}

// NewManager creates a session manager backed by the PostgreSQL sessions table.
func NewManager(db *sql.DB, cookieName string, ttl time.Duration, secure bool) (*Manager, error) {
sm := scs.New()
sm.Store = &pgStore{db: db, table: "sessions"}
sm.Cookie.Name = cookieName
sm.Cookie.HttpOnly = true
sm.Cookie.SameSite = http.SameSiteLaxMode
sm.Cookie.Secure = secure
sm.Lifetime = ttl
sm.CleanupInterval = 10 * time.Minute
return &Manager{Session: sm}, nil
}

// OpenSQLDB opens a database/sql handle to Postgres using the pgx stdlib driver.
func OpenSQLDB(ctx context.Context, url string) (*sql.DB, error) {
db, err := sql.Open("pgx", url)
if err != nil {
return nil, fmt.Errorf("auth: open sql db: %w", err)
}
pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
if err := db.PingContext(pingCtx); err != nil {
db.Close()
return nil, fmt.Errorf("auth: ping sql db: %w", err)
}
return db, nil
}

func HashPassword(password string) (string, error) {
b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
return "", fmt.Errorf("auth: hash: %w", err)
}
return string(b), nil
}

func CheckPassword(hash, password string) bool {
return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Login establishes a new session for the user (session-fixation safe).
func (m *Manager) Login(ctx context.Context, userID int64, role, fullName string) error {
if err := m.Session.RenewToken(ctx); err != nil {
return fmt.Errorf("auth: renew token: %w", err)
}
m.Session.Put(ctx, userIDKey, userID)
m.Session.Put(ctx, userRoleKey, role)
m.Session.Put(ctx, userNameKey, fullName)
token := make([]byte, 32)
if _, err := rand.Read(token); err != nil {
return fmt.Errorf("auth: csrf rand: %w", err)
}
m.Session.Put(ctx, csrfTokenKey, hex.EncodeToString(token))
return nil
}

func (m *Manager) Logout(ctx context.Context) error {
return m.Session.Destroy(ctx)
}

func (m *Manager) CurrentUserID(ctx context.Context) (int64, bool) {
id := m.Session.GetInt64(ctx, userIDKey)
return id, id != 0
}

func (m *Manager) CurrentRole(ctx context.Context) string {
return m.Session.GetString(ctx, userRoleKey)
}

func (m *Manager) CurrentName(ctx context.Context) string {
return m.Session.GetString(ctx, userNameKey)
}

func (m *Manager) IsAuthenticated(ctx context.Context) bool {
_, ok := m.CurrentUserID(ctx)
return ok
}

func (m *Manager) IsAdmin(ctx context.Context) bool {
return m.CurrentRole(ctx) == RoleAdmin
}

// CSRFToken returns the per-session CSRF token, creating one if missing.
func (m *Manager) CSRFToken(ctx context.Context) string {
tok := m.Session.GetString(ctx, csrfTokenKey)
if tok == "" {
buf := make([]byte, 32)
_, _ = rand.Read(buf)
tok = hex.EncodeToString(buf)
m.Session.Put(ctx, csrfTokenKey, tok)
}
return tok
}

// ValidateCSRF compares tokens in constant time.
func (m *Manager) ValidateCSRF(ctx context.Context, submitted string) bool {
expected := m.CSRFToken(ctx)
if expected == "" || submitted == "" {
return false
}
return subtle.ConstantTimeCompare([]byte(expected), []byte(submitted)) == 1
}

// hashToken converts a session token to its SHA-256 digest for storage.
func hashToken(token string) []byte {
sum := sha256.Sum256([]byte(token))
return sum[:]
}

var ErrInvalidCredentials = errors.New("invalid email or password")
