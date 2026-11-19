package web

import (
"net/http"
"strings"

"github.com/sukamoto/app/internal/auth"
)

// upper returns s upper-cased (used for package badges).
func upper(s string) string { return strings.ToUpper(s) }

// ifElse is a tiny ternary helper usable inside templ expressions.
func ifElse(cond bool, a, b string) string {
if cond {
return a
}
return b
}

// BuildPageData assembles the common page payload from the request context.
func BuildPageData(r *http.Request, m *auth.Manager, title, desc string) PageData {
pd := PageData{Title: title, Description: desc}
if m != nil {
pd.CSRF = m.CSRFToken(r.Context())
pd.Nav = NavProps{
Authenticated: m.IsAuthenticated(r.Context()),
IsAdmin:       m.IsAdmin(r.Context()),
Name:          m.CurrentName(r.Context()),
}
}
return pd
}
