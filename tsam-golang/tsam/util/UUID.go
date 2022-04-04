package util

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
)

// // Future for uuid pkg #Niranjan
// type UUID uuid.UUID

// func (id UUID) IsValid() bool {
// 	if uuid.UUID(id) == uuid.Nil {
// 		return false
// 	}
// 	return true
// }

// func New() UUID {
// 	return UUID(uuid.NewV4())
// }

// GenerateUUID Return uuid
func GenerateUUID() uuid.UUID {
	return uuid.NewV4()
}

// GetCurrentDate In YYYY-MM-DD format
func GetCurrentDate() *string {
	date := time.Now().Format("2006-01-02")
	return &date
}

// ParseUUID Parse uuid From String
func ParseUUID(input string) (uuid.UUID, error) {
	if len(input) == 0 {
		return uuid.Nil, errors.NewValidationError("Empty ID")
	}
	id, err := uuid.FromString(input)
	if err != nil {
		return uuid.Nil, errors.NewValidationError(input + ": Invalid ID")
	}
	return id, nil
}

// IsUUIDValid Check weather UUID Valid Or Not
func IsUUIDValid(id uuid.UUID) bool {
	if id == uuid.Nil {
		return false
	}
	return true
}

// GetCredentialIDFromToken Extract Token From Request & return ID From Token
func GetCredentialIDFromToken(r *http.Request, JWTKey string) (uuid.UUID, error) {
	tokenStr, err := request.HeaderExtractor{"token"}.ExtractToken(r)
	if err != nil {
		return uuid.Nil, errors.NewHTTPError("empty token", http.StatusBadRequest)
	}
	token, err := jwt.Parse(tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
			}
			return []byte(JWTKey), nil
		})
	if err != nil {
		return uuid.Nil, errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, _ := claims["credentialID"].(string)
		credentialID, err := ParseUUID(id)
		if err != nil {
			return uuid.Nil, errors.NewHTTPError(err.Error(), http.StatusBadRequest)
		}
		return credentialID, nil
	}
	return uuid.Nil, errors.NewHTTPError("Invalid token.", http.StatusBadRequest)
}
