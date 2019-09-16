package common

import (
	"bytes"
	"encoding/json"
	"errors"
)

type settingChangeEntry struct {
	Type *ChangeType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (entry *SettingChangeEntry) UnmarshalJSON(data []byte) error {
	var tmp settingChangeEntry
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return nil
	}
	if tmp.Type == nil {
		return errors.New("'type' field of change is required")
	}
	entry.Type = *tmp.Type
	obj, err := SettingChangeFromType(*tmp.Type)
	if err != nil {
		return err
	}
	decode := json.NewDecoder(bytes.NewBuffer(tmp.Data))
	decode.DisallowUnknownFields()

	if err = decode.Decode(obj); err != nil {
		if te, ok := err.(*json.UnmarshalTypeError); ok {
			return errors.New(te.Error())
		}
		return err
	}
	entry.Data = obj
	return nil
}
