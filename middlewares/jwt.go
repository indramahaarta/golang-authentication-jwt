package middlewares

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/indramahaarta/golang-authentication-jwt/config"
	"github.com/indramahaarta/golang-authentication-jwt/helpers"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := map[string]interface{}{
					"message": "Unauthorized",
				}
				helpers.ResponseJSON(w, http.StatusUnauthorized, response)
				return 
			}
		}

		// mengambil token value
		tokenString := c.Value

		// claims
		claims := &config.JWTClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWK_KEY, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				response := map[string]interface{}{
					"message": "Unauthorized",
				}
				helpers.ResponseJSON(w, http.StatusUnauthorized, response)
				return

			case jwt.ValidationErrorExpired:
				response := map[string]interface{}{
					"message": "Unauthorized, token is expired",
				}
				helpers.ResponseJSON(w, http.StatusUnauthorized, response)
				return

			default:
				response := map[string]interface{}{
					"message": "Unauthorized",
				}
				helpers.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}

		if !token.Valid {
			response := map[string]interface{}{
				"message": "Unauthorized, token is expired",
			}
			helpers.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		next.ServeHTTP(w, r)
	})
}
