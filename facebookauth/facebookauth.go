package facebookauth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// FacebookClaims represents the claims we expect in the Facebook JWT
type FacebookClaims struct {
	Iss     string      `json:"iss"`
	Sub     string      `json:"sub"`     // Facebook user ID
	Aud     interface{} `json:"aud"`     // Your app ID
	Exp     int64       `json:"exp"`     // Expiration time
	Iat     int64       `json:"iat"`     // Issued at time
	UserID  string      `json:"user_id"` // Same as sub
	Email   string      `json:"email"`   // User's email (if requested and granted)
	Name    string      `json:"name"`    // User's full name
	Picture string      `json:"picture"` // URL of profile picture
}

const (
	facebookJWKSUrl = "https://www.facebook.com/.well-known/oauth/openid/jwks/"
	facebookIssuer  = "https://www.facebook.com"
	maxClockSkew    = 10 * time.Minute
)

func VerifyFacebookJWT(tokenString, appID string) (*FacebookClaims, error) {
	// Fetch keys
	set, err := jwk.Fetch(context.Background(), facebookJWKSUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	// Parse with validation
	token, err := jwt.ParseString(
		tokenString,
		jwt.WithKeySet(set),
		jwt.WithValidate(true),
		jwt.WithIssuer(facebookIssuer),
		jwt.WithAudience(appID),
		jwt.WithAcceptableSkew(maxClockSkew),
	)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	// Additional manual expiration check with more context
	now := time.Now()
	if token.Expiration().Before(now.Add(-maxClockSkew)) {
		return nil, fmt.Errorf("token expired at %v (current time %v)",
			token.Expiration(), now)
	}

	// Extract claims
	claims := &FacebookClaims{}
	if err := json.Unmarshal(tokenToJSON(token), claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return claims, nil
}

// Helper function to convert jwt.Token to JSON bytes
func tokenToJSON(token jwt.Token) []byte {
	buf, _ := json.Marshal(token)
	return buf
}

// GetUserInfoFromToken is a convenience function that verifies the token and returns user info
func GetUserInfoFromToken(tokenString, appID string) (userID, email, name, picture string, err error) {
	claims, err := VerifyFacebookJWT(tokenString, appID)
	if err != nil {
		return "", "", "", "", err
	}
	return claims.Sub, claims.Email, claims.Name, claims.Picture, nil
}
