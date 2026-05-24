package session

import (
	"testing"
)

type mockStore struct {
	data map[string]map[string]any
}

func (s *mockStore) Read(id string) (map[string]any, error) {
	if data, ok := s.data[id]; ok {
		return data, nil
	}
	return nil, nil
}

func (s *mockStore) Write(id string, data map[string]any) error {
	s.data[id] = data
	return nil
}

func (s *mockStore) Destroy(id string) error {
	delete(s.data, id)
	return nil
}

func TestSessionManager(t *testing.T) {
	store := &mockStore{data: make(map[string]map[string]any)}
	
	// Request 1: Start and Flash
	manager1 := NewManager(store, "test-session")
	manager1.Start()
	
	manager1.Put("user", "Alice")
	manager1.Flash("message", "Welcome!")
	manager1.FlashInput(map[string]any{"email": "test@test.com"})
	
	if manager1.Get("user") != "Alice" {
		t.Errorf("expected Alice")
	}
	
	manager1.Save()

	// Request 2: Read flashed data
	manager2 := NewManager(store, "test-session")
	manager2.Start()
	
	if manager2.Get("user") != "Alice" {
		t.Errorf("expected Put data to persist")
	}
	
	if manager2.Get("message") != "Welcome!" {
		t.Errorf("expected flash message to be available, got %v", manager2.Get("message"))
	}
	
	if manager2.Old("email", "") != "test@test.com" {
		t.Errorf("expected old input to be available")
	}

	// Reflash the message
	manager2.Keep("message")
	manager2.Save()

	// Request 3: Check Keep
	manager3 := NewManager(store, "test-session")
	manager3.Start()

	if manager3.Get("message") != "Welcome!" {
		t.Errorf("expected kept flash message to be available")
	}

	if manager3.Old("email", "") == "test@test.com" {
		t.Errorf("did not expect old input to be available (was not kept)")
	}

	manager3.Save()

	// Request 4: Ensure flash data is destroyed
	manager4 := NewManager(store, "test-session")
	manager4.Start()

	if manager4.Get("message") != nil {
		t.Errorf("expected flash message to be destroyed")
	}

	// Test Regenerate
	oldID := manager4.ID()
	manager4.Regenerate()
	
	if manager4.ID() == oldID {
		t.Errorf("expected new session ID after regenerate")
	}
	
	if _, exists := store.data[oldID]; exists {
		t.Errorf("expected old session data to be destroyed in store")
	}
	
	if _, exists := store.data[manager4.ID()]; !exists {
		t.Errorf("expected new session to be saved in store")
	}
}

