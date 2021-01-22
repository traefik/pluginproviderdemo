package dynamic

import "encoding/json"

type JSONPayload struct {
	*Configuration
}

func (c JSONPayload) MarshalJSON() ([]byte, error) {
	if c.Configuration == nil {
		return nil, nil
	}

	return json.Marshal(c.Configuration)
}
