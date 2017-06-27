package GoJSON

import "github.com/jeffail/gabs"

// JSONWrapper Struct wrapper to perform JSON operations
type JSONWrapper struct {
	ct *gabs.Container
}

// New Creates a new instance of JSONWrapper
func New(data string) (*JSONWrapper, error) {
	if container, err := gabs.ParseJSON([]byte(data)); err == nil {
		return &JSONWrapper{container}, nil
	} else {
		return nil, err
	}
}

// GetStringFromPath Returns a string from the path passed as parameter
func (wrapper *JSONWrapper) GetStringFromPath(path string) (string, bool) {
	value, ok := wrapper.ct.Path(path).Data().(string)
	return value, ok
}

// GetString Returns a string from the current object
func (wrapper *JSONWrapper) GetString() (string, bool) {
	value, ok := wrapper.ct.Data().(string)
	return value, ok
}

// GetIntFromPath Returns an int from the path passed as parameter
func (wrapper *JSONWrapper) GetIntFromPath(path string) (int, bool) {
	value, ok := wrapper.ct.Path(path).Data().(int)
	return value, ok
}

// GetInt Returns an int from the current object
func (wrapper *JSONWrapper) GetInt() (int, bool) {
	value, ok := wrapper.ct.Data().(int)
	return value, ok
}

// GetInt64FromPath Returns an int64 from the path passed as parameter
func (wrapper *JSONWrapper) GetInt64FromPath(path string) (int64, bool) {
	value, ok := wrapper.ct.Path(path).Data().(int64)
	return value, ok
}

// GetInt64 Returns an int64 from the current object
func (wrapper *JSONWrapper) GetInt64() (int64, bool) {
	value, ok := wrapper.ct.Data().(int64)
	return value, ok
}

// GetFloatFromPath Returns a float from the path passed as parameter
func (wrapper *JSONWrapper) GetFloatFromPath(path string) (float64, bool) {
	value, ok := wrapper.ct.Path(path).Data().(float64)
	return value, ok
}

// GetFloat Returns a float from the current object
func (wrapper *JSONWrapper) GetFloat() (float64, bool) {
	value, ok := wrapper.ct.Data().(float64)
	return value, ok
}

// GetBoolFromPath Returns a bool from the path passed as parameter
func (wrapper *JSONWrapper) GetBoolFromPath(path string) (bool, bool) {
	value, ok := wrapper.ct.Path(path).Data().(bool)
	return value, ok
}

// GetBool Returns a bool from the current object
func (wrapper *JSONWrapper) GetBool() (bool, bool) {
	value, ok := wrapper.ct.Data().(bool)
	return value, ok
}

// GetJSONObjectFromPath Returns a json object from the path passed as parameter
func (wrapper *JSONWrapper) GetJSONObjectFromPath(path string) *JSONWrapper {
	if wrapper.HasPath(path) {
		return &JSONWrapper{wrapper.ct.Path(path)}
	}
	return nil
}

// GetJSONObject Returns the current json object
func (wrapper *JSONWrapper) GetJSONObject() *JSONWrapper {
	return &JSONWrapper{wrapper.ct}
}

// CopyJSONObjectFromPath Returns a copy of the json object in the path passed
// as parameter
func (wrapper *JSONWrapper) CopyJSONObjectFromPath(path string) *JSONWrapper {
	if wrapper.HasPath(path) {
		return &JSONWrapper{wrapper.ct.Path(path)}
	}
	return nil
}

// GetArrayFromPath Returns an array from the path passed as parameter
func (wrapper *JSONWrapper) GetArrayFromPath(path string) []*JSONWrapper {
	array, _ := wrapper.ct.Path(path).Children()
	return createArrayWrapper(array)
}

// GetArray Returns an array from the current object
func (wrapper *JSONWrapper) GetArray() []*JSONWrapper {
	array, _ := wrapper.ct.Children()
	return createArrayWrapper(array)
}

func createArrayWrapper(array []*gabs.Container) []*JSONWrapper {
	var result []*JSONWrapper
	for _, item := range array {
		result = append(result, &JSONWrapper{item})
	}
	return result
}

// ArrayAppendCopy Appends and element to the array
func (wrapper *JSONWrapper) ArrayAppendCopy(json *JSONWrapper) error {
	return wrapper.ct.ArrayAppend(json.ct)
}

// SetValueAtPath Sets the value at the path passed as parameter
func (wrapper *JSONWrapper) SetValueAtPath(path string, value interface{}) {
	wrapper.ct.SetP(value, path)
}

// SetObjectAtPath Sets the object at the path passes as parameter
func (wrapper *JSONWrapper) SetObjectAtPath(path string, object *JSONWrapper) {
	wrapper.ct.SetP(object.ct.Data(), path)
}

// ToString Returns the string representation of the object
func (wrapper *JSONWrapper) ToString() string {
	return wrapper.ct.String()
}

// HasPath Returns true if the path exists false otherwise
func (wrapper *JSONWrapper) HasPath(path string) bool {
	return wrapper.ct.ExistsP(path)
}

func (wrapper *JSONWrapper) CreateJSONArrayAtPathWithArray(
	path string, array []*JSONWrapper) error {

	if _, err := wrapper.ct.ArrayP(path); err != nil {
		return err
	}
	for _, element := range array {
		wrapper.ct.ArrayAppendP(element.ct.Data(), path)
	}
	return nil
}

func (wrapper *JSONWrapper) CreateJSONArrayAtPath(
	path string) error {

	_, err := wrapper.ct.ArrayP(path)
	return err
}

func (wrapper *JSONWrapper) ArrayAppendInPath(
	path string, element *JSONWrapper) error {
	return wrapper.ct.ArrayAppendP(element.ct.Data(), path)
}

// FreeJSON Releases resources associated to the object
func (wrapper *JSONWrapper) FreeJSON() {
	wrapper = nil
}
