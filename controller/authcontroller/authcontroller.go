package authcontroller

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/indramahaarta/golang-authentication-jwt/config"
	"github.com/indramahaarta/golang-authentication-jwt/database"
	"github.com/indramahaarta/golang-authentication-jwt/helpers"
	"github.com/indramahaarta/golang-authentication-jwt/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	// mengambil body dari request
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response := map[string]string{
			"message": "success",
		}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	// hash pasword menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err.Error())
		response := map[string]string{
			"message": "password failed to hashed",
		}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	user.Password = string(hashedPassword)

	// insert ke database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = database.DB.ExecContext(ctx, `INSERT INTO "user"(user_name, full_name, password) values ($1, $2, $3)`, user.UserName, user.FullName, user.Password)

	if err != nil {
		log.Println(err.Error())
		response := map[string]string{
			"message": "cannot insert into database",
		}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// response success
	response := map[string]string{
		"message": "success",
	}
	helpers.ResponseJSON(w, http.StatusAccepted, response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// mengambil body dari request
	var userInput models.User
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		response := map[string]string{
			"message": "success",
		}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	// ambil data user berdasarkan username
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	row := database.DB.QueryRowContext(ctx, `select id, full_name, user_name, password from "user" where user_name=$1`, userInput.UserName)
	err = row.Scan(&user.ID, &user.FullName, &user.UserName, &user.Password)

	switch {
	case err == sql.ErrNoRows:
		log.Println("Data username tidak ditemukan")
		response := map[string]string{
			"message": "Data username tidak ditemukan",
		}
		helpers.ResponseJSON(w, http.StatusNotFound, response)
		return
	case err != nil:
		log.Println(err.Error())
		response := map[string]string{
			"message": "Internal server error",
		}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{
			"message": "Password anda salah",
		}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// proses pembuatan claims jwt
	expTime := time.Now().Add(time.Minute * 30)
	claims := &config.JWTClaim{
		Username: user.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// mendeklarasikan algoritma untuk signed in
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signed token
	token, err := tokenAlgo.SignedString(config.JWK_KEY)
	if err != nil {
		log.Println(err.Error())
		response := map[string]string{
			"message": "Internal server error",
		}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// set token ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	response := map[string]string{
		"message": "Seccessfully Login",
	}
	helpers.ResponseJSON(w, http.StatusAccepted, response)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	// hapus token ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge: -1,
	})

	response := map[string]string{
		"message": "Seccessfully Logout",
	}
	helpers.ResponseJSON(w, http.StatusAccepted, response)
}
