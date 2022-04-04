package security

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	uuid "github.com/satori/go.uuid"
	conf "github.com/techlabs/swabhav/tsam/config"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/util"
)

// GenerateToken take userID, email, tablename as Role  Return Token
func (auth *Authentication) GenerateToken(credentialID string, loginID string, email string, roleID string) (string, error) {
	// Create a claims map
	// claims based on which token should be created
	claims := jwt.MapClaims{
		"credentialID": credentialID,
		"loginID":      loginID,
		"emailID":      email,
		"roleID":       roleID,
		"exp":          time.Now().Add(time.Hour * 900000).Unix(),
	}

	// NewWithClaims returns token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// access token string based on token
	tokenString, err := token.SignedString([]byte(auth.Config.GetString(conf.JWTKey)))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return "", errors.NewHTTPError("unable to generate Token", http.StatusInternalServerError)
	}
	return tokenString, nil
}

// ExtractIDFromToken will check the token.Claims for entityName and extract the ID from token.
// 	eg: jwt.Mapclaims["credentialID"] when entityName = "credentialID"
func (auth *Authentication) ExtractIDFromToken(r *http.Request, entityName string) (uuid.UUID, error) {
	tokenStr, err := request.HeaderExtractor{"token"}.ExtractToken(r)
	if err != nil {
		return uuid.Nil, errors.NewHTTPError("empty token", http.StatusBadRequest)
	}
	token, err := jwt.Parse(tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
			}
			return []byte(auth.Config.GetString(conf.JWTKey)), nil
		})
	if err != nil {
		return uuid.Nil, errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, _ := claims[entityName].(string)
		entityID, err := util.ParseUUID(id)
		if err != nil {
			return uuid.Nil, err
		}
		return entityID, nil
	}
	return uuid.Nil, errors.NewHTTPError("Invalid token.", http.StatusBadRequest)
}

// ExtractCredentialIDFromToken will extract the credentialID from token.
func (auth *Authentication) ExtractCredentialIDFromToken(r *http.Request) (uuid.UUID, error) {
	tokenStr, err := request.HeaderExtractor{"token"}.ExtractToken(r)
	if err != nil {
		return uuid.Nil, errors.NewHTTPError("empty token", http.StatusBadRequest)
	}
	token, err := jwt.Parse(tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
			}
			return []byte(auth.Config.GetString(conf.JWTKey)), nil
		})
	if err != nil {
		return uuid.Nil, errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, _ := claims["credentialID"].(string)
		credentialID, err := util.ParseUUID(id)
		if err != nil {
			return uuid.Nil, err
		}
		return credentialID, nil
	}
	return uuid.Nil, errors.NewHTTPError("Invalid token.", http.StatusBadRequest)
}
