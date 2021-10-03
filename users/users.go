package users

import (
	"time"

	"banking-app/database"
	"banking-app/helpers"
	"banking-app/interfaces"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func prepareToken(user *interfaces.User) string {
	// Sign Token
	tokenContent := jwt.MapClaims {
		"user_id": user.ID,
		"expiry": time.Now().Add(time.Minute ^ 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	return token
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount, withToken bool) map[string]interface{} {
	// Atur Respon
	responseUser := &interfaces.ResponseUser{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		Accounts: accounts,
	}

	// Prepare Response
	var response = map[string]interface{}{"message": "Semuanya Sukses!"}
	if withToken {
		var token = prepareToken(user)
		response["jwt"] = token
	}
	response["data"] = responseUser

	return response
}

func Login(username string, password string) map[string]interface{} {
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: password, Valid: "password"},
		},
	)

	if valid {
		// Koneksi ke DB
		// db := helpers.ConnectDB()
		user := &interfaces.User{}
		if database.DB.Where("username = ?", username).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User tidak ditemukan!"}
		}

		// Verifikasi Password
		passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if passwordErr == bcrypt.ErrMismatchedHashAndPassword && passwordErr != nil {
			return map[string]interface{}{"message": "Password anda salah!"}
		}

		// Cari akun user
		accounts := []interfaces.ResponseAccount{}
		database.DB.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

		// defer db.Close()

		var response = prepareResponse(user, accounts, true)

		return response
	} else {
		return map[string]interface{}{"message": "Username atau Password Tidak Valid"}
	}
}

func Register(username string, email string, password string) map[string]interface{} {
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: password, Valid: "password"},
		},
	)

	if valid {
		// db := helpers.ConnectDB()
		generatePassword := helpers.HashAndSalt([]byte(password))
		user := &interfaces.User{Username: username, Email: email, Password: generatePassword}
		// db.Create(&user)
		database.DB.Create(&user)

		account := &interfaces.Account{Type: "Daily Account", Name: string(username + "'s" + " account"), Balance: 0, UserID: user.ID}
		// db.Create(&account)
		database.DB.Create(&account)

		// defer db.Close()

		accounts := []interfaces.ResponseAccount{}
		responseAccount := interfaces.ResponseAccount{ID: account.ID, Name: account.Name, Balance: int(account.Balance)}
		accounts = append(accounts, responseAccount)
		var response = prepareResponse(user, accounts, true)

		return response
	} else {
		return map[string]interface{}{"message": "Username atau Password Tidak Valid"}
	}
}

func GetUser(id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken(id, jwt)

	if isValid {
		// db := helpers.ConnectDB()
		
		user := &interfaces.User{}
		if database.DB.Where("id = ?", id).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User tidak ditemukan!"}
		}
		accounts := []interfaces.ResponseAccount{}
		database.DB.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

		// defer db.Close()

		var response = prepareResponse(user, accounts, false)
		return response

	} else {
		return map[string]interface{}{"message": "Maaf Token Tidak Valid"}
	}
}