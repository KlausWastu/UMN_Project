package authController

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jeypc/go-jwt-mux/config"
	"github.com/jeypc/go-jwt-mux/helper"
	"github.com/jeypc/go-jwt-mux/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// menagambil inputan json
	var userInput models.User
	// mendapat json dari body
	decoderJSON := json.NewDecoder(r.Body)
	err := decoderJSON.Decode(&userInput)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// mengambil data user berdasarkan email
	var user models.User
	err = models.DB.Where("email = ?", userInput.Email).First(&user).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": "Email or password not found"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// cek password valid
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		response := map[string]string{"message": "Email or password not found"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// proses pembuatan token
	expTime := time.Now().Add(time.Minute * 15)
	claims := &config.JWTClaim{
		UserID:   user.UserID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.Fullname,
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// penggunaan algo untuk signin
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// sign token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)

	if err != nil {
		response := map[string]string{"message": "err.Error()"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// set token ke cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	response := map[string]string{"message": "login successfully"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// menagambil inputan json
	var userInput models.User
	// mendapat json dari body
	decoderJSON := json.NewDecoder(r.Body)
	err := decoderJSON.Decode(&userInput)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// hash pass menggunakan bcrypt

	hashpass, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashpass)

	// insert ke database
	err = models.DB.Create(&userInput).Error
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "success"}
	helper.ResponseJSON(w, http.StatusOK, response)

}

func Logout(w http.ResponseWriter, r *http.Request) {

	// set token ke cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "logout successfully"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {

	// Ambil token dari cookies dan parsing claims token
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value

	// Parsing token dan klaim
	claims := &config.JWTClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verifikasi apakah algoritma yang digunakan untuk sign token cocok
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.JWT_KEY, nil
	})

	// Jika token tidak valid atau ada kesalahan
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// menagambil inputan json
	var userInput models.ChangePassword

	// mendapat json dari body
	decoderJSON := json.NewDecoder(r.Body)
	err = decoderJSON.Decode(&userInput)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// pengecekan password lama
	var user models.User

	// pengeceekan user didatabase
	err = models.DB.Where("user_id = ?", claims.UserID).First(&user).Error
	if err != nil {
		response := map[string]string{"message": "user not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	// cek password valid
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.OldPassword))

	if err != nil {
		response := map[string]string{"message": "old password is incorrect"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// pengecekan password baru dengan konfirmasi password
	if userInput.NewPassword != userInput.ConfirmPassword {
		response := map[string]string{"message": "new password and confirm password not match"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// hash new password menggunakan bcrypt
	hashpass, _ := bcrypt.GenerateFromPassword([]byte(userInput.NewPassword), bcrypt.DefaultCost)
	userInput.NewPassword = string(hashpass)

	// insert ke database
	err = models.DB.Model(user).Where("user_id = ?", user.UserID).Update("password", &userInput.NewPassword).Error
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "change password successfully"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
