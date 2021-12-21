package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
)

func CurrentUser(ur app.UserRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := r.Header.Get("Authorization")

			if len(token) > 0 {
				splitToken := strings.Split(token, "Bearer ")
				token = splitToken[1]

				claims, err := utils.GetTokenClaims(token)

				if err != nil {
					utils.SendJsonError(w, "invalid token", http.StatusUnauthorized)
					return
				}

				if claims != nil {
					user, err := ur.Get(claims.Uid)

					if err == nil {
						ctx := context.WithValue(r.Context(), "user", user)

						r = r.WithContext(ctx)
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
