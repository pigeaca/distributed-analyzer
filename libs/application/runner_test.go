package application

import (
	"context"
	"errors"
	"testing"
	"time"
)

// MockComponent is a test implementation of the Component interface
type MockComponent struct {
	name          string
	startCalled   bool
	stopCalled    bool
	startError    error
	stopError     error
	startBlocking bool
	stopBlocking  bool
}

func NewMockComponent(name string) *MockComponent {
	return &MockComponent{
		name: name,
	}
}

func (m *MockComponent) Start(ctx context.Context) error {
	m.startCalled = true
	if m.startBlocking {
		<-ctx.Done()
	}
	return m.startError
}

func (m *MockComponent) Stop(ctx context.Context) error {
	m.stopCalled = true
	if m.stopBlocking {
		<-ctx.Done()
	}
	return m.stopError
}

func (m *MockComponent) Name() string {
	return m.name
}

func (m *MockComponent) WithStartError(err error) *MockComponent {
	m.startError = err
	return m
}

func (m *MockComponent) WithStopError(err error) *MockComponent {
	m.stopError = err
	return m
}

func (m *MockComponent) WithStartBlocking() *MockComponent {
	m.startBlocking = true
	return m
}

func (m *MockComponent) WithStopBlocking() *MockComponent {
	m.stopBlocking = true
	return m
}

func TestRunnerRegister(t *testing.T) {
	runner := NewApplicationRunner()
	comp1 := NewMockComponent("comp1")
	comp2 := NewMockComponent("comp2")

	runner.RegisterComponent(comp1)
	runner.RegisterComponent(comp2)

	if len(runner.components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(runner.components))
	}
}

func TestRunnerDefer(t *testing.T) {
	runner := NewApplicationRunner()
	called := false

	runner.Defer(func() error {
		called = true
		return nil
	})

	if len(runner.defers) != 1 {
		t.Errorf("Expected 1 defer function, got %d", len(runner.defers))
	}

	// Execute the defer function
	ctx := context.Background()
	runner.StopAll(ctx)

	if !called {
		t.Error("Defer function was not called")
	}
}

func TestRunnerStopAll(t *testing.T) {
	comp1 := NewMockComponent("comp1")
	comp2 := NewMockComponent("comp2")
	runner := NewApplicationRunner(comp1, comp2)

	ctx := context.Background()
	runner.StopAll(ctx)

	if !comp1.stopCalled {
		t.Error("Stop was not called on comp1")
	}

	if !comp2.stopCalled {
		t.Error("Stop was not called on comp2")
	}
}

func TestRunnerStopAllWithTimeout(t *testing.T) {
	comp1 := NewMockComponent("comp1")
	comp2 := NewMockComponent("comp2").WithStopBlocking()
	runner := NewApplicationRunner(comp1, comp2)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	runner.StopAll(ctx)

	if !comp1.stopCalled {
		t.Error("Stop was not called on comp1")
	}

	if !comp2.stopCalled {
		t.Error("Stop was not called on comp2")
	}
}

func TestRunnerStopAllWithErrors(t *testing.T) {
	expectedError := errors.New("stop error")
	comp1 := NewMockComponent("comp1")
	comp2 := NewMockComponent("comp2").WithStopError(expectedError)
	runner := NewApplicationRunner(comp1, comp2)

	ctx := context.Background()
	runner.StopAll(ctx)

	if !comp1.stopCalled {
		t.Error("Stop was not called on comp1")
	}

	if !comp2.stopCalled {
		t.Error("Stop was not called on comp2")
	}
}
