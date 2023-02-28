package meta

import (
	"encoding/json"
	"testing"
)

func TestJSON(t *testing.T) {
	text := `{}`

	var r Resource

	json.Unmarshal([]byte(text), &r)

	t.Log(r.Enable)
}

func GoBool(b bool) *bool {
	return &b
}
