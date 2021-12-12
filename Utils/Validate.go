package Utils

import (
	"fmt"
	"koinfolio/Models"
)

func ValidateResponse(status Models.Status) bool {
	if status.ErrorMessage != nil {
		fmt.Println(status.ErrorMessage)
		return false
	}

	return true
}
