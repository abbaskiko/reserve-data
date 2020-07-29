// Code generated by "enumer -type=ChangeStatus -linecomment -json=true"; DO NOT EDIT.

//
package common

import (
	"encoding/json"
	"fmt"
)

const _ChangeStatusName = "pendingacceptedrejected"

var _ChangeStatusIndex = [...]uint8{0, 7, 15, 23}

func (i ChangeStatus) String() string {
	if i < 0 || i >= ChangeStatus(len(_ChangeStatusIndex)-1) {
		return fmt.Sprintf("ChangeStatus(%d)", i)
	}
	return _ChangeStatusName[_ChangeStatusIndex[i]:_ChangeStatusIndex[i+1]]
}

var _ChangeStatusValues = []ChangeStatus{0, 1, 2}

var _ChangeStatusNameToValueMap = map[string]ChangeStatus{
	_ChangeStatusName[0:7]:   0,
	_ChangeStatusName[7:15]:  1,
	_ChangeStatusName[15:23]: 2,
}

// ChangeStatusString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ChangeStatusString(s string) (ChangeStatus, error) {
	if val, ok := _ChangeStatusNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to ChangeStatus values", s)
}

// ChangeStatusValues returns all values of the enum
func ChangeStatusValues() []ChangeStatus {
	return _ChangeStatusValues
}

// IsAChangeStatus returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ChangeStatus) IsAChangeStatus() bool {
	for _, v := range _ChangeStatusValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for ChangeStatus
func (i ChangeStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for ChangeStatus
func (i *ChangeStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ChangeStatus should be a string, got %s", data)
	}

	var err error
	*i, err = ChangeStatusString(s)
	return err
}