package security

import (
	"net/http"
	"time"

	"github.com/techlabs/swabhav/tsam/config"
	conf "github.com/techlabs/swabhav/tsam/config"

	"github.com/jinzhu/gorm"

	"github.com/techlabs/swabhav/tsam/log"

	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// Authentication Provide Method AuthUser.
type Authentication struct {
	DB     *gorm.DB
	Config config.ConfReader
}

// NewAuthentication returns new instance of Authentication
func NewAuthentication(db *gorm.DB, config config.ConfReader) *Authentication {
	return &Authentication{
		DB:     db,
		Config: config,
	}
}

// ValidateToken verifies the user login.
func (auth *Authentication) ValidateToken(w http.ResponseWriter, r *http.Request) error {
	// log.NewLogger().Info("==============================ValidateToken call==============================")
	tokenStr, err := request.HeaderExtractor{"token"}.ExtractToken(r)
	if err != nil {
		// web.RespondError(w, errors.NewHTTPError("invalid Session ID", http.StatusBadRequest))
		return errors.NewHTTPError("Invalid Session ID", http.StatusBadRequest)
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenStr, &claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
			}
			return []byte(auth.Config.GetString(config.JWTKey)), nil
		})

	if err != nil {
		// checks if the jwt token has expired or not
		now := time.Now().Unix()
		if !claims.VerifyExpiresAt(now, false) {
			return errors.NewValidationError("Session expired! Please login again")
		}
		log.NewLogger().Error(err.Error())
		// web.RespondError(w, errors.NewHTTPError("Invalid Session ID", http.StatusInternalServerError))
		return errors.NewHTTPError("Invalid Session ID", http.StatusInternalServerError)
	}

	// prints all the claims
	// for key, val := range claims {
	// 	fmt.Printf("Key: %v, value: %v\n", key, val)
	// }

	// if token is valid then it will be redirected to the endpoint
	if token.Valid {
		credential := general.Credential{}
		id, _ := claims["credentialID"].(string)
		credential.ID, err = util.ParseUUID(id)
		if err != nil {
			// web.RespondError(w, errors.NewHTTPError("Invalid Login!", http.StatusForbidden))
			log.NewLogger().Error(err.Error())
			return err
		}
		if r.Method == "OPTION" {
			w.WriteHeader(http.StatusOK)
			return err
		}
		return nil
	}
	// returns error if token is not valid
	return errors.NewHTTPError("Invalid Session ID", http.StatusInternalServerError)
}

// TestMiddleware verifies the user login.
func TestMiddleware(w http.ResponseWriter, r *http.Request) error {
	// log.NewLogger().Info("==============================Middleware call==============================")
	tokenStr, err := request.HeaderExtractor{"token"}.ExtractToken(r)
	if err != nil {
		// web.RespondError(w, errors.NewHTTPError("invalid Session ID", http.StatusBadRequest))
		return errors.NewHTTPError("Invalid Session ID", http.StatusBadRequest)
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenStr, &claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
			}
			var str string
			return []byte(str), nil
		})

	if err != nil {
		// checks if the jwt token has expired or not
		now := time.Now().Unix()
		if !claims.VerifyExpiresAt(now, false) {
			return errors.NewValidationError("Session expired! Please login again")
		}
		log.NewLogger().Error(err.Error())
		// web.RespondError(w, errors.NewHTTPError("Invalid Session ID", http.StatusInternalServerError))
		return errors.NewHTTPError("Invalid Session ID", http.StatusInternalServerError)
	}

	// prints all the claims
	// for key, val := range claims {
	// 	fmt.Printf("Key: %v, value: %v\n", key, val)
	// }

	// if token is valid then it will be redirected to the endpoint
	if token.Valid {
		credential := general.Credential{}
		id, _ := claims["credentialID"].(string)
		credential.ID, err = util.ParseUUID(id)
		if err != nil {
			// web.RespondError(w, errors.NewHTTPError("Invalid Login!", http.StatusForbidden))
			log.NewLogger().Error(err.Error())
			return err
		}
		if r.Method == "OPTION" {
			w.WriteHeader(http.StatusOK)
			return err
		}
		return nil
	}
	// returns error if token is not valid
	return errors.NewHTTPError("Invalid Session ID", http.StatusInternalServerError)
}

// CheckSessionValidity checks if current session is valid
func (auth *Authentication) CheckSessionValidity(r *http.Request) error {

	tokenStr, err := request.HeaderExtractor{"token"}.ExtractToken(r)
	if err != nil {
		return errors.NewHTTPError("Invalid Session ID", http.StatusBadRequest)
	}

	claims := jwt.MapClaims{}

	// ParseWithClaims returns token and err
	_, err = jwt.ParseWithClaims(
		tokenStr, &claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
			}
			return []byte(auth.Config.GetString(conf.JWTKey)), nil
		})

	now := time.Now().Unix()
	if !claims.VerifyExpiresAt(now, false) {
		return errors.NewValidationError("Session expired! Please login again")
	}

	return nil

	// fmt.Println("Token ->", token)
	// // prints all the claims
	// for key, val := range claims {
	// 	fmt.Printf("Key: %v, value: %v\n", key, val)
	// }

	// fmt.Println("Valid ->", claims.VerifyExpiresAt(now, false))
	// fmt.Println("Valid ->", claims.VerifyExpiresAt(now, true))
	// fmt.Println("time.Unix ->", claims["exp"])

	// i := claims["exp"].(float64)
	// fmt.Println("expires at ->", time.Unix(int64(i), 0))
}
