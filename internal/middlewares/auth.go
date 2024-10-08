package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"silent-notes/internal/types"
	"silent-notes/internal/utils"
	"time"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := r.Cookie("token")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Unauthorized", Error: "Unauthorized"}
			json.NewEncoder(w).Encode(res)
			return
		}

		claims, err := utils.VerifyJWT(token.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Unauthorized", Error: "Unauthorized"}
			json.NewEncoder(w).Encode(res)
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Unauthorized", Error: "Invalid expiration claim"}
			json.NewEncoder(w).Encode(res)
			return
		}

		if time.Now().After(time.Unix(int64(exp), 0)) {
			w.WriteHeader(http.StatusUnauthorized)
			res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Unauthorized", Error: "Cookie Expired"}
			json.NewEncoder(w).Encode(res)
			return
		}

		userId := claims["user_id"].(string)
		ctx := context.WithValue(r.Context(), types.UserIDKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
