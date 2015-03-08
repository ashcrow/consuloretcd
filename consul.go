package consuloretcd

// Consul specific implementation.

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// The useable structure for Consul
type Consul struct {
	Config
}

/*{	Endpoint string
	Client   http.Client
	Port     int
}*/

// Makes the URI from the Consul struct
// Returns the full URI as a string
func (c Consul) makeURI(name string) string {
	return c.Endpoint + ":" + strconv.Itoa(c.Port) + "/v1/kv/" + name
}

func (c Consul) checkAndReturn(resp *http.Response, kv KeyValue) (KeyValue, error) {
	kv.StatusCode = resp.StatusCode
	if resp.StatusCode != 200 {
		kv.Exists = false
		kv.Error = 2
		return kv, errors.New(Errors[kv.Error])
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		kv.Error = 3
		return kv, errors.New(Errors[kv.Error])
	}
	// Unmarshal the response
	var result []interface{}
	json.Unmarshal(body, &result)

	nr := result[0].(map[string]interface{})

	data, err := base64.StdEncoding.DecodeString(nr["Value"].(string))
	if err != nil {
		kv.Error = 4
		return kv, errors.New(Errors[kv.Error])
	}
	// Turn the result into a map of interfaces
	kv.CreateIndex = int(nr["CreateIndex"].(float64))
	kv.ModifyIndex = int(nr["ModifyIndex"].(float64))
	kv.Key = nr["Key"].(string)
	kv.Error = 0
	kv.Exists = true
	kv.Value = string(data)
	return kv, nil
}

// Gets a key from the remote Consul server.
// Returns KeyValue, nil on success
// Returns KeyValue, int (lookup via Errors) when unable to get a value
func (c Consul) GetKey(name string) (KeyValue, error) {
	kv := KeyValue{
		Name:   name,
		Exists: false}
	resp, err := c.Client.Get(c.makeURI(name))
	if err != nil {
		kv.Error = 1

		return kv, errors.New(Errors[kv.Error])
	}
	// Close the body at the end
	defer resp.Body.Close()
	return c.checkAndReturn(resp, kv)
}

// Gets a key from the remote Consul server.
// Returns KeyValue, nil on success
// Returns KeyValue, int (lookup via Errors) when unable to get a value
func (c Consul) PutKey(name string, value string) (KeyValue, error) {
	kv := KeyValue{
		Name:   name,
		Exists: false}
	req, _ := http.NewRequest("PUT", c.makeURI(name), strings.NewReader(value))
	resp, err := c.Client.Do(req)
	if err != nil {
		kv.Error = 1
		return kv, errors.New(Errors[kv.Error])
	}
	// Close the body at the end
	defer resp.Body.Close()
	kv.StatusCode = resp.StatusCode
	if resp.StatusCode != 200 {
		kv.Error = 2
		return kv, errors.New(Errors[kv.Error])
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		kv.Error = 3
		return kv, errors.New(Errors[kv.Error])
	}
	if string(body) == "true" {
		return c.GetKey(name)
	}
	kv.Error = 5
	return kv, errors.New(Errors[kv.Error])
}
