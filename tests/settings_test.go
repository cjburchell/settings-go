package tests

import (
	"github.com/cjburchell/settings-go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
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
	s := settings.Get("")
	result1 := s.Get("testEnvString", "")
	result2 := s.GetInt("testEnvInt", 0)
	result3 := s.GetBool("testEnvBool", false)
	result4 := s.GetInt64("testEnvInt64", 0)
	assert.Equal(t, "test", result1, "String not equal")
	assert.Equal(t, 1, result2, "int not equal")
	assert.Equal(t, true, result3, "Bool not equal")
	assert.Equal(t, int64(2), result4, "Int64 not equal")
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
	d1 := []byte("{\n   \"testJsonString\":\"test\",\n   \"testJsonInt\":1,\n   \"testJsonBool\":true,\n   \"testJsonInt64\":2\n}")
	err := ioutil.WriteFile("config.json", d1, 0644)
	defer os.Remove("config.json")
	assert.Nil(t, err)
	s := settings.Get("config.json")
	result1 := s.Get("testJsonString", "")
	result2 := s.GetInt("testJsonInt", 0)
	result3 := s.GetBool("testJsonBool", false)
	result4 := s.GetInt64("testJsonInt64", 0)
	assert.Equal(t, "test", result1, "String not equal")
	assert.Equal(t, 1, result2, "int not equal")
	assert.Equal(t, true, result3, "Bool not equal")
	assert.Equal(t, int64(2), result4, "Int64 not equal")
}

func TestYaml(t *testing.T) {
	d1 := []byte("testYamlString: test\ntestYamlInt: 1\ntestYamlBool: true\ntestYamlInt64: 2")
	err := ioutil.WriteFile("config.yaml", d1, 0644)
	defer os.Remove("config.yaml")
	assert.Nil(t, err)
	s := settings.Get("config.yaml")
	result1 := s.Get("testYamlString", "")
	result2 := s.GetInt("testYamlInt", 0)
	result3 := s.GetBool("testYamlBool", false)
	result4 := s.GetInt64("testYamlInt64", 0)
	assert.Equal(t, "test", result1, "String not equal")
	assert.Equal(t, 1, result2, "int not equal")
	assert.Equal(t, true, result3, "Bool not equal")
	assert.Equal(t, int64(2), result4, "Int64 not equal")
}