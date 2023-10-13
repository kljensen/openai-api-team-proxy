package main

import (
	"bytes"
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3"
)

// Function takes an http.Request

func payload_is_allowed(payload string, logic string) (error, bool) {
	var result bytes.Buffer
	payloadReader := bytes.NewReader([]byte(payload))
	logicReader := bytes.NewReader([]byte(logic))
	err := jsonlogic.Apply(logicReader, payloadReader, &result)
	return err, strings.TrimSpace(result.String()) == "true"
}
