package federation

import (
	"encoding/json"
	"strconv"
)

// NameResponse represents the result of a federation request
// for `name` and `forward` requests.
type NameResponse struct {
	AccountID string `json:"account_id"`
	MemoType  string `json:"memo_type,omitempty"`
	Memo      Memo   `json:"memo,omitempty"`
}

// IDResponse represents the result of a federation request
// for `id` request.
type IDResponse struct {
	Address string `json:"payshares_address"`
}

// Memo value can be either integer or string in JSON. This struct
// allows marshaling and unmarshaling both types.
type Memo struct {
	Value string
}

func (m Memo) MarshalJSON() ([]byte, error) {
	// Attempt to ParseUint. If this marshal m.Value
	// and return JSON string.
	_, err := strconv.ParseUint(m.Value, 10, 64)
	if err != nil {
		jsonString, err := json.Marshal(m.Value)
		if err != nil {
			return []byte{}, err
		}
		return jsonString, nil
	}
	return []byte(m.Value), nil
}

func (m *Memo) UnmarshalJSON(value []byte) error {
	// Try to unmarshal value into uint64. If that fails
	// unmarshal into string.
	var uintValue uint64
	err := json.Unmarshal(value, &uintValue)
	if err == nil {
		m.Value = strconv.FormatUint(uintValue, 10)
		return nil
	}
	err = json.Unmarshal(value, &m.Value)
	if err != nil {
		return err
	}
	return nil
}

func (m *Memo) String() string {
	return m.Value
}
