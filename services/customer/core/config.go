package core

// core packages
import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"strconv"

	"github.com/golang/groupcache"
	"github.com/joho/godotenv"
)

// third-party packages

// NewConfig comment
func NewConfig(appDir string) (*Config, error) {
	var err error

	cfg := new(Config)

	cfg.GroupCache = groupcache.NewGroup("global_config", 64<<20, groupcache.GetterFunc(
		func(ctx context.Context, key string, dest groupcache.Sink) error {
			data, err := cfg.Load(key)
			if err != nil {
				return errors.New("config not found")
			}
			dest.SetBytes(data)
			return nil
		},
	))

	if appDir == "" {
		appDir, err = GetAppDir()
		if err != nil {
			return cfg, err
		}
	}
	cfg.AppDir = appDir

	envFile := path.Join(appDir, ".env")
	if _, err := os.Stat(envFile); err == nil {
		if err = godotenv.Load(envFile); err != nil {
			return cfg, err
		}
	}

	if err := cfg.Init(os.Getenv("GO_ENV")); err != nil {
		return cfg, err
	}

	port, err := strconv.ParseInt(cfg.Port, 10, 32)
	hostAddress := BuildString("http://", cfg.ClusterHostName, ":", string(port+1))
	peers := groupcache.NewHTTPPool(hostAddress)
	peers.Set(hostAddress)

	return cfg, err
}

// // GlobalGroupCache comment
// var GlobalGroupCache

// Init comment
func (cfg *Config) Init(scopeName string) (err error) {
	var data []byte

	if scopeName == "" {
		scopeName = "development"
	}

	err = cfg.GroupCache.Get(nil, scopeName, groupcache.AllocatingByteSliceSink(&data))

	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cfg)
}

// Load comment
func (cfg *Config) Load(key string) (data []byte, err error) {
	appDir := cfg.AppDir
	if err != nil {
		return data, err
	}

	cfgPath := path.Join(
		appDir,
		"config",
		BuildString("config.", key, ".json"),
	)
	configFile, err := os.Open(cfgPath)
	if err != nil {
		return data, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return json.Marshal(cfg)
}
