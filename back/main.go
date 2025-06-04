package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	fbAppID     string
	fbAppSecret string
	jwtSecret   []byte
	db          *sql.DB
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	fbAppID = os.Getenv("FB_APP_ID")
	fbAppSecret = os.Getenv("FB_APP_SECRET")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/test_app_fb_auth")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/auth/facebook/callback", facebookCallbackHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func facebookCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	accessToken, err := getFacebookAccessToken(code)
	if err != nil {
		http.Error(w, "Failed to get access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fbUser, err := getFacebookUser(accessToken)
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE facebook_id = ?)", fbUser.ID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		_, err = db.Exec(`
			UPDATE users SET name=?, email=?, first_name=?, last_name=?, profile_pic=?, gender=?, locale=?
			WHERE facebook_id=?
		`, fbUser.Name, fbUser.Email, fbUser.FirstName, fbUser.LastName, fbUser.Picture.Data.Url, fbUser.Gender, fbUser.Locale, fbUser.ID)
	} else {
		_, err = db.Exec(`
			INSERT INTO users (facebook_id, name, email, first_name, last_name, profile_pic, gender, locale)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, fbUser.ID, fbUser.Name, fbUser.Email, fbUser.FirstName, fbUser.LastName, fbUser.Picture.Data.Url, fbUser.Gender, fbUser.Locale)
	}
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"facebook_id": fbUser.ID,
		"name":        fbUser.Name,
		"email":       fbUser.Email,
		"exp":         time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Token generation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to frontend with token in query param
	redirectURL := fmt.Sprintf("http://localhost:5173/?token=%s", url.QueryEscape(tokenString))
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func getFacebookAccessToken(code string) (string, error) {
	params := url.Values{}
	params.Set("client_id", fbAppID)
	params.Set("client_secret", fbAppSecret)
	params.Set("redirect_uri", "http://localhost:8080/auth/facebook/callback")
	params.Set("code", code)

	resp, err := http.Get("https://graph.facebook.com/v17.0/oauth/access_token?" + params.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.AccessToken, nil
}

type FacebookUser struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   struct {
		Data struct {
			Url string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
	Gender *string `json:"gender,omitempty"`
	Locale *string `json:"locale,omitempty"`
}

func getFacebookUser(accessToken string) (*FacebookUser, error) {
	endpoint := "https://graph.facebook.com/me?fields=id,name,email,first_name,last_name,picture,locale,gender&access_token=" + url.QueryEscape(accessToken)

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user FacebookUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
