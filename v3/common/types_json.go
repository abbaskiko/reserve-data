package common

import "encoding/json"

type settingChangeEntry struct {
	Type ChangeType      `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (entry *SettingChangeEntry) UnmarshalJSON(data []byte) error {
	var tmp settingChangeEntry
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return nil
	}
	entry.Type = tmp.Type
	obj, err := SettingChangeFromType(tmp.Type)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(tmp.Data, obj); err != nil {
		return err
	}
	entry.Data = obj
	return nil
}
