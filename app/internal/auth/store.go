package auth

import (
"context"
"database/sql"
"errors"
"time"

"github.com/alexedwards/scs/v2"
)

// pgStore implements scs.SessionStore backed by the PostgreSQL `sessions` table.
// Session tokens are stored as SHA-256 hashes so a database leak cannot be
// replayed as a valid cookie.
type pgStore struct {
db    *sql.DB
table string
}

func (s *pgStore) Find(ctx context.Context, token string) ([]byte, error) {
var data []byte
err := s.db.QueryRowContext(ctx,
"SELECT data FROM "+s.table+" WHERE id = $1 AND expires_at > now()", hashToken(token)).Scan(&data)
if errors.Is(err, sql.ErrNoRows) {
return nil, nil
}
return data, err
}

func (s *pgStore) Commit(ctx context.Context, token string, data []byte, expiry time.Time) error {
_, err := s.db.ExecContext(ctx,
`INSERT INTO `+s.table+` (id, user_id, data, expires_at)
 VALUES ($1, NULL, $2, $3)
 ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data, expires_at = EXCLUDED.expires_at`,
hashToken(token), data, expiry.UTC())
return err
}

func (s *pgStore) Delete(ctx context.Context, token string) error {
_, err := s.db.ExecContext(ctx, "DELETE FROM "+s.table+" WHERE id = $1", hashToken(token))
return err
}

var _ scs.SessionStore = (*pgStore)(nil)
