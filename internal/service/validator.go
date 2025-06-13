package service

import (
	"encoding/json"
)

type Validator struct {
}

func (v *Validator) ValidateScript(script string) any {
	panic("unimplemented")
}

func (v *Validator) validateJson(data *json.RawMessage) error {
	panic("unimplemented")
}

func (v *Validator) validateScript(script string) error {
	panic("unimplemented")
}
