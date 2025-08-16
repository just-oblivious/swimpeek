package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/charmbracelet/huh"
	"github.com/google/uuid"
)

// initPrompts shows the user prompts for configuration initialization.
func initPrompts(cfg *Config) error {
	confForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&cfg.SwimlaneRegion).
				Title("Swimlane Region (e.g. us1, de1)").
				Placeholder("us1").
				Description("To find the region, visit any Swimlane tenant and look at the domain name.").
				Validate(func(s string) error {
					if !regexp.MustCompile(`^[a-z]{2}\d$`).MatchString(s) {
						return errors.New("please provide a valid Swimlane region (e.g., us1, de1)")
					}
					return nil
				}),
			huh.NewInput().
				Value(&cfg.SwimlaneAccountId).
				Title("Swimlane Account UUID").
				Placeholder("123e4567-e89b-12d3-a456-426614174000").
				Description("To find your account UUID, visit any Swimlane tenant and copy the value after '/account/' from the URL.").
				Validate(func(s string) error {
					if err := uuid.Validate(s); err != nil {
						return errors.New("please provide a valid UUID for the Swimlane Account ID")
					}
					return nil
				}),
			huh.NewInput().
				Value(&cfg.SwimlaneAccessToken).
				Title("Swimlane Access Token").
				EchoMode(huh.EchoModePassword).
				Placeholder("your-access-token").
				Description("To create a token, visit any Swimlane tenant and go to 'Profile & user settings' → 'Personal access token'.").
				Validate(func(s string) error {
					if len(s) != 64 {
						return errors.New("please provide a valid Swimlane Access Token (exactly 64 characters)")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeDracula()).WithLayout(huh.LayoutStack)

	return confForm.Run()
}

// InitConfig initializes the SwimPeek configuration.
func InitConfig(cfgDir string) (*Config, error) {
	// Load the existing configuration or create a fresh one
	cfg, err := ReadConfig(cfgDir)
	if errors.Is(err, os.ErrNotExist) {
		cfg = &Config{}
	} else if err != nil {
		return nil, fmt.Errorf("failed to read existing configuration; delete file to start over: %w", err)
	}

	// Show the initialization prompts
	if err := initPrompts(cfg); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, fmt.Errorf("setup was aborted by the user")
		}
		return nil, fmt.Errorf("failed to show configuration form: %w", err)
	}

	return cfg, nil
}
