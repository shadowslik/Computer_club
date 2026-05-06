package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port              string        `yaml:"port"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	CORSAllowedOrigin string        `yaml:"cors_allowed_origin"`
}

type ListsConfig struct {
	WhitelistFile string        `yaml:"whitelist_file"`
	BlacklistFile string        `yaml:"blacklist_file"`
	GraylistFile  string        `yaml:"graylist_file"`
	RateLimitFile string        `yaml:"ratelimit_file"`
	DefaultDeny   bool          `yaml:"default_deny"`
	CacheSize     int           `yaml:"cache_size"`
	CacheTTL      time.Duration `yaml:"cache_ttl"`
}

type RateLimitConfig struct {
	DefaultRPS            int           `yaml:"default_rps"`
	DefaultRPM            int           `yaml:"default_rpm"`
	DefaultRPH            int           `yaml:"default_rph"`
	DefaultRPD            int           `yaml:"default_rpd"`
	DefaultMaxUploadBps   int64         `yaml:"default_max_upload_bps"`
	DefaultMaxDownloadBps int64         `yaml:"default_max_download_bps"`
	DefaultMaxTrafficDay  int64         `yaml:"default_max_traffic_day"`
	DefaultMaxConcurrent  int           `yaml:"default_max_concurrent"`
	DefaultMaxNewPerSec   int           `yaml:"default_max_new_per_second"`
	BanDuration           time.Duration `yaml:"ban_duration"`
}

type UpstreamConfig struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Type string `yaml:"type"` // "http" or "tcp"
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	Server    ServerConfig     `yaml:"server"`
	Lists     ListsConfig      `yaml:"lists"`
	RateLimit RateLimitConfig  `yaml:"rate_limit"`
	Upstreams []UpstreamConfig `yaml:"upstreams"`
	Log       LogConfig        `yaml:"log"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	applyDefaults(cfg)
	return cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8081"
	}
	if cfg.Server.IdleTimeout == 0 {
		cfg.Server.IdleTimeout = 60 * time.Second
	}
	if cfg.Server.CORSAllowedOrigin == "" {
		cfg.Server.CORSAllowedOrigin = "http://localhost:3000"
	}
	if cfg.Lists.CacheSize == 0 {
		cfg.Lists.CacheSize = 1000
	}
	if cfg.Lists.CacheTTL == 0 {
		cfg.Lists.CacheTTL = 30 * time.Second
	}
	if cfg.RateLimit.BanDuration == 0 {
		cfg.RateLimit.BanDuration = 5 * time.Minute
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
}
