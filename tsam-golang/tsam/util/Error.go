package util

import (
	"fmt"

	"github.com/techlabs/swabhav/tsam/errors"
)

// HandleError will handle the error given by DoesRecordExists by specifying the entity name in returned error
// 	handleError("Tenant",false,nil)
// will return Invalid Tenant ID
func HandleError(errMsg string, exists bool, err error) error {
	// exists will have zero value if error occurs.
	if !exists {
		if err != nil {
			return errors.NewValidationError(err.Error())
		}
		return errors.NewValidationError(errMsg)
	}
	return nil
}

// HandleIfExistsError will handle the error given by DoesRecordExists by specifying the entity name in returned error
// 	handleError("Tenant",false,nil)
// will return Invalid Tenant ID
func HandleIfExistsError(errMsg string, exists bool, err error) error {
	fmt.Println("HandleIfExistsError======================", errMsg, exists)
	if err != nil {
		return errors.NewValidationError(err.Error())
	}
	if exists {
		return errors.NewValidationError(errMsg)
	}
	return nil
}
