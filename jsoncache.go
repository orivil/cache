package cache

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"os"
	"gopkg.in/orivil/helper.v0"
	"fmt"
	"gopkg.in/yaml.v2"
)

type JsonCache struct {
	file string
	isYamlFile bool
}

func NewJsonCache(dir string, fileName string) (c *JsonCache, err error) {
	if !helper.IsExist(dir) {
		fmt.Println(dir)
		err = os.MkdirAll(dir, os.ModePerm)
	}
	var isYaml bool
	if ext := filepath.Ext(fileName); ext == ".yml" || ext == ".yaml" {
		isYaml = true
	}
	c = &JsonCache{
		file: filepath.Join(dir, fileName),
		isYamlFile: isYaml,
	}
	return
}

func (this *JsonCache) Write(inst interface{}) error {
	var (
		jsonData []byte
		err error
	)

	if !this.isYamlFile {
		jsonData, err = json.Marshal(inst)
	} else {
		jsonData, err = yaml.Marshal(inst)
	}

	if err == nil {
		err = ioutil.WriteFile(this.file, jsonData, 0777)
	}

	return err
}

func (this *JsonCache) Read(inst interface{}) error {
	var (
		jsonData []byte
		err error
	)
	if jsonData, err = ioutil.ReadFile(this.file); err == nil {
		// 有可能缓存文件不存在
		if !this.isYamlFile {
			err = json.Unmarshal(jsonData, inst)
		} else {
			err = yaml.Unmarshal(jsonData, inst)
		}
		return err
	}
	return nil
}
