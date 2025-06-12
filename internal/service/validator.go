package service

import (
	"encoding/json"
)

type Validator struct {
}

func (v *Validator) validateJson(data *json.RawMessage) error {
	return nil
}

func (v *Validator) validateScript(script string) error {

}
