package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type PluginSettings struct {
	Region    string                `json:"region"`
	AccessKey string                `json:"accessKey,omitempty"`
	NumMaxRetries int               `json:"maxRetries,omitemtpy"`
	MinRetryDelay int               `json:"minRetryDelay,omitempty"`
	MinThrottleDelay int            `json:"minThrottleDelay,omitempty"`
	MaxRetryDelay int               `json:"maxRetryDelay,omitempty"`
	MaxThrottleDelay int            `json:"maxThrottleDelay,omitempty"`
	CacheSeconds int                `json:"cacheSeconds,omitemtpy"`
	Secrets   *SecretPluginSettings `json:"-"`
}

type SecretPluginSettings struct {
	AccessSecret string `json:"accessSecret,omitempty"`
	AccessToken string `json:"accessToken,omitempty"`
}

func LoadPluginSettings(source backend.DataSourceInstanceSettings) (*PluginSettings, error) {
	settings := PluginSettings{}
	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal PluginSettings json: %w", err)
	}

	settings.Secrets = loadSecretPluginSettings(source.DecryptedSecureJSONData)

	return &settings, nil
}

func loadSecretPluginSettings(source map[string]string) *SecretPluginSettings {
	return &SecretPluginSettings{
		AccessSecret: source["accessSecret"],
		AccessToken: source["accessToken"],
	}
}
