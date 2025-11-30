// Package gemini provides authentication and token management functionality
// for Google's Gemini AI services. It handles OAuth2 token storage, serialization,
// and retrieval for maintaining authenticated sessions with the Gemini API.
package gemini

import (
	"fmt"
	"strings"

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
	misc.LogSavingCredentials(authFilePath)
	ts.Type = "gemini"

	if ts.Email == "" {
		return fmt.Errorf("email is required for Gemini account identification")
	}

	// Distinguish accounts by project when provided (supports multi-project installs).
	accountID := strings.TrimSpace(ts.Email)
	if proj := strings.TrimSpace(ts.ProjectID); proj != "" {
		norm := proj
		if strings.EqualFold(norm, "all") || strings.Contains(norm, ",") {
			norm = "all"
		}
		accountID = fmt.Sprintf("%s-%s", accountID, norm)
	}

	return accounts.SaveProviderAccount("gemini", accountID, func(existing map[string]any) map[string]any {
		for k, v := range map[string]any{
			"token":      ts.Token,
			"project_id": ts.ProjectID,
			"email":      ts.Email,
			"auto":       ts.Auto,
			"checked":    ts.Checked,
		} {
			switch val := v.(type) {
			case string:
				if val == "" {
					continue
				}
			case nil:
				continue
			}
			existing[k] = v
		}
		return existing
	})
}

// CredentialFileName returns the filename used to persist Gemini CLI credentials.
// When projectID represents multiple projects (comma-separated or literal ALL),
// the suffix is normalized to "all" and a "gemini-" prefix is enforced to keep
// web and CLI generated files consistent.
func CredentialFileName(email, projectID string, includeProviderPrefix bool) string {
	email = strings.TrimSpace(email)
	project := strings.TrimSpace(projectID)
	if strings.EqualFold(project, "all") || strings.Contains(project, ",") {
		return fmt.Sprintf("gemini-%s-all.json", email)
	}
	prefix := ""
	if includeProviderPrefix {
		prefix = "gemini-"
	}
	return fmt.Sprintf("%s%s-%s.json", prefix, email, project)
}
