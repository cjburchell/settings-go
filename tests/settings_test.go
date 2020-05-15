package tests

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cjburchell/settings-go"
	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	err := os.Setenv("testEnvString", "test")
	assert.Nil(t, err)
	err = os.Setenv("testEnvInt", "1")
	assert.Nil(t, err)
	err = os.Setenv("testEnvBool", "true")
	assert.Nil(t, err)
	err = os.Setenv("testEnvInt64", "2")
	assert.Nil(t, err)
	err = os.Setenv("test_SubEnvString", "test5")
	assert.Nil(t, err)

	s := settings.Get("")
	result1 := s.Get("testEnvString", "")
	result2 := s.GetInt("testEnvInt", 0)
	result3 := s.GetBool("testEnvBool", false)
	result4 := s.GetInt64("testEnvInt64", 0)

	sub := s.GetSection("test")
	result5 := sub.Get("SubEnvString", "")

	assert.Equal(t, "test", result1, "String not equal")
	assert.Equal(t, 1, result2, "int not equal")
	assert.Equal(t, true, result3, "Bool not equal")
	assert.Equal(t, int64(2), result4, "Int64 not equal")
	assert.Equal(t, "test5", result5, "sub string not equal")
}

func TestFallback(t *testing.T) {
	s := settings.Get("")
	result1 := s.Get("testEnvStringFallback", "fallback")
	result2 := s.GetInt("testEnvIntFallback", 42)
	result3 := s.GetBool("testEnvBoolFallback", false)
	result4 := s.GetInt64("testEnvInt64Fallback", 43)
	assert.Equal(t, "fallback", result1, "String not equal")
	assert.Equal(t, 42, result2, "int not equal")
	assert.Equal(t, false, result3, "Bool not equal")
	assert.Equal(t, int64(43), result4, "Int64 not equal")
}

func TestCashed(t *testing.T) {
	err := os.Setenv("testEnvString2", "test")
	assert.Nil(t, err)
	s := settings.Get("")
	result1 := s.Get("testEnvString2", "")
	result2 := s.Get("testEnvString2", "")
	assert.Equal(t, "test", result1, "String not equal")
	assert.Equal(t, "test", result2, "Cashed value not equal")
}

func TestJson(t *testing.T) {
	d1 := []byte("{\n  \"testJsonString\": \"test\",\n  \"testJsonInt\": 1,\n  \"testJsonBool\": true,\n  \"testJsonInt64\": 2,\n  \"testJsonSub\": {\n    \"testJsonStringSub\": \"another test\"\n  }\n}")
	err := ioutil.WriteFile("config.json", d1, 0644)
	defer os.Remove("config.json")
	assert.Nil(t, err)
	s := settings.Get("config.json")
	result1 := s.Get("testJsonString", "")
	result2 := s.GetInt("testJsonInt", 0)
	result3 := s.GetBool("testJsonBool", false)
	result4 := s.GetInt64("testJsonInt64", 0)
	subSetting := s.GetSection("testJsonSub")
	assert.NotNil(t, subSetting)
	result5 := subSetting.Get("testJsonStringSub", "")
	assert.Equal(t, "test", result1, "String not equal")
	assert.Equal(t, 1, result2, "int not equal")
	assert.Equal(t, true, result3, "Bool not equal")
	assert.Equal(t, int64(2), result4, "Int64 not equal")
	assert.Equal(t, "another test", result5, "Substring not equal")
}

func TestYaml(t *testing.T) {
	d1 := []byte("testYamlString: test\ntestYamlInt: 1\ntestYamlBool: true\ntestYamlInt64: 2\ntestYamlSubValue:\n    testYamlStringSub: 'another test'")
	err := ioutil.WriteFile("config.yaml", d1, 0644)
	defer os.Remove("config.yaml")
	assert.Nil(t, err)
	s := settings.Get("config.yaml")
	result1 := s.Get("testYamlString", "")
	result2 := s.GetInt("testYamlInt", 0)
	result3 := s.GetBool("testYamlBool", false)
	result4 := s.GetInt64("testYamlInt64", 0)
	subSetting := s.GetSection("testYamlSubValue")
	assert.NotNil(t, subSetting)
	result5 := subSetting.Get("testYamlStringSub", "")
	assert.Equal(t, "test", result1, "String not equal")
	assert.Equal(t, 1, result2, "int not equal")
	assert.Equal(t, true, result3, "Bool not equal")
	assert.Equal(t, int64(2), result4, "Int64 not equal")
	assert.Equal(t, "another test", result5, "Substring not equal")
}
