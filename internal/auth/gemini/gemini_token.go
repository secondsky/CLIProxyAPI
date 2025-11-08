// Package gemini provides authentication and token management functionality
// for Google's Gemini AI services. It handles OAuth2 token storage, serialization,
// and retrieval for maintaining authenticated sessions with the Gemini API.
package gemini

import (
	"fmt"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/accounts"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/misc"
)

// GeminiTokenStorage stores OAuth2 token information for Google Gemini API authentication.
// It maintains compatibility with the existing auth system while adding Gemini-specific fields
// for managing access tokens, refresh tokens, and user account information.
type GeminiTokenStorage struct {
	// Token holds the raw OAuth2 token data, including access and refresh tokens.
	Token any `json:"token"`

	// ProjectID is the Google Cloud Project ID associated with this token.
	ProjectID string `json:"project_id"`

	// Email is the email address of the authenticated user.
	Email string `json:"email"`

	// Auto indicates if the project ID was automatically selected.
	Auto bool `json:"auto"`

	// Checked indicates if the associated Cloud AI API has been verified as enabled.
	Checked bool `json:"checked"`

	// Type indicates the authentication provider type, always "gemini" for this storage.
	Type string `json:"type"`
}

// SaveTokenToFile serializes the Gemini token storage to a JSON file.
// This method creates the necessary directory structure and writes the token
// data in JSON format to the specified file path for persistent storage.
//
// Parameters:
//   - authFilePath: The full path where the token file should be saved
//
// Returns:
//   - error: An error if the operation fails, nil otherwise
func (ts *GeminiTokenStorage) SaveTokenToFile(authFilePath string) error {
	if ts.Email == "" {
		return fmt.Errorf("email is required for Gemini account identification")
	}
	misc.LogSavingCredentials(authFilePath)
	accountID := ts.Email
	return accounts.SaveProviderAccount("gemini", accountID, func(existing map[string]any) map[string]any {
		for k, v := range map[string]any{
			"token":      ts.Token,
			"project_id": ts.ProjectID,
			"email":      ts.Email,
		} {
			if v == nil {
				continue
			}
			existing[k] = v
		}
		return existing
	})
}
