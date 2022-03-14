package middleware

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func InitFirebaseClient(fileName string) *auth.Client {

	opt := option.WithCredentialsFile(fileName)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {

		return nil
	}

	auth, err2 := app.Auth(context.Background())
	if err2 != nil {

		return nil
	}

	return auth
}

func AuthWithFirebase(client *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("HeaderAuthorization")
			idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
			_, err := client.VerifyIDToken(context.Background(), idToken)
			if err != nil {

				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
