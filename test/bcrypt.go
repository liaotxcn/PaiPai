package main

import (
	"PaiPai/pkg/encrypt"
	"testing"
)

// Bcrypt加密测试

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "paipai123456", false},
		{"empty password", "", false},
		{"long password", "thisIsAVeryLongPasswordWithMoreThanSeventyTwoBytesLengthWhichIsTheBcryptLimitButItShouldStillWorkBecauseBcryptWillHashItAnyway", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encrypt.HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("hashPassword() returned empty string")
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "correctPassword"
	wrongPassword := "wrongPassword"

	hashed, err := encrypt.HashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword failed: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hashed   string
		want     bool
	}{
		{"correct password", password, hashed, true},
		{"wrong password", wrongPassword, hashed, false},
		{"empty password", "", hashed, false},
		{"empty hash", password, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encrypt.CheckPassword(tt.password, tt.hashed); got != tt.want {
				t.Errorf("checkPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkPassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := encrypt.HashPassword(password)
		if err != nil {
			b.Fatalf("hashPassword failed: %v", err)
		}
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "benchmarkPassword123"
	hashed, err := encrypt.HashPassword(password)
	if err != nil {
		b.Fatalf("hashPassword failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encrypt.CheckPassword(password, hashed)
	}
}
