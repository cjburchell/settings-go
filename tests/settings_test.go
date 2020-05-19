package tests

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cjburchell/settings-go"
	"github.com/stretchr/testify/assert"
)

func TestEnvString(t *testing.T) {
	err := os.Setenv("testEnvString", "test")
	assert.Nil(t, err)
	s := settings.Get("")
	result1 := s.Get("testEnvString", "")
	assert.Equal(t, "test", result1, "String not equal")
}

func TestEnvInt(t *testing.T) {
	err := os.Setenv("testEnvInt", "1")
	assert.Nil(t, err)
	s := settings.Get("")
	result := s.GetInt("testEnvInt", 0)
	assert.Equal(t, 1, result, "int not equal")
}

func TestEnvBool(t *testing.T) {
	err := os.Setenv("testEnvBool", "true")
	assert.Nil(t, err)
	s := settings.Get("")
	result := s.GetBool("testEnvBool", false)
	assert.Equal(t, true, result, "Bool not equal")
}

func TestEnvInt64(t *testing.T) {
	err := os.Setenv("testEnvInt64", "2")
	assert.Nil(t, err)
	s := settings.Get("")
	result4 := s.GetInt64("testEnvInt64", 0)
	assert.Equal(t, int64(2), result4, "Int64 not equal")
}

func TestEnvSubEnv(t *testing.T) {
	err := os.Setenv("test_SubEnvString", "test5")
	assert.Nil(t, err)
	s := settings.Get("")
	sub := s.GetSection("test")
	result5 := sub.Get("SubEnvString", "")
	assert.Equal(t, "test5", result5, "sub string not equal")
}

func TestEnvObjectYaml(t *testing.T) {
	err := os.Setenv("testEnvYamlObj", "testYaml: 'another test1'")
	assert.Nil(t, err)
	s := settings.Get("")

	yamlObj := testObject{}
	err = s.GetObject("testEnvYamlObj", &yamlObj)
	assert.Nil(t, err)

	assert.Equal(t, "another test1", yamlObj.TheValue, "yaml obj not equal")
}

func TestObjectJson(t *testing.T) {
	err := os.Setenv("testEnvJsonObj", "{\n    \"testJson\": \"another test3\"\n  }")
	assert.Nil(t, err)
	s := settings.Get("")

	jsonObj := testObject{}
	err = s.GetObject("testEnvJsonObj", &jsonObj)
	assert.Nil(t, err)

	assert.Equal(t, "another test3", jsonObj.TheValue, "json obj not equal")
}

func TestEnvBadObject(t *testing.T) {
	err := os.Setenv("testEnvBadObj", "this is a bad object")
	assert.Nil(t, err)
	badObj := testObject{}
	s := settings.Get("")
	err = s.GetObject("testEnvBadObj", &badObj)
	assert.NotNil(t, err)
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
	d1 := []byte(`
{
   "testJsonString":"test",
   "testJsonInt":1,
   "testJsonBool":true,
   "testJsonInt64":2,
   "testJsonSub":{
      "testJsonStringSub":"another test"
   },
   "testJsonObject":{
      "testJson":"another test3"
   },
   "testJsonArray":[
      {
         "testJson":"another test4"
      }
   ]
}`)
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
	obj := testObject{}
	err = s.GetObject("testJsonObject", &obj)
	assert.Nil(t, err)

	objArray := make([]testObject, 0)
	err = s.GetObject("testJsonArray", &objArray)
	assert.Nil(t, err)

	assert.Equal(t, "test", result1, "String not equal")
	assert.Equal(t, 1, result2, "int not equal")
	assert.Equal(t, true, result3, "Bool not equal")
	assert.Equal(t, int64(2), result4, "Int64 not equal")
	assert.Equal(t, "another test", result5, "Substring not equal")
	assert.Equal(t, "another test3", obj.TheValue, "Object value not equal")
	assert.Equal(t, 1, len(objArray), "Array Size not equal")
	assert.Equal(t, "another test4", objArray[0].TheValue, "Object array value not equal")
}

type testObject struct {
	TheValue string `json:"testJson" yaml:"testYaml"`
}

func TestYaml(t *testing.T) {
	d1 := []byte("testYamlString: test\n" +
		"testYamlInt: 1\n" +
		"testYamlBool: true\n" +
		"testYamlInt64: 2\n" +
		"testYamlSubValue:\n" +
		"    testYamlStringSub: 'another test'\n" +
		"testYamlObject:\n" +
		"    testYaml: 'another test1'")
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
	obj := testObject{}
	err = s.GetObject("testYamlObject", &obj)
	assert.Nil(t, err)
	assert.Equal(t, "test", result1, "String not equal")
	assert.Equal(t, 1, result2, "int not equal")
	assert.Equal(t, true, result3, "Bool not equal")
	assert.Equal(t, int64(2), result4, "Int64 not equal")
	assert.Equal(t, "another test", result5, "Substring not equal")
	assert.Equal(t, "another test1", obj.TheValue, "Object value not equal")
}
