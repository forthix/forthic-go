package forthic

import (
	"errors"
	"testing"
)

func TestWordErrorHandler_AddHandler(t *testing.T) {
	word := NewBaseWord("TEST")
	handler := func(err error, word Word, interp *Interpreter) error {
		return nil
	}

	if len(word.GetErrorHandlers()) != 0 {
		t.Errorf("Expected 0 handlers, got %d", len(word.GetErrorHandlers()))
	}

	word.AddErrorHandler(handler)

	if len(word.GetErrorHandlers()) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(word.GetErrorHandlers()))
	}
}

func TestWordErrorHandler_ClearHandlers(t *testing.T) {
	word := NewBaseWord("TEST")
	word.AddErrorHandler(func(err error, word Word, interp *Interpreter) error {
		return nil
	})
	word.AddErrorHandler(func(err error, word Word, interp *Interpreter) error {
		return nil
	})

	if len(word.GetErrorHandlers()) != 2 {
		t.Errorf("Expected 2 handlers, got %d", len(word.GetErrorHandlers()))
	}

	word.ClearErrorHandlers()

	if len(word.GetErrorHandlers()) != 0 {
		t.Errorf("Expected 0 handlers after clear, got %d", len(word.GetErrorHandlers()))
	}
}

func TestWordErrorHandler_GetHandlersReturnsCopy(t *testing.T) {
	word := NewBaseWord("TEST")
	handler := func(err error, word Word, interp *Interpreter) error {
		return nil
	}
	word.AddErrorHandler(handler)

	handlers := word.GetErrorHandlers()
	// Try to modify the returned slice
	handlers = append(handlers, func(err error, word Word, interp *Interpreter) error {
		return nil
	})

	// Original should be unchanged
	if len(word.GetErrorHandlers()) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(word.GetErrorHandlers()))
	}
}

func TestWordErrorHandler_ModuleWordCallsHandler(t *testing.T) {
	interp := NewInterpreter()
	handlerCalled := false
	var receivedError error

	word := NewModuleWord("FAILING-WORD", func(interp *Interpreter) error {
		return errors.New("Test error")
	})

	word.AddErrorHandler(func(err error, w Word, i *Interpreter) error {
		handlerCalled = true
		receivedError = err
		return nil
	})

	err := word.Execute(interp)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}

	if receivedError == nil || receivedError.Error() != "Test error" {
		t.Errorf("Expected error message 'Test error', got %v", receivedError)
	}
}

func TestWordErrorHandler_ErrorSuppressedWhenHandlerSucceeds(t *testing.T) {
	interp := NewInterpreter()

	word := NewModuleWord("FAILING-WORD", func(interp *Interpreter) error {
		return errors.New("Test error")
	})

	word.AddErrorHandler(func(err error, w Word, i *Interpreter) error {
		// Handler succeeds by returning nil
		return nil
	})

	err := word.Execute(interp)

	if err != nil {
		t.Errorf("Expected error to be suppressed, got %v", err)
	}
}

func TestWordErrorHandler_ErrorPropagatesWhenHandlerThrows(t *testing.T) {
	interp := NewInterpreter()

	word := NewModuleWord("FAILING-WORD", func(interp *Interpreter) error {
		return errors.New("Original error")
	})

	word.AddErrorHandler(func(err error, w Word, i *Interpreter) error {
		return errors.New("Handler also failed")
	})

	err := word.Execute(interp)

	if err == nil {
		t.Error("Expected error to propagate")
	}

	// Should get original error, not handler error
	if err.Error() != "Original error" {
		t.Errorf("Expected 'Original error', got %v", err)
	}
}

func TestWordErrorHandler_HandlersCalledInOrder(t *testing.T) {
	interp := NewInterpreter()
	callOrder := []int{}

	word := NewModuleWord("FAILING-WORD", func(interp *Interpreter) error {
		return errors.New("Test error")
	})

	word.AddErrorHandler(func(err error, w Word, i *Interpreter) error {
		callOrder = append(callOrder, 1)
		return errors.New("Handler 1 failed")
	})

	word.AddErrorHandler(func(err error, w Word, i *Interpreter) error {
		callOrder = append(callOrder, 2)
		// Handler 2 succeeds
		return nil
	})

	word.AddErrorHandler(func(err error, w Word, i *Interpreter) error {
		callOrder = append(callOrder, 3)
		// Should not be called
		return nil
	})

	_ = word.Execute(interp)

	if len(callOrder) != 2 {
		t.Errorf("Expected 2 handlers called, got %d", len(callOrder))
	}

	if callOrder[0] != 1 || callOrder[1] != 2 {
		t.Errorf("Expected call order [1, 2], got %v", callOrder)
	}
}

func TestWordErrorHandler_IntentionalStopErrorBypassesHandlers(t *testing.T) {
	interp := NewInterpreter()
	handlerCalled := false

	word := NewModuleWord("STOPPING-WORD", func(interp *Interpreter) error {
		return NewIntentionalStopError("Intentional stop")
	})

	word.AddErrorHandler(func(err error, w Word, i *Interpreter) error {
		handlerCalled = true
		return nil
	})

	err := word.Execute(interp)

	if err == nil {
		t.Error("Expected IntentionalStopError to propagate")
	}

	if _, ok := err.(*IntentionalStopError); !ok {
		t.Errorf("Expected IntentionalStopError, got %T", err)
	}

	if handlerCalled {
		t.Error("Handler should not be called for IntentionalStopError")
	}
}

func TestWordErrorHandler_AddModuleWordCreatesWordWithHandlerSupport(t *testing.T) {
	interp := NewInterpreter()
	module := NewModule("test-module", "")
	handlerCalled := false

	module.AddModuleWord("FAILING", func(interp *Interpreter) error {
		return errors.New("Test error")
	})

	word := module.FindWord("FAILING")
	if word == nil {
		t.Fatal("Expected to find FAILING word")
	}

	if _, ok := word.(*ModuleWord); !ok {
		t.Errorf("Expected ModuleWord, got %T", word)
	}

	word.AddErrorHandler(func(err error, w Word, i *Interpreter) error {
		handlerCalled = true
		return nil
	})

	_ = word.Execute(interp)

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
}

func TestWordErrorHandler_HandlerReceivesWordAndInterpreter(t *testing.T) {
	interp := NewInterpreter()
	var receivedWord Word
	var receivedInterp *Interpreter

	word := NewModuleWord("TEST-WORD", func(interp *Interpreter) error {
		return errors.New("Test error")
	})

	word.AddErrorHandler(func(err error, w Word, i *Interpreter) error {
		receivedWord = w
		receivedInterp = i
		return nil
	})

	_ = word.Execute(interp)

	if receivedWord != word {
		t.Error("Handler did not receive correct word")
	}

	if receivedInterp != interp {
		t.Error("Handler did not receive correct interpreter")
	}
}
