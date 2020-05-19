package settings

import (
	"encoding/json"
	"fmt"
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
	GetObject(key string, obj interface{}) error
}

type fileType string

const (
	fileTypeJSON fileType = "json"
	fileTypeYAML fileType = "yaml"
)

type configFile struct {
	cash     map[string]interface{}
	fileType fileType
}

type settings struct {
	section    string
	configFile *configFile
	cash       map[string]interface{}
}

func getMapInterface(value interface{}) map[string]interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		return valueMap
	} else if valueInterface, ok := value.(map[interface{}]interface{}); ok {
		valueMap := make(map[string]interface{})
		for subKey, subValue := range valueInterface {
			valueMap[subKey.(string)] = subValue
		}
		return valueMap
	}

	return nil
}

func (s *settings) GetObject(key string, obj interface{}) error {
	if s.configFile != nil {
		if value, ok := s.configFile.cash[key]; ok {
			if s.configFile.fileType == fileTypeJSON {
				jsonBody, err := json.Marshal(value)
				if err != nil {
					return err
				}

				return json.Unmarshal(jsonBody, obj)

			} else if s.configFile.fileType == fileTypeYAML {
				yamlBody, err := yaml.Marshal(value)
				if err != nil {
					return err
				}

				return yaml.Unmarshal(yamlBody, obj)
			}
		}
	}

	keyName := key
	if s.section != "" {
		keyName = s.section + "_" + key
	}

	if value, ok := os.LookupEnv(keyName); ok {
		err := json.Unmarshal([]byte(value), obj)
		if err != nil {
			err2 := yaml.Unmarshal([]byte(value), obj)
			if err2 != nil {
				return fmt.Errorf("unable to Unmarshal %s, as json:%s, or yaml:%s", key, err.Error(), err2.Error())
			}
		}
	}

	return nil
}

func (s *settings) GetSection(key string) ISettings {

	sectionName := key
	if s.section != "" {
		sectionName = s.section + "_" + key
	}

	var settings = &settings{cash: map[string]interface{}{}, configFile: s.configFile, section: sectionName}

	if s.configFile != nil {
		if value, ok := s.configFile.cash[key]; ok {
			valueMap := getMapInterface(value)
			if valueMap != nil {
				settings.configFile = &configFile{cash: valueMap, fileType: s.configFile.fileType}
			}
		}
	}

	return settings
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

	if s.configFile != nil {
		if value, ok := s.configFile.cash[key]; ok {
			s.cash[key] = value
			return s.cash[key]
		}
	}

	keyName := key
	if s.section != "" {
		keyName = s.section + "_" + key
	}

	if value, ok := os.LookupEnv(keyName); ok {
		s.cash[key] = value
		return s.cash[key]
	}

	return fallback

}

func loadConfigFile(configFileName string) *configFile {
	if configFileName != "" {
		file, err := os.Open(configFileName)
		if err == nil {
			defer file.Close()
			byteValue, err := ioutil.ReadAll(file)
			if err == nil {
				var result map[string]interface{}
				var fileType fileType
				if strings.HasSuffix(configFileName, ".json") {
					err = json.Unmarshal(byteValue, &result)
					fileType = fileTypeJSON

				} else if strings.HasSuffix(configFileName, ".yaml") {
					err = yaml.Unmarshal(byteValue, &result)
					fileType = fileTypeYAML
				}

				if err == nil {
					return &configFile{cash: result, fileType: fileType}
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
	settings.configFile = loadConfigFile(configFile)
	return settings
}
