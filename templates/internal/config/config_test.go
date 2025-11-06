package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/a1ostudio/nova/internal/config"
	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	return p
}

func TestLoadConfig_FromFile(t *testing.T) {
	dir := t.TempDir()
	// prepare a minimal app.env
	content := `TOKEN_SYMMETRIC_KEY=secretkey
ACCESS_TOKEN_DURATION=30m
REFRESH_TOKEN_DURATION=24h
LOCK_TTL=3s
MAX_WAIT_TIME=1500ms
INITIAL_WAIT_TIME=100ms
MAX_SINGLE_WAIT=300ms
`
	writeFile(t, dir, "app.env", content)

	cfg, err := config.LoadConfig(dir)
	require.NoError(t, err)
	require.Equal(t, "secretkey", cfg.TokenSymmetricKey)
	require.Equal(t, 30*time.Minute, cfg.AccessTokenDuration)
	require.Equal(t, 24*time.Hour, cfg.RefreshTokenDuration)
	require.Equal(t, 3*time.Second, cfg.LockTTL)
	require.Equal(t, 1500*time.Millisecond, cfg.MaxWaitTime)
	require.Equal(t, 100*time.Millisecond, cfg.InitialWaitTime)
	require.Equal(t, 300*time.Millisecond, cfg.MaxSingleWait)
}

func TestLoadConfig_EnvOverridesFile(t *testing.T) {
	dir := t.TempDir()
	content := `TOKEN_SYMMETRIC_KEY=filekey
ACCESS_TOKEN_DURATION=15m
`
	writeFile(t, dir, "app.env", content)

	// env should override file
	os.Setenv("TOKEN_SYMMETRIC_KEY", "envkey")
	defer os.Unsetenv("TOKEN_SYMMETRIC_KEY")

	cfg, err := config.LoadConfig(dir)
	require.NoError(t, err)
	require.Equal(t, "envkey", cfg.TokenSymmetricKey)
	require.Equal(t, 15*time.Minute, cfg.AccessTokenDuration)
}

func TestLoadConfig_InvalidDuration_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	content := `ACCESS_TOKEN_DURATION=notaduration
`
	writeFile(t, dir, "app.env", content)

	_, err := config.LoadConfig(dir)
	require.Error(t, err)
}

func TestLoadConfig_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		env         map[string]string
		wantErr     bool
		wantToken   string
		wantAccess  time.Duration
		wantLockTTL time.Duration
	}{
		{
			name: "valid_file",
			fileContent: `TOKEN_SYMMETRIC_KEY=secretkey
ACCESS_TOKEN_DURATION=30m
LOCK_TTL=3s
`,
			env:         nil,
			wantErr:     false,
			wantToken:   "secretkey",
			wantAccess:  30 * time.Minute,
			wantLockTTL: 3 * time.Second,
		},
		{
			name: "env_override",
			fileContent: `TOKEN_SYMMETRIC_KEY=filekey
ACCESS_TOKEN_DURATION=15m
`,
			env: map[string]string{
				"TOKEN_SYMMETRIC_KEY": "envkey",
			},
			wantErr:    false,
			wantToken:  "envkey",
			wantAccess: 15 * time.Minute,
		},
		{
			name: "invalid_duration",
			fileContent: `ACCESS_TOKEN_DURATION=bad
`,
			env:     nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, dir, "app.env", tc.fileContent)

			// set envs
			for k, v := range tc.env {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			cfg, err := config.LoadConfig(dir)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tc.wantToken != "" {
				require.Equal(t, tc.wantToken, cfg.TokenSymmetricKey)
			}
			if tc.wantAccess != 0 {
				require.Equal(t, tc.wantAccess, cfg.AccessTokenDuration)
			}
			if tc.wantLockTTL != 0 {
				require.Equal(t, tc.wantLockTTL, cfg.LockTTL)
			}
		})
	}
}
