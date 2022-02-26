package validator

import (
	"errors"
	"fmt"

	"github.com/shuntagami/dojo1/kadai1/shuntagami/helper"
)

func ValidateInput(from, to string) error {
	if from == to {
		return errors.New("from and to cannot be same")
	}
	if from != helper.JPG && from != helper.PNG {
		return fmt.Errorf("arg FROM must be %s or %s", helper.JPG, helper.PNG)
	}
	if to != helper.JPG && to != helper.PNG {
		return fmt.Errorf("arg TO must be %s or %s", helper.JPG, helper.PNG)
	}
	return nil
}
