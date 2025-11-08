package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	DefaultDir         = ".cli-proxy-api"
	ActiveAccountsFile = "active-accounts.json"
)

type Account struct {
	Type            string         `json:"type"`
	AccountID       string         `json:"accountId"`
	AccountNickname string         `json:"accountNickname,omitempty"`
	Email           string         `json:"email,omitempty"`
	Expired         *time.Time     `json:"expired,omitempty"`
	CreatedAt       *time.Time     `json:"createdAt,omitempty"`
	Provider        string         `json:"-"`
	FileBase        string         `json:"-"`
	Raw             map[string]any `json:"-"`
}

type ProviderAccounts map[string][]*Account

type ActiveAccounts map[string]string

func LoadAllAccounts() (ProviderAccounts, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, DefaultDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ProviderAccounts{}, nil
		}
		return nil, err
	}

	result := ProviderAccounts{}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		if name == ActiveAccountsFile {
			continue
		}

		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil || len(data) == 0 {
			continue
		}

		var raw map[string]any
		if err := json.Unmarshal(data, &raw); err != nil {
			continue
		}

		base := strings.TrimSuffix(name, ".json")
		provider, accountID := inferProviderAndAccountID(base, raw)
		if provider == "" {
			continue
		}

		acc := &Account{
			Type:      provider,
			Provider:  provider,
			AccountID: accountID,
			FileBase:  base,
			Raw:       raw,
		}

		if v, ok := raw["accountNickname"].(string); ok {
			acc.AccountNickname = v
		}
		if v, ok := raw["email"].(string); ok {
			acc.Email = v
		}
		if v, ok := raw["expired"].(string); ok {
			if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
				acc.Expired = &t
			}
		}
		if v, ok := raw["createdAt"].(string); ok {
			if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
				acc.CreatedAt = &t
			}
		}

		if acc.AccountID == "" {
			acc.AccountID = base
		}

		result[provider] = append(result[provider], acc)
	}

	return result, nil
}

func inferProviderAndAccountID(base string, raw map[string]any) (string, string) {
	if t, ok := raw["type"].(string); ok && t != "" {
		if id, ok := raw["accountId"].(string); ok && id != "" {
			return strings.ToLower(t), id
		}
	}

	parts := strings.SplitN(base, "-", 2)
	if len(parts) == 2 {
		provider := strings.ToLower(parts[0])
		id := parts[1]
		if provider != "" && id != "" {
			return provider, id
		}
	}

	if t, ok := raw["type"].(string); ok && t != "" {
		return strings.ToLower(t), base
	}

	return "", ""
}

func LoadActiveAccounts() ActiveAccounts {
	home, err := os.UserHomeDir()
	if err != nil {
		return ActiveAccounts{}
	}
	path := filepath.Join(home, DefaultDir, ActiveAccountsFile)
	b, err := os.ReadFile(path)
	if err != nil || len(b) == 0 {
		return ActiveAccounts{}
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return ActiveAccounts{}
	}
	return ActiveAccounts(m)
}

func ResolveActiveAccount(provider string, accounts []*Account, active ActiveAccounts) *Account {
	if len(accounts) == 0 {
		return nil
	}

	provider = strings.ToLower(provider)
	var hint string
	if active != nil {
		if v, ok := active[provider]; ok {
			hint = v
		}
	}

	pick := func(filter func(*Account) bool) *Account {
		for _, a := range accounts {
			if a == nil || !strings.EqualFold(a.Provider, provider) {
				continue
			}
			if filter(a) {
				return a
			}
		}
		return nil
	}

	if hint != "" {
		if acc := pick(func(a *Account) bool { return !IsExpired(a) && a.AccountID == hint }); acc != nil {
			return acc
		}

		if strings.HasPrefix(hint, provider+"-") {
			trim := strings.TrimPrefix(hint, provider+"-")
			if acc := pick(func(a *Account) bool { return !IsExpired(a) && a.AccountID == trim }); acc != nil {
				return acc
			}
		}

		if acc := pick(func(a *Account) bool { return !IsExpired(a) && a.Email != "" && strings.EqualFold(a.Email, hint) }); acc != nil {
			return acc
		}

		if acc := pick(func(a *Account) bool {
			if !IsExpired(a) && a.FileBase == hint {
				return true
			}
			if !IsExpired(a) && strings.HasPrefix(a.FileBase, provider+"-") && strings.TrimPrefix(a.FileBase, provider+"-") == hint {
				return true
			}
			return false
		}); acc != nil {
			return acc
		}
	}

	if acc := pick(func(a *Account) bool { return !IsExpired(a) }); acc != nil {
		return acc
	}

	return accounts[0]
}

func IsExpired(a *Account) bool {
	if a == nil || a.Expired == nil {
		return false
	}
	return a.Expired.Before(time.Now())
}

func SaveProviderAccount(provider, accountID string, update func(existing map[string]any) map[string]any) error {
	if provider == "" || accountID == "" {
		return fmt.Errorf("provider and accountID required")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, DefaultDir)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	filename := fmt.Sprintf("%s-%s.json", strings.ToLower(provider), accountID)
	path := filepath.Join(dir, filename)

	var existing map[string]any
	if b, err := os.ReadFile(path); err == nil && len(b) > 0 {
		_ = json.Unmarshal(b, &existing)
	}
	if existing == nil {
		existing = map[string]any{}
	}

	existing["type"] = strings.ToLower(provider)
	existing["accountId"] = accountID

	if update != nil {
		existing = update(existing)
		if existing == nil {
			return fmt.Errorf("update func returned nil")
		}
	}

	b, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0o600)
}
