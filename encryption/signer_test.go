package encryption

import (
	"testing"
	"time"
)

func TestNewSigner(t *testing.T) {
	signer, err := NewSigner("test-secret-key")
	if err != nil {
		t.Fatalf("NewSigner failed: %v", err)
	}
	if signer == nil {
		t.Fatal("expected non-nil signer")
	}
}

func TestNewSigner_EmptyKey(t *testing.T) {
	_, err := NewSigner("")
	if err == nil {
		t.Error("expected error for empty key")
	}
}

func TestSigner_SignVerify(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	message := "hello world"
	signed := signer.Sign(message)

	if signed == "" {
		t.Error("expected non-empty signature")
	}

	if !signer.Verify(message, signed) {
		t.Error("expected Verify to return true for valid signature")
	}
}

func TestSigner_VerifyInvalid(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	if signer.Verify("hello", "invalid-signature") {
		t.Error("expected Verify to return false for invalid signature")
	}
}

func TestSigner_VerifyTampered(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	message := "hello world"
	signed := signer.Sign(message)

	// Tamper with the message
	if signer.Verify("hello world!", signed) {
		t.Error("expected Verify to return false for tampered message")
	}
}

func TestSigner_DifferentKeys(t *testing.T) {
	signer1, _ := NewSigner("key-1")
	signer2, _ := NewSigner("key-2")

	message := "hello world"
	signed := signer1.Sign(message)

	// Verify with different key should fail
	if signer2.Verify(message, signed) {
		t.Error("expected Verify to return false with different key")
	}
}

func TestSigner_SignHMAC(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	message := "hello world"
	signed := signer.SignHMAC(message)

	if signed == "" {
		t.Error("expected non-empty HMAC signature")
	}

	if !signer.VerifyHMAC(message, signed) {
		t.Error("expected VerifyHMAC to return true")
	}
}

func TestSigner_VerifyHMACInvalid(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	message := "hello world"
	signed := signer.SignHMAC(message)

	if signer.VerifyHMAC("hello world!", signed) {
		t.Error("expected VerifyHMAC to return false for tampered message")
	}
}

func TestGenerateRandomKey(t *testing.T) {
	key1 := GenerateRandomKey(32)
	key2 := GenerateRandomKey(32)

	if key1 == "" {
		t.Error("expected non-empty key")
	}
	if key1 == key2 {
		t.Error("expected two random keys to be different")
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	bytes1, err := GenerateRandomBytes(32)
	if err != nil {
		t.Fatalf("GenerateRandomBytes failed: %v", err)
	}

	if len(bytes1) != 32 {
		t.Errorf("expected 32 bytes, got %d", len(bytes1))
	}

	bytes2, _ := GenerateRandomBytes(32)
	if string(bytes1) == string(bytes2) {
		t.Error("expected two random byte slices to be different")
	}
}

func TestConstantTimeEqual(t *testing.T) {
	a := []byte("hello")
	b := []byte("hello")
	c := []byte("world")

	if !ConstantTimeEqual(a, b) {
		t.Error("expected ConstantTimeEqual to return true for equal slices")
	}

	if ConstantTimeEqual(a, c) {
		t.Error("expected ConstantTimeEqual to return false for different slices")
	}
}

func TestSigner_SignWithTimestamp(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	message := "hello world"
	timestamp := time.Now().Unix()
	signed := signer.SignWithTimestamp(message, timestamp)

	if signed == nil {
		t.Fatal("expected non-nil SignedMessage")
	}

	if signed.Payload != message {
		t.Errorf("expected Payload %q, got %q", message, signed.Payload)
	}

	if signed.Timestamp != timestamp {
		t.Errorf("expected Timestamp %d, got %d", timestamp, signed.Timestamp)
	}
}

func TestSigner_VerifySignedMessage(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	message := "hello world"
	timestamp := time.Now().Unix()
	signed := signer.SignWithTimestamp(message, timestamp)

	if !signer.VerifySignedMessage(signed, 3600) {
		t.Error("expected VerifySignedMessage to return true")
	}
}

func TestSigner_VerifySignedMessage_Nil(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	if signer.VerifySignedMessage(nil, 3600) {
		t.Error("expected VerifySignedMessage to return false for nil")
	}
}

func TestSigner_SignURL(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	url := "https://example.com/download?file=test.pdf"
	signed := signer.SignURL(url)

	if signed == nil {
		t.Fatal("expected non-nil SignedURL")
	}

	if signed.URL != url {
		t.Errorf("expected URL %q, got %q", url, signed.URL)
	}

	if signed.Signature == "" {
		t.Error("expected non-empty signature")
	}
}

func TestSigner_VerifyURL(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	url := "https://example.com/download?file=test.pdf"
	signed := signer.SignURL(url)

	if !signer.VerifyURL(signed.URL, signed.Signature) {
		t.Error("expected VerifyURL to return true")
	}

	if signer.VerifyURL(url, "invalid-signature") {
		t.Error("expected VerifyURL to return false for invalid signature")
	}
}

func TestSigner_Hash(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	message := "hello world"
	hashed := signer.Hash(message)

	if hashed == "" {
		t.Error("expected non-empty hash")
	}

	// Same input should produce same hash
	hashed2 := signer.Hash(message)
	if hashed != hashed2 {
		t.Error("expected same hash for same input")
	}
}

func TestSigner_HashHMAC(t *testing.T) {
	signer, _ := NewSigner("test-secret-key")

	message := "hello world"
	hashed := signer.HashHMAC(message)

	if hashed == "" {
		t.Error("expected non-empty HMAC hash")
	}

	// Same input should produce same hash
	hashed2 := signer.HashHMAC(message)
	if hashed != hashed2 {
		t.Error("expected same HMAC hash for same input")
	}
}
