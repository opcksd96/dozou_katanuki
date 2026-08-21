package main

import (
	"strings"
	"testing"
)

func TestUpdateStashConfigYAML(t *testing.T) {
	input := `blobs_path: blobs
host: 127.0.0.1
port: 9999
dangerous_allow_public_without_auth: "false"
database: stash-go.sqlite
`
	updated, changed := updateStashConfigYAML(input, "0.0.0.0", 9999)
	if !changed {
		t.Fatalf("expected changed to be true")
	}

	if !strings.Contains(updated, "host: 0.0.0.0") {
		t.Errorf("expected host: 0.0.0.0, got: %s", updated)
	}
	if !strings.Contains(updated, `dangerous_allow_public_without_auth: "true"`) {
		t.Errorf("expected dangerous_allow_public_without_auth: \"true\", got: %s", updated)
	}
	if !strings.Contains(updated, "blobs_path: blobs") || !strings.Contains(updated, "database: stash-go.sqlite") {
		t.Errorf("other fields should be preserved, got: %s", updated)
	}

	// 2回目の実行では変更なしになることを検証
	_, changed2 := updateStashConfigYAML(updated, "0.0.0.0", 9999)
	if changed2 {
		t.Errorf("expected no changes on second run")
	}
}
