package mcp

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidAPIKey = errors.New("invalid or missing API key")
	ErrUnauthorized  = errors.New("API key does not have access to this location")
	ErrExpiredKey    = errors.New("API key has expired")
	ErrInactiveKey   = errors.New("API key is inactive")
)

// Context key for API key
type contextKey string

const apiKeyContextKey contextKey = "mcp_api_key"

// AuthHandler handles API key authentication for MCP
type AuthHandler struct{}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// ExtractAPIKeyFromHeader extracts the API key from Authorization header
// Supports: "Bearer <key>" or just "<key>"
func ExtractAPIKeyFromHeader(authHeader string) string {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return ""
	}

	// Check for Bearer prefix
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return authHeader
}

// ContextWithAPIKey returns a context with the API key embedded
func ContextWithAPIKey(ctx context.Context, apiKey string) context.Context {
	return context.WithValue(ctx, apiKeyContextKey, apiKey)
}

// APIKeyFromContext retrieves the API key from context
func APIKeyFromContext(ctx context.Context) string {
	if key, ok := ctx.Value(apiKeyContextKey).(string); ok {
		return key
	}
	return ""
}

// ValidateAndGetKey validates the API key from context and returns the key record
func (a *AuthHandler) ValidateAndGetKey(ctx context.Context) (*models.MCPApiKey, error) {
	apiKey := APIKeyFromContext(ctx)
	if apiKey == "" {
		return nil, ErrInvalidAPIKey
	}

	keyHash := models.HashMCPApiKey(apiKey)

	var mcpKey models.MCPApiKey
	err := db.DB.
		Where("key_hash = ?", keyHash).
		Preload("AllowedLocations").
		Preload("Profile").
		First(&mcpKey).Error

	if err != nil {
		return nil, ErrInvalidAPIKey
	}

	if !mcpKey.IsActive {
		return nil, ErrInactiveKey
	}

	if mcpKey.ExpiresAt != nil && mcpKey.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredKey
	}

	// Update last used timestamp asynchronously
	go func() {
		now := time.Now()
		db.DB.Model(&mcpKey).Update("last_used_at", now)
	}()

	return &mcpKey, nil
}

// CanAccessLocation checks if the API key can access the given location
func (a *AuthHandler) CanAccessLocation(key *models.MCPApiKey, locationID string) bool {
	return key.CanAccessLocation(locationID)
}
