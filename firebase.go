package middleware

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func AuthWithFirebase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opt := option.WithCredentialsFile("credentials.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error on credentials " + err.Error()))
			return
		}

		auth, err2 := app.Auth(context.Background())
		if err2 != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error on credentials " + err2.Error()))
			return
		}

		header := r.Header.Get("HeaderAuthorization")
		idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
		_, err = auth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte("Error getting the token.\n" + err.Error()))
			return
		}

		//compara uid
		//crear cliente
		//crear libreria

		next.ServeHTTP(w, r)
	})
}
