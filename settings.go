package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// ISettings interface
type ISettings interface {
	GetInt(key string, fallback int) int
	GetInt64(key string, fallback int64) int64
	Get(key string, fallback string) string
	GetBool(key string, fallback bool) bool
	GetSection(key string) ISettings
}

type settings struct {
	configFileCash map[string]interface{}
	cash           map[string]interface{}
}

func (s *settings) GetSection(key string) ISettings {
	if s.configFileCash != nil {
		if value, ok := s.configFileCash[key]; ok {
			var settings = &settings{cash: map[string]interface{}{}}
			if valueMap, ok := value.(map[string]interface{}); ok {
				settings.configFileCash = valueMap
				return settings
			}

			if valueMap, ok := value.(map[interface{}]interface{}); ok {
				settings.configFileCash = make(map[string]interface{})
				for subKey, subValue := range valueMap {
					settings.configFileCash[subKey.(string)] = subValue
				}

				return settings
			}
		}
	}

	return s
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

func loadConfigFile(configFileName string) map[string]interface{} {
	if configFileName != "" {
		file, err := os.Open(configFileName)
		if err == nil {
			defer file.Close()
			byteValue, err := ioutil.ReadAll(file)
			if err == nil {
				var result map[string]interface{}
				if strings.HasSuffix(configFileName, ".json") {
					err = json.Unmarshal(byteValue, &result)

				} else if strings.HasSuffix(configFileName, "yaml") {
					err = yaml.Unmarshal(byteValue, &result)
				}

				if err == nil {
					return result
				}
			}
		}
	}

	return nil
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

// Get the settings object
func Get(configFile string) ISettings {
	var settings = &settings{cash: map[string]interface{}{}}
	settings.configFileCash = loadConfigFile(configFile)
	return settings
}
