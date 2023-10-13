package main

import (
	"testing"
)

// "model" must be "gpt-3.5-turbo" or "gpt-3.5-turbo-instruct"
var canOnlyUseTurbo = `{
	"or": [
		{"==": [{"var": "model"}, "gpt-3.5-turbo"]},
		{"==": [{"var": "model"}, "gpt-3.5-turbo-instruct"]}
	]
}`

func TestPermissions(t *testing.T) {
	payload1 := `{
		"model": "gpt-3.5-turbo-instruct",
		"prompt": "Say this is a test",
		"max_tokens": 7,
		"temperature": 0
}`
	err, allowed := payload_is_allowed(payload1, canOnlyUseTurbo)
	if err != nil {
		t.Fatalf("Should not be an error")
	}
	if !allowed {
		t.Fatalf("payload should be allowed")
	}
	payload2 := `{
		"model": "gpt-4",
		"prompt": "Say this is a test",
		"max_tokens": 7,
		"temperature": 0
}`
	err, allowed = payload_is_allowed(payload2, canOnlyUseTurbo)
	if err != nil {
		t.Fatalf("Should not be an error")
	}
	if allowed {
		t.Fatalf("payload should not be allowed")
	}

}
