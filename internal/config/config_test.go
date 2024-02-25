package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	LogFile  string `yaml:"log-file" validate:"required"`
	WarnFile string `yaml:"warn-file" validate:"required"`
}

var validTestData = `
log-file: Charles
warn-file: Bronson
`

var invalidTypeTestData = `
log-file: 12
warn-file: Bronson
`

var expectedData testConfig = testConfig{
	LogFile:  "Charles",
	WarnFile: "Bronson",
}

func TestReadConfigFromYAML(t *testing.T) {
	tempDir := t.TempDir()
	configFile, err := os.Create(tempDir + "/config.yaml")
	t.Cleanup(
		func() {
			configFile.Close()
		},
	)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	_, err = configFile.Write([]byte(validTestData))
	if !assert.NoError(t, err) {
		assert.Fail(t, err.Error())
	}

	conf, err := ReadConfigFromYAML[testConfig](configFile.Name())
	if assert.NoError(t, err) {
		assert.Equal(t, conf, &expectedData)
	}

	_, err = configFile.Write([]byte(invalidTypeTestData))
	if !assert.NoError(t, err) {
		assert.Fail(t, err.Error())
	}

	conf, err = ReadConfigFromYAML[testConfig](configFile.Name())
	if assert.Error(t, err) {
		assert.Nil(t, conf)
	}

	conf, err = ReadConfigFromYAML[testConfig]("file.txt")
	if assert.Error(t, err) {
		assert.Nil(t, conf)
	}
}

func TestValidateConfig(t *testing.T) {
	validConf := &testConfig{
		LogFile:  "somefile",
		WarnFile: "somefile",
	}
	invalidConf := &testConfig{
		WarnFile: "somefile",
	}

	err := ValidateConfig(validConf)
	assert.NoError(t, err)

	err = ValidateConfig(invalidConf)
	assert.Error(t, err)
}
