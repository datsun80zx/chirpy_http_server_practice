package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	t.Run("successfully creates a valid JWT token", func(t *testing.T) {
		userID := uuid.New()
		secret := "test-secret-key"
		expiresIn := time.Hour

		token, err := MakeJWT(userID, secret, expiresIn)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if token == "" {
			t.Fatal("expected non-empty token")
		}
	})

	t.Run("token contains correct claims", func(t *testing.T) {
		userID := uuid.New()
		secret := "test-secret-key"
		expiresIn := time.Hour
		now := time.Now()

		token, err := MakeJWT(userID, secret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		// Parse the token to verify claims
		claims := &jwt.RegisteredClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			t.Fatalf("failed to parse token: %v", err)
		}

		if !parsedToken.Valid {
			t.Fatal("expected valid token")
		}

		// Verify issuer
		if claims.Issuer != "chirpy" {
			t.Errorf("expected issuer 'chirpy', got %s", claims.Issuer)
		}

		// Verify subject (user ID)
		if claims.Subject != userID.String() {
			t.Errorf("expected subject %s, got %s", userID.String(), claims.Subject)
		}

		// Verify IssuedAt is recent
		if claims.IssuedAt == nil {
			t.Fatal("expected IssuedAt to be set")
		}
		issuedAt := claims.IssuedAt.Time
		if issuedAt.Before(now.Add(-time.Minute)) || issuedAt.After(now.Add(time.Minute)) {
			t.Errorf("IssuedAt time %v is not within expected range of %v", issuedAt, now)
		}

		// Verify ExpiresAt
		if claims.ExpiresAt == nil {
			t.Fatal("expected ExpiresAt to be set")
		}
		expiresAt := claims.ExpiresAt.Time
		expectedExpiry := now.Add(expiresIn)
		// Allow 1 minute tolerance for timing
		if expiresAt.Before(expectedExpiry.Add(-time.Minute)) || expiresAt.After(expectedExpiry.Add(time.Minute)) {
			t.Errorf("ExpiresAt time %v is not within expected range of %v", expiresAt, expectedExpiry)
		}
	})

	t.Run("creates different tokens for different users", func(t *testing.T) {
		userID1 := uuid.New()
		userID2 := uuid.New()
		secret := "test-secret-key"
		expiresIn := time.Hour

		token1, err := MakeJWT(userID1, secret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token1: %v", err)
		}

		token2, err := MakeJWT(userID2, secret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token2: %v", err)
		}

		if token1 == token2 {
			t.Error("expected different tokens for different users")
		}
	})

	t.Run("creates different tokens with different secrets", func(t *testing.T) {
		userID := uuid.New()
		secret1 := "secret-one"
		secret2 := "secret-two"
		expiresIn := time.Hour

		token1, err := MakeJWT(userID, secret1, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token1: %v", err)
		}

		token2, err := MakeJWT(userID, secret2, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token2: %v", err)
		}

		if token1 == token2 {
			t.Error("expected different tokens for different secrets")
		}
	})

	t.Run("handles different expiration durations", func(t *testing.T) {
		userID := uuid.New()
		secret := "test-secret-key"

		testCases := []struct {
			name      string
			expiresIn time.Duration
		}{
			{"1 minute", time.Minute},
			{"1 hour", time.Hour},
			{"24 hours", 24 * time.Hour},
			{"1 second", time.Second},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				token, err := MakeJWT(userID, secret, tc.expiresIn)
				if err != nil {
					t.Fatalf("failed to create token: %v", err)
				}
				if token == "" {
					t.Fatal("expected non-empty token")
				}
			})
		}
	})

	t.Run("handles empty secret", func(t *testing.T) {
		userID := uuid.New()
		secret := ""
		expiresIn := time.Hour

		token, err := MakeJWT(userID, secret, expiresIn)
		// Should still create a token (though not secure)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token == "" {
			t.Fatal("expected non-empty token")
		}
	})

	t.Run("handles zero expiration", func(t *testing.T) {
		userID := uuid.New()
		secret := "test-secret-key"
		expiresIn := time.Duration(0)

		token, err := MakeJWT(userID, secret, expiresIn)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token == "" {
			t.Fatal("expected non-empty token")
		}
	})
}

func TestValidateJWT(t *testing.T) {
	t.Run("successfully validates a valid token", func(t *testing.T) {
		userID := uuid.New()
		secret := "test-secret-key"
		expiresIn := time.Hour

		token, err := MakeJWT(userID, secret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		validatedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if validatedID != userID {
			t.Errorf("expected user ID %s, got %s", userID, validatedID)
		}
	})

	t.Run("rejects token with wrong secret", func(t *testing.T) {
		userID := uuid.New()
		secret := "correct-secret"
		wrongSecret := "wrong-secret"
		expiresIn := time.Hour

		token, err := MakeJWT(userID, secret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		_, err = ValidateJWT(token, wrongSecret)
		if err == nil {
			t.Fatal("expected error for wrong secret, got nil")
		}
	})

	t.Run("rejects expired token", func(t *testing.T) {
		userID := uuid.New()
		secret := "test-secret-key"
		expiresIn := -time.Hour // Already expired

		token, err := MakeJWT(userID, secret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Fatal("expected error for expired token, got nil")
		}
	})

	t.Run("rejects malformed token", func(t *testing.T) {
		secret := "test-secret-key"

		testCases := []struct {
			name  string
			token string
		}{
			{"empty string", ""},
			{"random string", "not-a-jwt-token"},
			{"incomplete token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"},
			{"invalid format", "invalid.token.format.extra"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := ValidateJWT(tc.token, secret)
				if err == nil {
					t.Fatal("expected error for malformed token, got nil")
				}
			})
		}
	})

	t.Run("rejects token with invalid signing method", func(t *testing.T) {
		// Create a token with RS256 instead of HS256
		userID := uuid.New()
		secret := "test-secret-key"

		token := jwt.NewWithClaims(jwt.SigningMethodHS512, // Wrong signing method
			jwt.RegisteredClaims{
				Issuer:    "chirpy",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				Subject:   userID.String(),
			})

		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		_, err = ValidateJWT(tokenString, secret)
		if err == nil {
			t.Fatal("expected error for wrong signing method, got nil")
		}
	})

	t.Run("rejects token with invalid UUID in subject", func(t *testing.T) {
		secret := "test-secret-key"

		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.RegisteredClaims{
				Issuer:    "chirpy",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				Subject:   "not-a-valid-uuid",
			})

		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		_, err = ValidateJWT(tokenString, secret)
		if err == nil {
			t.Fatal("expected error for invalid UUID, got nil")
		}
	})

	t.Run("rejects token with empty subject", func(t *testing.T) {
		secret := "test-secret-key"

		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.RegisteredClaims{
				Issuer:    "chirpy",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				Subject:   "",
			})

		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		_, err = ValidateJWT(tokenString, secret)
		if err == nil {
			t.Fatal("expected error for empty subject, got nil")
		}
	})

	t.Run("validates token created just before expiration", func(t *testing.T) {
		userID := uuid.New()
		secret := "test-secret-key"
		expiresIn := 1000 * time.Millisecond

		token, err := MakeJWT(userID, secret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		// Validate immediately
		validatedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if validatedID != userID {
			t.Errorf("expected user ID %s, got %s", userID, validatedID)
		}
	})

	t.Run("returns uuid.Nil on error", func(t *testing.T) {
		secret := "test-secret-key"

		validatedID, err := ValidateJWT("invalid-token", secret)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if validatedID != uuid.Nil {
			t.Errorf("expected uuid.Nil, got %s", validatedID)
		}
	})
}

func TestMakeJWTAndValidateJWTIntegration(t *testing.T) {
	t.Run("round trip with valid token", func(t *testing.T) {
		userID := uuid.New()
		secret := "integration-test-secret"
		expiresIn := time.Hour

		// Create token
		token, err := MakeJWT(userID, secret, expiresIn)
		if err != nil {
			t.Fatalf("MakeJWT failed: %v", err)
		}

		// Validate token
		validatedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("ValidateJWT failed: %v", err)
		}

		// Verify the ID matches
		if validatedID != userID {
			t.Errorf("expected user ID %s, got %s", userID, validatedID)
		}
	})

	t.Run("multiple users round trip", func(t *testing.T) {
		secret := "integration-test-secret"
		expiresIn := time.Hour

		for i := 0; i < 10; i++ {
			userID := uuid.New()

			token, err := MakeJWT(userID, secret, expiresIn)
			if err != nil {
				t.Fatalf("MakeJWT failed for iteration %d: %v", i, err)
			}

			validatedID, err := ValidateJWT(token, secret)
			if err != nil {
				t.Fatalf("ValidateJWT failed for iteration %d: %v", i, err)
			}

			if validatedID != userID {
				t.Errorf("iteration %d: expected user ID %s, got %s", i, userID, validatedID)
			}
		}
	})
}
