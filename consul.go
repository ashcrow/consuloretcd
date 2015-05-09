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

// Makes the URI from the Consul struct
// Returns the full URI as a string
func (c Consul) makeURI(name string, opts KeyOptions) string {
	url := c.Endpoint + ":" + strconv.Itoa(c.Port) + "/v1/kv/" + name
	// TODO(ashcrow): This is a hack to avoid colliding with int:0. Fix it.
	if opts.CASet != "" {
		url = url + "?cas=" + opts.CASet
	}
	return url
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
func (c Consul) GetKey(name string, opts KeyOptions) (KeyValue, error) {
	kv := KeyValue{
		Name:   name,
		Exists: false}
	resp, err := c.Client.Get(c.makeURI(name, opts))
	if err != nil {
		kv.Error = 1

		return kv, errors.New(Errors[kv.Error])
	}
	// Close the body at the end
	defer resp.Body.Close()
	return c.checkAndReturn(resp, kv)
}

// Puts a key from the remote Consul server.
// Returns KeyValue, nil on success
// Returns KeyValue, int (lookup via Errors) when unable to get a value
func (c Consul) PutKey(name string, value string, opts KeyOptions) (KeyValue, error) {
	kv := KeyValue{
		Name:   name,
		Exists: false}
	req, _ := http.NewRequest("PUT", c.makeURI(name, opts), strings.NewReader(value))
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
		return c.GetKey(name, KeyOptions{})
	}
	kv.Error = 5
	return kv, errors.New(Errors[kv.Error])
}

// Deletes a key from the remote Consul server.
// Returns nil on success
// Returns Error when unable to delete
func (c Consul) DeleteKey(name string, opts KeyOptions) error {
	req, _ := http.NewRequest("DELETE", c.makeURI(name, opts), nil)
	resp, err := c.Client.Do(req)
	if err != nil {
		return errors.New(Errors[6])
	}
	// Close the body at the end
	defer resp.Body.Close()
	// It's weird but 200 means it deleted ...
	if resp.StatusCode != 200 {
		return errors.New(Errors[6])
	}
	return nil
}
