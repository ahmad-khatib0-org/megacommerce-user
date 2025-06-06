package logger

import "encoding/json"

func toJSON(v any) any {
	if v == nil {
		return nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return v
	}

	var result any
	if err := json.Unmarshal(b, &result); err != nil {
		return v
	}
	return result
}
