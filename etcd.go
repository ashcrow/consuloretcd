package consuloretcd

// Etcd specific implementation.

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// The useable structure for Etcd
type Etcd struct {
	Config
}

// Makes the URI from the Consul struct
// Returns the full URI as a string
func (c Etcd) makeURI(name string) string {
	return c.Endpoint + ":" + strconv.Itoa(c.Port) + "/v2/keys/" + name
}

// Gets a key from the remote Etcd server.
// Returns KeyValue, nil on success
// Returns KeyValue, int (lookup via Errors) when unable to get a value
func (c Etcd) GetKey(name string) (KeyValue, error) {
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
	if resp.StatusCode != 200 {
		kv.Error = 2
		kv.StatusCode = resp.StatusCode
		return kv, errors.New(Errors[kv.Error])
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		kv.Error = 3
		return kv, errors.New(Errors[kv.Error])
	}
	// Unmarshal the response
	var result interface{}
	json.Unmarshal(body, &result)

	nr := result.(map[string]interface{})
	node := nr["node"].(map[string]interface{})

	// Turn the result into a map of interfaces
	kv.CreateIndex = int(node["createdIndex"].(float64))
	kv.ModifyIndex = int(node["modifiedIndex"].(float64))
	kv.Key = node["key"].(string)
	kv.Error = 0
	kv.Exists = true
	kv.Value = node["value"]
	return kv, nil
}

// Puts a key on the remote Etcd server.
// Returns KeyValue, nil on success
// Returns KeyValue, ERROR_CODE when unable to get a value
func (c Etcd) PutKey(name string, value string) (KeyValue, error) {
	kv := KeyValue{
		Name:   name,
		Exists: false}
	values := url.Values{}
	values.Add("value", value)
	req, _ := http.NewRequest(
		"PUT",
		c.makeURI(name),
		strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded; param=value")

	resp, err := c.Client.Do(req)
	if err != nil {
		kv.Error = 1
		return kv, errors.New(Errors[kv.Error])
	}
	// Close the body at the end
	defer resp.Body.Close()
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
	var result interface{}
	json.Unmarshal(body, &result)

	nr := result.(map[string]interface{})

	if nr["action"] == "set" {
		node := nr["node"].(map[string]interface{})
		kv.Value = node["value"]
		kv.Exists = true
		kv.Key = node["key"].(string)
		kv.ModifyIndex = int(node["modifiedIndex"].(float64))
		kv.CreateIndex = int(node["createdIndex"].(float64))
		return kv, nil
	}
	kv.Error = 5
	return kv, errors.New(Errors[kv.Error])
}
