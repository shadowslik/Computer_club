package config

import "os"

type Config struct {
	Port string

	// 2.2.1 — запросы (-1 = отключён, 0 = запрещено всё)
	DefaultRPS int
	DefaultRPM int
	DefaultRPH int
	DefaultRPD int

	// 2.2.2 — трафик
	DefaultMaxUploadBps   int64
	DefaultMaxDownloadBps int64
	DefaultMaxTrafficDay  int64

	// 2.2.3 — соединения
	DefaultMaxConcurrent   int
	DefaultMaxNewPerSecond int
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		Port: port,

		DefaultRPS: 100, // 100 запросов в секунду
		DefaultRPM: -1,  // отключён
		DefaultRPH: -1,
		DefaultRPD: -1,

		DefaultMaxUploadBps:   -1, // отключён
		DefaultMaxDownloadBps: -1,
		DefaultMaxTrafficDay:  -1,

		DefaultMaxConcurrent:   100, // 100 одновременных соединений
		DefaultMaxNewPerSecond: -1,
	}
}
