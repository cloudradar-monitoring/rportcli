package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoadFromEnv(t *testing.T) {
	vp := NewValuesProvider("somepath")
	err := changeEnvs(map[string]string{
		LoginEnvVar:     "admin",
		PasswordEnvVar:  "somePass",
		ServerURLEnvVar: "10.9.8.7",
	}, true)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer func() {
		err = changeEnvs(map[string]string{
			LoginEnvVar:     "admin",
			PasswordEnvVar:  "somePass",
			ServerURLEnvVar: "10.9.8.7",
		}, false)
		if err != nil {
			logrus.Error(err)
		}
	}()

	val, found := vp.Read(Login)
	assert.Equal(t, "admin", val)
	assert.True(t, found)

	val, found = vp.Read(Password)
	assert.Equal(t, "somePass", val)
	assert.True(t, found)

	val2, found2 := vp.Read(ServerURL)
	assert.Equal(t, "10.9.8.7", val2)
	assert.True(t, found2)

	_, found3 := vp.Read("unknownVal")
	assert.False(t, found3)
}

func TestLoadFromEnvAndFile(t *testing.T) {
	fileName := "testingFile.json"
	fileData := []byte(`{"login":"root","password":"root","server_url":"https://10.10.10.10:3000"}`)
	err := ioutil.WriteFile(fileName, fileData, 0600)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer func() {
		e := os.Remove(fileName)
		if e != nil {
			logrus.Error(e)
		}
	}()

	vp := NewValuesProvider(fileName)

	err = changeEnvs(map[string]string{
		LoginEnvVar:    "admin",
		PasswordEnvVar: "somePass",
	}, true)

	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer func() {
		err = changeEnvs(map[string]string{
			LoginEnvVar:    "admin",
			PasswordEnvVar: "somePass",
		}, false)
		if err != nil {
			logrus.Error(err)
		}
	}()

	val, found := vp.Read(Login)
	assert.Equal(t, "admin", val)
	assert.True(t, found)

	val1, found1 := vp.Read(Password)
	assert.Equal(t, "somePass", val1)
	assert.True(t, found1)

	val2, found2 := vp.Read(ServerURL)
	assert.Equal(t, "https://10.10.10.10:3000", val2)
	assert.True(t, found2)

	_, found3 := vp.Read("unknownVal")
	assert.False(t, found3)
}

func TestLoadFromBrokenFile(t *testing.T) {
	fileName := "testingFile2.json"
	fileData := []byte(`dlfasdj`)
	err := ioutil.WriteFile(fileName, fileData, 0600)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			logrus.Error(err)
		}
	}()

	vp := NewValuesProvider(fileName)

	_, found := vp.Read(Login)
	assert.False(t, found)
}

func changeEnvs(vals map[string]string, isCreate bool) error {
	var err error
	for k, v := range vals {
		if isCreate {
			err = os.Setenv(k, v)
		} else {
			err = os.Unsetenv(k)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
