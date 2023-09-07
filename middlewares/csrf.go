package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"

	"github.com/bookandmusic/docs/global"
)

func CSRFMiddleware() gin.HandlerFunc {
	CSRF := csrf.Protect(
		[]byte(global.GVA_CONFIG.Server.SecretKey),
		csrf.RequestHeader("Authenticity-Token"),
		csrf.FieldName("authenticity_token"),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message": "csrf token error"}`))
		})),
	)

	return adapter.Wrap(CSRF)
}
