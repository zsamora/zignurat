package utils

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// DropInvalidSessionCookie removes a session cookie that can no longer be
// decoded, e.g. because it was signed with an old SESSION_SECRET.
// Without this, gin-contrib/sessins returns a nil session and any access to
// it panics (HTTP 500 on every page until the user clears their cookies).
// Register it BEFORE sessions.Sessions.
func DropInvalidSessionCookie(name string, store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := c.Request.Cookie(name); err == nil {
			// Decode on a copy of the request, so the session registry that
			// gorilla/sessions attaches to the request is not reused later.
			probe := c.Request.Clone(c.Request.Context())
			if _, err := store.Get(probe, name); err != nil {
				log.Printf("!! Dropping invalid %q cookie: %s", name, err)
				removeRequestCookie(c.Request, name)
				// Tell the browser to delete it too
				http.SetCookie(c.Writer, &http.Cookie{Name: name, Value: "", Path: "/", MaxAge: -1})
			}
		}
		c.Next()
	}
}

// removeRequestCookie rebuilds the Cookie header without the named cookie.
func removeRequestCookie(r *http.Request, name string) {
	cookies := r.Cookies()
	r.Header.Del("Cookie")
	for _, ck := range cookies {
		if ck.Name != name {
			r.AddCookie(ck)
		}
	}
}
