package usercreation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	aclient "github.com/rupeshshewalkar/users-client"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		//as this terraform provider is created for learning purpose, i kept it very simple.
		//we can access all APIs available at users-backend without any authentication and authorization.
		//so it will only contain ResourcesMap & DataSourcesMap, no schema
		ResourcesMap: map[string]*schema.Resource{
			"users_resource": resourceUsers(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"users_datasource": dataSourceUsers(),
		},
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("USERS_BACKEND_HOST_URL", "http://localhost:8000"),
			},
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return NewApiClient(d)

}

type ApiClient struct {
	data        *schema.ResourceData
	usersclient *aclient.Client
}

// NewApiClient will return a new instance of ApiClient using which we can communicate with Users-backend
func NewApiClient(d *schema.ResourceData) (*ApiClient, diag.Diagnostics) {
	c := &ApiClient{data: d}
	client, err := c.NewUsersClient()
	if err != nil {
		return c, diag.FromErr(err)
	}
	c.usersclient = client
	return c, nil

}

func (a *ApiClient) NewUsersClient() (*aclient.Client, error) {
	host := a.data.Get("host").(string)
	c, err := aclient.NewClient(&host)
	if err != nil {
		return c, err
	}
	return c, nil
}
