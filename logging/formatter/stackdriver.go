package formatter

import (
	"encoding/json"
)

type timestampField struct {
	Seconds json.Number `json:"seconds"`
	Nanos   json.Number `json:"nanos"`
}

// StackdriverFormatter does stuff
type StackdriverFormatter struct {
	Severity string                 `json:"severity"`
	Time     json.Number            `json:"time"`
	Message  map[string]interface{} `json:"message,omitempty"`
}

// Format creates a
// func (f *StackdriverFormatter) Format(e *logrus.Entry) ([]byte, error) {
// 	var b bytes.Buffer

// 	// Fill in the respective fields

// 	enc := json.NewEncoder(&b)
// 	if err := enc.Encode(f); err != nil {
// 		return nil, fmt.Errorf("failed to marshal log entry: %v", err)
// 	}

// 	return b.Bytes(), nil
// }

// func timestamp(duration time.Duration) *timestampField {
// 	return &timestampField{
// 		Seconds: "",
// 		Nanos:   "",
// 	}
// }
