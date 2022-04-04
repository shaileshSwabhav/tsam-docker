package util

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/errors"
)

// IsEmpty Check str is Empty Or Not.
func IsEmpty(str string) bool {
	if len(strings.TrimSpace(str)) == 0 {
		return true
	}
	return false
}

// IsNil Check str is Empty Or Not (accept only pointer of string).
func IsNil(str *string) bool {
	if str == nil {
		return true
	}

	if len(strings.TrimSpace(*str)) == 0 {
		return true
	}
	return false
}

// RemoveTimeStampFromDate removes the timestamp from string.
func RemoveTimeStampFromDate(str string) *string {
	str = str[:10]
	//s := strings.TrimSuffix(str, "T00:00:00Z")
	return &str
}

// RemoveTimeStampFromTime removes the seconds from string.
func RemoveTimeStampFromTime(str string) *string {
	str = str[:16]
	//s := strings.TrimSuffix(str, ":00Z")
	return &str
}

// GenerateUniqueCode generates random code with first 3 alphabets of the name string
// & checks if the value is unique in DB based on condition.
// 	GenerateUniqueCode(uow.DB,"Ravi","`code` = ?",model.Talent{})
// will generate a code "RAV" + 7 digits, check the talent table and will return a unique code
// by doing necessary checks.
func GenerateUniqueCode(db *gorm.DB, name string, condition string, out interface{}) (string, error) {
	code := name
	for {
		if err := appendRandomText(&code); err != nil {
			return "", err
		}
		// exists, err := repository.DoesRecordExist(db, out, repository.Filter(condition, code))
		// if err != nil {
		// 	return "", err
		// }
		var count uint
		err := db.Model(out).Where(condition, code).Count(&count).Error
		if err != nil {
			return "", err
		}
		if count > 0 {
			if err := appendRandomText(&code); err != nil {
				return "", err
			}
			continue
		}
		break
	}
	return code, nil
}

// appendRandomText generates random text according to the given string
func appendRandomText(str *string) error {
	if str == nil {
		return errors.NewValidationError("Empty string value to convert to code.")
	}
	if IsEmpty(*str) {
		return errors.NewValidationError("Empty string value to convert to code.")
	}
	rand.Seed(time.Now().UTC().UnixNano())
	tempRandom := randInt(1000000, 10000000)
	var result strings.Builder
	var n = 0
	for i := 0; i < len(*str); i++ {
		b := (*str)[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') {
			n++
			err := result.WriteByte(b)
			if err != nil {
				return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
			}
			if n == 3 {
				*str = strings.ToUpper(fmt.Sprintf("%v%d", result.String(), tempRandom))
				return nil
			}
		}
	}
	if n == 2 {
		*str = strings.ToUpper(fmt.Sprintf("%v0%d", result.String(), tempRandom))
		return nil
	}
	if n == 1 {
		*str = strings.ToUpper(fmt.Sprintf("%v00%d", result.String(), tempRandom))
		return nil
	}
	*str = strings.ToUpper(fmt.Sprintf("%v000%d", result.String(), tempRandom))
	return nil
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// IsContain check weather base string contain pattern string and Return true otherwise false
func IsContain(base string, pattern string) bool {
	return strings.Contains(base, pattern)
}

// GeneratePassword returns random password string
func GeneratePassword() string {

	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "%/!@#$?"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + digits + specials
	length := 8

	password := make([]byte, length)
	password[0] = digits[rand.Intn(len(digits))]
	password[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		password[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}
