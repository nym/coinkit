package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/lacker/coinkit/util"
)

// JSONObject is a modifiable set of key-value mappings.
// It is designed to be a convenient albeit possibly less efficient
// way of representing a JSON object of unknown format.
// After calling any exposed method, bytes and content should be equivalent.
// JSONObject works with both go's built-in JSON encoding and SQL encoding.
type JSONObject struct {
	bytes   []byte
	content map[string]interface{}
}

// Sets bytes based on content
func (ob *JSONObject) encode() {
	bytes, err := json.Marshal(ob.content)
	if err != nil {
		panic(err)
	}
	ob.bytes = bytes
}

func NewJSONObject(content map[string]interface{}) *JSONObject {
	answer := &JSONObject{
		content: content,
	}
	answer.encode()
	return answer
}

func ReadJSONObject(bytes []byte) (*JSONObject, error) {
	ob := &JSONObject{}
	err := ob.UnmarshalJSON(bytes)
	if err == nil {
		return ob, nil
	}
	return nil, err
}

func NewEmptyJSONObject() *JSONObject {
	content := make(map[string]interface{})
	return NewJSONObject(content)
}

func (ob *JSONObject) MarshalJSON() ([]byte, error) {
	return ob.bytes, nil
}

func (ob *JSONObject) UnmarshalJSON(bytes []byte) error {
	ob.bytes = bytes
	return json.Unmarshal(ob.bytes, &ob.content)
}

func (ob *JSONObject) Value() (driver.Value, error) {
	return driver.Value(ob.bytes), nil
}

func (ob *JSONObject) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("expected []byte")
	}
	return ob.UnmarshalJSON(bytes)
}

func (ob *JSONObject) Set(key string, value interface{}) {
	ob.content[key] = value
	ob.encode()
}

// Returns (nil, false) if the key does not exist
func (ob *JSONObject) Get(key string) (interface{}, bool) {
	value, ok := ob.content[key]
	return value, ok
}

// Returns (0, false) if the key does not exist, or is not int-y
func (ob *JSONObject) GetInt(key string) (int, bool) {
	value, ok := ob.Get(key)
	if ok {
		floatValue, ok := value.(float64)
		if ok {
			return int(floatValue), true
		}
		intValue, ok := value.(int)
		if ok {
			return intValue, true
		}
	}
	return 0, false
}

func (ob *JSONObject) DefaultInt(key string, def int) int {
	answer, ok := ob.GetInt(key)
	if ok {
		return answer
	}
	return def
}

func (ob *JSONObject) String() string {
	return string(util.PrettyJSON(ob.content))
}