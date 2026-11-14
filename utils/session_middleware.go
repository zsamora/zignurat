package utils

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func DropInvalidSessionCookie(name string, store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := c.Request.Cookie(name); err == nil {
			probe := c.Request.Clone(c.Request.Context())
			if _, err := store.Get(probe, name); err != nil {
				log.Printf("!! Dropping invalid %q cookie: %s", name, err)
				removeRequestCookie(c.Request, name)
				http.SetCookie(c.Writer, &http.Cookie{Name: name, Value: "", Path: "/", MaxAge: -1})
			}
		}
		c.Next()
	}
}

func removeRequestCookie(r *http.Request, name string) {
	cookies := r.Cookies()
	r.Header.Del("Cookie")
	for _, ck := range cookies {
		if ck.Name != name {
			r.AddCookie(ck)
		}
	}
}
