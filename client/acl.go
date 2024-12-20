package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/splunk/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

// https://docs.splunk.com/Documentation/Splunk/8.0.4/RESTUM/RESTusing#Access_Control_List
func (client *Client) GetAcl(owner, app, name string, resources ...string) (*http.Response, error) {
	resourcePath := []string{"servicesNS", owner, app}
	resourcePath = append(resourcePath, resources...)
	resourcePath = append(resourcePath, name, "acl")
	endpoint := client.BuildSplunkURL(nil, resourcePath...)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("GET failed for endpoint %s: %s", endpoint.Path, err)
	}

	return resp, nil
}

func (client *Client) ResourcesAndNameForPath(path string) (resources []string, name string, ok bool) {
	parts := strings.Split(path, "/")

	// an example path of apps/local/myapp would have parts:
	// * [apps, local] - parts[0:1]
	// * myapp         - parts[2]

	// 2 is the absolute minimum number of path parts that are valid to be parsed into resources and name
	if len(parts) < 2 {
		ok = false
		return
	}

	resources = parts[0 : len(parts)-1]
	name = parts[len(parts)-1]
	ok = true

	return
}

func (client *Client) UpdateAcl(owner, app, name string, acl *models.ACLObject, resources ...string) error {
	values, err := query.Values(&acl)
	if err != nil {
		return err
	}
	// remove app from url values during POST
	values.Del("app")
	values.Del("perms[read]")
	values.Del("perms[write]")
	// Flatten []string
	values.Set("perms.read", strings.Join(acl.Perms.Read, ","))
	values.Set("perms.write", strings.Join(acl.Perms.Write, ","))
	// Adding resources
	resourcePath := []string{"servicesNS", owner, app}
	resourcePath = append(resourcePath, resources...)
	resourcePath = append(resourcePath, name, "acl")
	endpoint := client.BuildSplunkURL(nil, resourcePath...)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return fmt.Errorf("GET failed for endpoint %s: %s", endpoint.Path, err)
	}

	defer resp.Body.Close()
	return nil
}

func (client *Client) Move(owner, app, name string, acl *models.ACLObject, resources ...string) error {
	values, err := query.Values(&acl)
	if err != nil {
		return err
	}
	// Adding resources
	resourcePath := []string{"servicesNS", owner, app}
	resourcePath = append(resourcePath, resources...)
	resourcePath = append(resourcePath, name, "move")
	endpoint := client.BuildSplunkURL(nil, resourcePath...)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
