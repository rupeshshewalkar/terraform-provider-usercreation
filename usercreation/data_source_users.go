package usercreation

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var HostURL string = "http://localhost:8000"

// dataSourceUsers is the Users data source which will pull information on all Users served by users-backend.
func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		//to read All Users, we can directly call resourceUsersRead()
		//which is implemented in resource_users.go file.
		//But watch the Schema, here KEYs are 'Computed: true' not 'Required: true'
		//because we don't want to provide these values while read.
		ReadContext: resourceUsersRead,
		Schema: map[string]*schema.Schema{
			"users": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the _id value returned from mongodb",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of user",
						},
						"alias": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "any alias of user",
						},
						"email": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "his/her email",
						},
						"username": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "his/her username",
						},
					},
				},
			},
		},
	}
}
