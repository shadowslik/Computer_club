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
	ProxyTarget       string        `yaml:"proxy_target"`
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
	Type string `yaml:"type"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type CacheRuleConfig struct {
	Match   string        `yaml:"match"`
	TTL     time.Duration `yaml:"ttl"`
	Tags    []string      `yaml:"tags"`
	Enabled *bool         `yaml:"enabled"` // nil = inherit global
}

type CacheCascadeRuleConfig struct {
	TriggerTag  string   `yaml:"trigger_tag"`
	CascadeTags []string `yaml:"cascade_tags"`
}

type CacheConfig struct {
	Enabled      bool                     `yaml:"enabled"`
	DefaultTTL   time.Duration            `yaml:"default_ttl"`
	MaxSize      int                      `yaml:"max_size"`
	MaxBodySize  int64                    `yaml:"max_body_size"`
	MinBodySize  int64                    `yaml:"min_body_size"`
	Cache2xx     bool                     `yaml:"cache_2xx"`
	Cache3xx     bool                     `yaml:"cache_3xx"`
	TTL3xx       time.Duration            `yaml:"ttl_3xx"`
	Cache4xx     bool                     `yaml:"cache_4xx"`
	TTL4xx       time.Duration            `yaml:"ttl_4xx"`
	Cache5xx     bool                     `yaml:"cache_5xx"`
	TTL5xx       time.Duration            `yaml:"ttl_5xx"`
	Rules        []CacheRuleConfig        `yaml:"rules"`
	CascadeRules []CacheCascadeRuleConfig `yaml:"cascade_rules"`
}

type AuthConfig struct {
	AdminKey   string   `yaml:"admin_key"`
	AdminPaths []string `yaml:"admin_paths"`
}

type Config struct {
	Server    ServerConfig     `yaml:"server"`
	Lists     ListsConfig      `yaml:"lists"`
	RateLimit RateLimitConfig  `yaml:"rate_limit"`
	Upstreams []UpstreamConfig `yaml:"upstreams"`
	Cache     CacheConfig      `yaml:"cache"`
	Auth      AuthConfig       `yaml:"auth"`
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
	if cfg.Cache.DefaultTTL == 0 {
		cfg.Cache.DefaultTTL = 5 * time.Minute
	}
	if cfg.Cache.MaxSize == 0 {
		cfg.Cache.MaxSize = 10000
	}
	if cfg.Cache.MaxBodySize == 0 {
		cfg.Cache.MaxBodySize = 1 << 20 // 1 MB
	}
	if cfg.Cache.TTL3xx == 0 {
		cfg.Cache.TTL3xx = time.Minute
	}
	if cfg.Cache.TTL4xx == 0 {
		cfg.Cache.TTL4xx = 30 * time.Second
	}
	if cfg.Cache.TTL5xx == 0 {
		cfg.Cache.TTL5xx = 10 * time.Second
	}
	if len(cfg.Auth.AdminPaths) == 0 {
		cfg.Auth.AdminPaths = []string{
			"/ip/",
			"/ratelimit",
			"/cache",
			"/metrics",
			"/top-clients",
			"/swagger/",
		}
	}
}
