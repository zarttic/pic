package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: "this-is-a-very-long-password-that-should-still-work-fine",
			wantErr:  false,
		},
		{
			name:     "password with special characters",
			password: "p@$$w0rd!#$%",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}
				if hash == tt.password {
					t.Error("HashPassword() returned plaintext password")
				}
				if len(hash) < 20 {
					t.Errorf("HashPassword() returned suspiciously short hash: %s", hash)
				}
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	// 预先生成一些测试用的密码哈希
	validPassword := "password123"
	validHash, err := HashPassword(validPassword)
	if err != nil {
		t.Fatalf("Failed to generate test hash: %v", err)
	}

	emptyPassword := ""
	emptyHash, err := HashPassword(emptyPassword)
	if err != nil {
		t.Fatalf("Failed to generate empty hash: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "correct password",
			password: validPassword,
			hash:     validHash,
			expected: true,
		},
		{
			name:     "incorrect password",
			password: "wrongpassword",
			hash:     validHash,
			expected: false,
		},
		{
			name:     "empty password with empty hash",
			password: emptyPassword,
			hash:     emptyHash,
			expected: true,
		},
		{
			name:     "empty password with non-empty hash",
			password: emptyPassword,
			hash:     validHash,
			expected: false,
		},
		{
			name:     "case sensitive password",
			password: "Password123",
			hash:     validHash,
			expected: false,
		},
		{
			name:     "invalid hash format",
			password: validPassword,
			hash:     "invalid-hash-string",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPassword(tt.password, tt.hash)
			if result != tt.expected {
				t.Errorf("CheckPassword() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestPasswordHashUniqueness(t *testing.T) {
	password := "samepassword123"

	// 同一个密码应该生成不同的哈希(bcrypt 自动加盐)
	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to generate first hash: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to generate second hash: %v", err)
	}

	if hash1 == hash2 {
		t.Error("HashPassword() generated identical hashes for the same password")
	}

	// 但是两个哈希都应该能验证原密码
	if !CheckPassword(password, hash1) {
		t.Error("CheckPassword() failed for hash1")
	}
	if !CheckPassword(password, hash2) {
		t.Error("CheckPassword() failed for hash2")
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkPassword123"
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "benchmarkPassword123"
	hash, _ := HashPassword(password)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CheckPassword(password, hash)
	}
}
