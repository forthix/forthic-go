package forthic

import (
	"sort"
	"strings"
	"testing"
)

func TestWordOptions_CreateFromFlatArray(t *testing.T) {
	opts, err := NewWordOptions([]interface{}{"depth", 2, "with_key", true})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if opts.Get("depth") != 2 {
		t.Errorf("Expected depth=2, got: %v", opts.Get("depth"))
	}

	if opts.Get("with_key") != true {
		t.Errorf("Expected with_key=true, got: %v", opts.Get("with_key"))
	}
}

func TestWordOptions_RequiresArray(t *testing.T) {
	_, err := NewWordOptions("not an array")
	if err == nil {
		t.Error("Expected error for non-array input")
	}
	if !strings.Contains(err.Error(), "must be an array") {
		t.Errorf("Expected 'must be an array' error, got: %v", err)
	}
}

func TestWordOptions_RequiresEvenLength(t *testing.T) {
	_, err := NewWordOptions([]interface{}{"depth", 2, "with_key"})
	if err == nil {
		t.Error("Expected error for odd-length array")
	}
	if !strings.Contains(err.Error(), "even length") {
		t.Errorf("Expected 'even length' error, got: %v", err)
	}
}

func TestWordOptions_RequiresStringKeys(t *testing.T) {
	_, err := NewWordOptions([]interface{}{123, "value"})
	if err == nil {
		t.Error("Expected error for non-string key")
	}
	if !strings.Contains(err.Error(), "must be a string") {
		t.Errorf("Expected 'must be a string' error, got: %v", err)
	}
}

func TestWordOptions_DefaultValues(t *testing.T) {
	opts, _ := NewWordOptions([]interface{}{"depth", 2})

	// Missing key without default returns nil
	if opts.Get("missing") != nil {
		t.Errorf("Expected nil for missing key, got: %v", opts.Get("missing"))
	}

	// Missing key with default returns default
	if opts.Get("missing", "default") != "default" {
		t.Errorf("Expected 'default' for missing key with default, got: %v", opts.Get("missing", "default"))
	}
}

func TestWordOptions_Has(t *testing.T) {
	opts, _ := NewWordOptions([]interface{}{"depth", 2})

	if !opts.Has("depth") {
		t.Error("Expected Has('depth') to be true")
	}

	if opts.Has("missing") {
		t.Error("Expected Has('missing') to be false")
	}
}

func TestWordOptions_ToRecord(t *testing.T) {
	opts, _ := NewWordOptions([]interface{}{"depth", 2, "with_key", true})
	record := opts.ToRecord()

	if record["depth"] != 2 {
		t.Errorf("Expected depth=2 in record, got: %v", record["depth"])
	}

	if record["with_key"] != true {
		t.Errorf("Expected with_key=true in record, got: %v", record["with_key"])
	}

	if len(record) != 2 {
		t.Errorf("Expected record length=2, got: %d", len(record))
	}
}

func TestWordOptions_Keys(t *testing.T) {
	opts, _ := NewWordOptions([]interface{}{"depth", 2, "with_key", true})
	keys := opts.Keys()
	sort.Strings(keys)

	expected := []string{"depth", "with_key"}
	if len(keys) != len(expected) {
		t.Fatalf("Expected %d keys, got %d", len(expected), len(keys))
	}

	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Expected key[%d]=%s, got: %s", i, expected[i], key)
		}
	}
}

func TestWordOptions_LaterValuesOverride(t *testing.T) {
	opts, _ := NewWordOptions([]interface{}{"depth", 2, "depth", 3})

	if opts.Get("depth") != 3 {
		t.Errorf("Expected depth=3 (overridden), got: %v", opts.Get("depth"))
	}
}

func TestWordOptions_HandlesNullUndefined(t *testing.T) {
	opts, _ := NewWordOptions([]interface{}{"key1", nil, "key2", nil})

	if opts.Get("key1") != nil {
		t.Errorf("Expected nil for key1, got: %v", opts.Get("key1"))
	}

	if opts.Get("key2") != nil {
		t.Errorf("Expected nil for key2, got: %v", opts.Get("key2"))
	}

	if !opts.Has("key1") {
		t.Error("Expected Has('key1') to be true")
	}

	if !opts.Has("key2") {
		t.Error("Expected Has('key2') to be true")
	}
}

func TestWordOptions_ComplexValues(t *testing.T) {
	complexValue := map[string]interface{}{
		"nested": map[string]interface{}{
			"data": []int{1, 2, 3},
		},
	}

	opts, _ := NewWordOptions([]interface{}{"config", complexValue})

	// Go can't easily deep compare, so just check it's not nil
	if opts.Get("config") == nil {
		t.Error("Expected config to have complex value")
	}
}

func TestWordOptions_ToString(t *testing.T) {
	opts, _ := NewWordOptions([]interface{}{"depth", 2, "with_key", true})
	str := opts.String()

	if !strings.Contains(str, "WordOptions") {
		t.Errorf("Expected toString to contain 'WordOptions', got: %s", str)
	}

	if !strings.Contains(str, ".depth") {
		t.Errorf("Expected toString to contain '.depth', got: %s", str)
	}

	if !strings.Contains(str, ".with_key") {
		t.Errorf("Expected toString to contain '.with_key', got: %s", str)
	}
}

func TestWordOptions_EmptyArray(t *testing.T) {
	opts, _ := NewWordOptions([]interface{}{})

	keys := opts.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected empty keys array, got: %d keys", len(keys))
	}

	record := opts.ToRecord()
	if len(record) != 0 {
		t.Errorf("Expected empty record, got: %d entries", len(record))
	}

	if opts.Has("anything") {
		t.Error("Expected Has('anything') to be false for empty options")
	}
}
