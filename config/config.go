package config

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"errors"
)

type cfg struct {
	vals map[string]interface{}
}
var c *cfg

// This gives us some sane typed JSON accessors
type JSONObject struct {
	vals map[string]interface{}
}

type JSONArray struct {
	vals []interface{}
}

func getObject(key string, index int, target interface{}) *JSONObject {
	obj := new(JSONObject)
	switch target := target.(type) {
		case map[string]interface{}:
			o, ok := target[key]
			if ok {
				msi, ok := o.(map[string]interface{})
				if ok {
					obj.vals = msi
					return obj
				}
			}
		case []interface{}:
			if index < len(target) {
				o := target[index]
				msi, ok := o.(map[string]interface{})
				if ok {
					obj.vals = msi
					return obj
				}
			}
	}
	return nil
}
func getArray(key string, index int, target interface{}) *JSONArray {
	arr := new(JSONArray)
	switch target := target.(type) {
		case map[string]interface{}:
			a, ok := target[key]
			if ok {
				ai, ok := a.([]interface{})
				if ok {
					arr.vals = ai
					return arr
				}
			}
		case []interface{}:
			if index < len(target) {
				a := target[index]
				ai, ok := a.([]interface{})
				if ok {
					arr.vals = ai
					return arr
				}
			}
	}
	return nil
}
func getString(key string, index int, target interface{}) *string {
	var s string
	switch target := target.(type) {
		case map[string]interface{}:
			s1, ok := target[key].(string)
			if ok {
				s = s1
			}
		case []interface{}:
			if index < len(target) {
				s2, ok := target[index].(string)
				if ok {
					s = s2
				}
			}
	}
	return &s
}

func getNum(key string, index int, target interface{}) *float64 {
	var n float64
	switch target := target.(type) {
		case map[string]interface{}:
			n1, ok := target[key].(float64)
			if ok {
				n = n1
			}
		case []interface{}:
			if index < len(target) {
				n2, ok := target[index].(float64)
				if ok {
					n = n2
				}
			}
	}
	return &n
}

// JSON Objects
func (self *JSONObject) GetObject(key string) *JSONObject {
	return getObject(key, 0, self.vals)
}
func (self *JSONObject) GetString(key string) *string {
	return getString(key, 0, self.vals)
}
func (self *JSONObject) GetNum(key string) *float64 {
	return getNum(key, 0, self.vals)
}
func (self *JSONObject) GetArray(key string) *JSONArray {
	return getArray(key, 0, self.vals)
}
// JSON arrays
func (self *JSONArray) GetObject(index int) *JSONObject {
	return getObject("", index, self.vals)
}
func (self *JSONArray) GetString(index int) *string {
	return getString("", index, self.vals)
}
func (self *JSONArray) GetNum(index int) *float64 {
	return getNum("", index, self.vals)
}
func (self *JSONArray) GetArray(index int) *JSONArray {
	return getArray("", index, self.vals)
}

func GetRoot() *JSONObject {
	root := new(JSONObject)
	root.vals = c.vals
	return root
}


func LoadFile(path, filetype *string) error {
	f, err := os.OpenFile(*path + "." + *filetype, os.O_RDONLY, 0655)
	if err != nil {
		return err
	}
	err, newc := parseConfigFile(f, filetype)
	if err != nil {
		return err
	}
	if c != nil {
		for key, value := range newc.vals {
			c.vals[key] = value
		}
	} else {
		c = newc
	}
	return nil
}

func parseConfigFile(config_file *os.File, filetype *string) (error, *cfg) {
	bytes, err := ioutil.ReadAll(config_file)
	if err != nil {
		return err, nil
	}
	newc := new(cfg)
	if *filetype == "json" {
		err = json.Unmarshal(bytes, &newc.vals)
		if err != nil {
			return err, nil
		}
	} else {
		return errors.New("Unrecognized config file type."), nil
	}

	return nil, newc
}

func LoadHierarchy(base_path string, filetype string, filenames ...string) error {
	for _, name := range filenames {
		fn := base_path + name
		err := LoadFile(&fn, &filetype)
		if err != nil {
			return err
		}
	}
	return nil
}