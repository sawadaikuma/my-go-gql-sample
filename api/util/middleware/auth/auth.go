package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// A stand-in for our database backed user object
type User struct {
	ID      string
	Name    string
	IsAdmin bool
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			c, err := r.Cookie("auth-cookie")

			cookie := &http.Cookie{
				Name:  "hoge", // ここにcookieの名前を記述
				Value: "bar",  // ここにcookieの値を記述
			}
			cookie2 := &http.Cookie{
				Name:  "auth-cookie", // ここにcookieの名前を記述
				Value: "bar",         // ここにcookieの値を記述
			}
			// 2
			http.SetCookie(w, cookie)
			http.SetCookie(w, cookie2)

			if c == nil {
				panic(fmt.Errorf("auth errors"))
			}
			// Allow unauthenticated users in
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			fmt.Println("[cokkie]=>", c)

			userId, err := validateAndGetUserID(c)
			if err != nil {
				fmt.Println("[Auth]Authenticate Error(403)")
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}

			// get the user from the database
			user, err := getUserByID(db, userId)
			if err != nil {
				panic(err)
			}

			fmt.Println("[Auth] User=>", user)

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func validateAndGetUserID(*http.Cookie) (userID int, err error) {
	// return 0, fmt.Errorf("401")
	return 1, nil
}

func getUserByID(db *sql.DB, userID int) (*User, error) {
	user := User{}
	err := db.QueryRow("SELECT id, name from user limit 1").Scan(&user.ID, &user.Name)
	fmt.Println("[user]=>", user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *User {
	raw, _ := ctx.Value(userCtxKey).(*User)
	return raw
}
