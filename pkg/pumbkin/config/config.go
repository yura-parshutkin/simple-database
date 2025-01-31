package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Logging struct {
	Level  string `mapstructure:"level"`
	Output string `mapstructure:"output"`
}

type Network struct {
	Address           string `mapstructure:"address"`
	MaxConnections    int    `mapstructure:"max_connections"`
	IdleTimeout       time.Duration
	MaxMessageSize    int
	MaxMessageSizeStr string `mapstructure:"max_message_size"`
	IdleTimeoutStr    string `mapstructure:"idle_timeout"`
}

type Engine struct {
	Type string `mapstructure:"type"`
}

type Config struct {
	Engine  Engine  `mapstructure:"engine"`
	Network Network `mapstructure:"network"`
	Logging Logging `mapstructure:"logging"`
}

func Parse(filename string) (*Config, error) {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	path := filepath.Dir(filename)
	baseName := filepath.Base(name)

	viper.SetConfigName(baseName)
	viper.SetConfigType(strings.TrimPrefix(ext, "."))
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}
	size, err := parseSizeString(config.Network.MaxMessageSizeStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing max_message_size: %v", err)
	}
	config.Network.MaxMessageSize = size
	duration, err := parseDurationString(config.Network.IdleTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing idle timeout: %v", err)
	}
	config.Network.IdleTimeout = duration
	return &config, nil
}

func parseDurationString(duration string) (time.Duration, error) {
	if duration == "" {
		return 0, nil
	}
	return time.ParseDuration(duration)
}

func parseSizeString(sizeString string) (int, error) {
	sizeString = strings.TrimSpace(sizeString)
	if sizeString == "" {
		return 0, nil
	}
	units := map[string]int{
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
		"PB": 1024 * 1024 * 1024 * 1024 * 1024,
		"B":  1,
	}
	for unit, multiplier := range units {
		if strings.HasSuffix(sizeString, unit) {
			numStr := strings.TrimSuffix(sizeString, unit)
			numStr = strings.TrimSpace(numStr)

			num, err := strconv.ParseInt(numStr, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number in size string: %w", err)
			}
			if num < 0 {
				return 0, fmt.Errorf("size cannot be negative")
			}
			return int(num) * multiplier, nil
		}
	}
	num, err := strconv.ParseInt(sizeString, 10, 64)
	if err == nil {
		if num < 0 {
			return 0, fmt.Errorf("size cannot be negative")
		}
		return int(num), nil
	}
	return 0, fmt.Errorf("invalid size format: %s", sizeString)
}
