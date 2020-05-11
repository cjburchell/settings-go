package settings

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// ISettings interface
type ISettings interface {
	GetInt(key string, fallback int) int
	GetInt64(key string, fallback int64) int64
	Get(key string, fallback string) string
	GetBool(key string, fallback bool) bool
}

type settings struct {
	configFile string
	configFileCash map[string]interface{}
	cash map[string]interface{}
}

func (s *settings) Get(key string, fallback string) string {
	value := s.get(key, fallback)
	if val, ok := value.(string); ok {
		return val
	}

	return fallback
}

func (s *settings) GetInt(key string, fallback int) int {
	value := s.get(key, fallback)
	if val, ok := value.(int); ok {
		return val
	} else if val, ok := value.(float64); ok {
		return int(val)
	} else if val, ok := value.(string); ok {
		if result, err := strconv.Atoi(val); err == nil {
			return result
		}
	}
	return fallback
}

func (s *settings) GetInt64(key string, fallback int64) int64 {
	value := s.get(key, fallback)
	if val, ok := value.(int64); ok {
		return val
	} else if val, ok := value.(float64); ok {
		return int64(val)
	} else if val, ok := value.(int); ok {
		return int64(val)
	} else if val, ok := value.(string); ok {
		if result, err := strconv.ParseInt(val, 10, 64); err == nil {
			return result
		}
	}
	return fallback
}

func (s *settings) get(key string, fallback interface{}) interface{} {
	if _, ok := s.cash[key]; ok {
		return s.cash[key]
	}

	if s.configFileCash == nil {
		if s.configFile != "" {
			file, err := os.Open(s.configFile)
			defer file.Close()
			if err == nil {
				byteValue, err := ioutil.ReadAll(file)
				if err == nil {
					var result map[string]interface{}
					if strings.HasSuffix(s.configFile, ".json") {
						err = json.Unmarshal(byteValue, &result)

					} else if strings.HasSuffix(s.configFile, "yaml") {
						err = yaml.Unmarshal(byteValue, &result)
					}

					if err == nil {
						s.configFileCash = result
					}
				}
			}
		}
	}

	if s.configFileCash != nil {
		if value, ok := s.configFileCash[key]; ok {
			s.cash[key] = value
			return s.cash[key]
		}
	}

	if value, ok := os.LookupEnv(key); ok {
		s.cash[key] = value
		return s.cash[key]
	}

	return fallback

}

func (s *settings) GetBool(key string, fallback bool) bool {
	value := s.get(key, fallback)
	if val, ok := value.(bool); ok {
		return val
	} else if val, ok := value.(string); ok {
		if result, err := strconv.ParseBool(val); err == nil {
			return result
		}
	}
	return fallback
}

// Gets the settings object
func Get(configFile string) ISettings {
	var settings = &settings{ cash: map[string]interface{}{}, configFile: configFile}
	return settings
}

