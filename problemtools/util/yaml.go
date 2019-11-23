package util

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ParseYamlString(content string, ptr interface{}) error {
	return yaml.Unmarshal([]byte(content), ptr)
}

func ParseYaml(path string, ptr interface{}) error {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(dat), ptr)
	if err != nil {
		return err
	}
	return nil
}
