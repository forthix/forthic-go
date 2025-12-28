package modules

import (
	"testing"

	"github.com/forthix/forthic-go/forthic"
)

func setupStringInterpreter() *forthic.Interpreter {
	interp := forthic.NewInterpreter()
	strMod := NewStringModule()
	interp.ImportModule(strMod.Module, "")
	return interp
}

func TestString_ToStr(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run("42 >STR")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "42" {
		t.Errorf("Expected '42', got %v", result)
	}
}

func TestString_ConcatTwoStrings(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"Hello" " World" CONCAT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "Hello World" {
		t.Errorf("Expected 'Hello World', got %v", result)
	}
}

func TestString_ConcatArray(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`["Hello" " " "World"] CONCAT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "Hello World" {
		t.Errorf("Expected 'Hello World', got %v", result)
	}
}

func TestString_Split(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"a,b,c" "," SPLIT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	arr := result.([]interface{})
	if len(arr) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(arr))
	}
	if arr[0].(string) != "a" || arr[1].(string) != "b" || arr[2].(string) != "c" {
		t.Errorf("Expected [a b c], got %v", arr)
	}
}

func TestString_Join(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`["a" "b" "c"] "," JOIN`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "a,b,c" {
		t.Errorf("Expected 'a,b,c', got %v", result)
	}
}

func TestString_SlashN(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run("/N")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "\n" {
		t.Errorf("Expected newline, got %v", result)
	}
}

func TestString_SlashR(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run("/R")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "\r" {
		t.Errorf("Expected carriage return, got %v", result)
	}
}

func TestString_SlashT(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run("/T")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "\t" {
		t.Errorf("Expected tab, got %v", result)
	}
}

func TestString_Lowercase(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"HELLO" LOWERCASE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "hello" {
		t.Errorf("Expected 'hello', got %v", result)
	}
}

func TestString_Uppercase(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"hello" UPPERCASE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "HELLO" {
		t.Errorf("Expected 'HELLO', got %v", result)
	}
}

func TestString_ASCII(t *testing.T) {
	interp := setupStringInterpreter()
	// \u0100 is character 256, should be filtered out
	interp.StackPush("Hello\u0100World")
	err := interp.Run("ASCII")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "HelloWorld" {
		t.Errorf("Expected 'HelloWorld', got %v", result)
	}
}

func TestString_Strip(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"  hello  " STRIP`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "hello" {
		t.Errorf("Expected 'hello', got %v", result)
	}
}

func TestString_Replace(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"hello world" "world" "there" REPLACE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "hello there" {
		t.Errorf("Expected 'hello there', got %v", result)
	}
}

func TestString_ReMatchSuccess(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"test123" "test[0-9]+" RE-MATCH`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()

	// Should be an array
	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected array, got %T", result)
	}

	if len(arr) == 0 {
		t.Error("Expected non-empty match")
	}

	if arr[0].(string) != "test123" {
		t.Errorf("Expected 'test123', got %v", arr[0])
	}
}

func TestString_ReMatchFailure(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"test" "[0-9]+" RE-MATCH`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result != false {
		t.Errorf("Expected false for no match, got %v", result)
	}
}

func TestString_ReMatchAll(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"test1 test2 test3" "test([0-9])" RE-MATCH-ALL`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	arr := result.([]interface{})

	if len(arr) != 3 {
		t.Fatalf("Expected 3 matches, got %d", len(arr))
	}

	if arr[0].(string) != "1" || arr[1].(string) != "2" || arr[2].(string) != "3" {
		t.Errorf("Expected ['1', '2', '3'], got %v", arr)
	}
}

func TestString_ReMatchGroup(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"test123" "test([0-9]+)" RE-MATCH 1 RE-MATCH-GROUP`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "123" {
		t.Errorf("Expected '123', got %v", result)
	}
}

func TestString_URLEncode(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"hello world" URL-ENCODE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "hello+world" {
		t.Errorf("Expected 'hello+world', got %v", result)
	}
}

func TestString_URLDecode(t *testing.T) {
	interp := setupStringInterpreter()
	err := interp.Run(`"hello+world" URL-DECODE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result.(string) != "hello world" {
		t.Errorf("Expected 'hello world', got %v", result)
	}
}
