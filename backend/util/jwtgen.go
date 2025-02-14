package util

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SessionTokenClaims struct {
	Account         string `json:"acc"`  // Account id of the connecting client
	PermissionLevel uint   `json:"plvl"` // Permission level of the account

	jwt.RegisteredClaims
}

// Generate a session token for an account
func SessionToken(account uuid.UUID, permLevel uint) (string, error) {

	// Create jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, SessionTokenClaims{

		Account:         account.String(),
		PermissionLevel: permLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(os.Getenv("MAGIC_JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
