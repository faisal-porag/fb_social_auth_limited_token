package main

import (
	"fmt"
	"github.com/faisal-porag/fb_social_auth_limited_token/facebookauth"
	"log"
)

func main() {
	// This is the JWT token you received from the client (Limited Login)
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6..............."

	// Your Facebook App ID
	appID := "346212......."

	// Verify the token
	claims, err := facebookauth.VerifyFacebookJWT(token, appID)
	if err != nil {
		log.Fatalf("Token verification failed: %v", err)
	}

	fmt.Printf("Authenticated user: %s\n", claims.Sub)
	fmt.Printf("Email: %s\n", claims.Email)
	fmt.Printf("Name: %s\n", claims.Name)

	// Or use the convenience function
	userID, email, name, picture, err := facebookauth.GetUserInfoFromToken(token, appID)
	if err != nil {
		log.Fatalf("Failed to get user info: %v", err)
	}

	fmt.Printf("User ID: %s, Email: %s, Name: %s, pic:%s\n", userID, email, name, picture)
}
