package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var (
	firebaseClient *FirebaseClient
	once           sync.Once
)

type FirebaseClient struct {
	AuthClient *auth.Client
}

func Newfirebase() *FirebaseClient {
	once.Do(InitFirebase)

	return firebaseClient
}

func InitFirebase() {

	opt := option.WithCredentialsFile("credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Error connecting to firebase" + err.Error())
	}

	auth, err2 := app.Auth(context.Background())
	if err2 != nil {
		log.Println("Error connecting to firebase" + err2.Error())
	}

	firebaseClient = &FirebaseClient{
		AuthClient: auth,
	}
}

func AuthWithFirebase(next http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("HeaderAuthorization")
			idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
			_, err := firebaseClient.AuthClient.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

/* func (f *FirebaseClient) Close() error {
	if firebaseClient == nil {
		return nil
	}

	return f
} */
