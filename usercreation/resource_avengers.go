package usercreation

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	aclient "github.com/rupeshshewalkar/users-client"
)

func resourceUsers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUsersCreate,
		ReadContext:   resourceUsersRead,
		UpdateContext: resourceUsersUpdate,
		DeleteContext: resourceUsersDelete,
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
							Description: "full name of user",
						},
						"alias": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "any alias/nickname of user",
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
			"_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the _id value returned from mongodb",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "full name of username",
			},
			"alias": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "any alias/nickname of name",
			},
			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "his/her email",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "his/her username",
			},
			"deleted_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "deleted item count",
			},
			"matched_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total matched item found",
			},
			"modified_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total item modified",
			},
			"upserted_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total item upserted",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUsersCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	log.Printf("[DEBUG] %s: Beginning resourceUsersCreate", d.Id())
	var diags diag.Diagnostics
	c := m.(*ApiClient)

	name := d.Get("name").(string)
	alias := d.Get("alias").(string)
	email := d.Get("email").(string)
	username := d.Get("username").(string)

	a := aclient.User{
		Name:     name,
		Alias:    alias,
		Email:    email,
		Username: username,
	}

	res, err := c.usersclient.CreateUser(a)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("_id", res.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", res.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("alias", res.Alias); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", res.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("username", res.Username); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(res.ID)
	log.Printf("[DEBUG] %s: resourceUsersCreate finished successfully", d.Id())
	return diags
}

func resourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	log.Printf("[DEBUG] %s: Beginning resourceUsersRead", d.Id())
	var diags diag.Diagnostics
	c := m.(*ApiClient)
	res, err := c.usersclient.GetAllUsers()
	if err != nil {
		return diag.FromErr(err)
	}
	if res != nil {
		//As the return item is a []Users, lets Unmarshal it into "users"
		resItems := flattenUsers(&res)
		if err := d.Set("users", resItems); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.Errorf("no data found in db, insert one")
	}
	log.Printf("[DEBUG] %s: resourceUsersRead finished successfully", d.Id())
	return diags
}

func resourceUsersUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	log.Printf("[DEBUG] %s: Beginning resourceUsersUpdate", d.Id())
	var diags diag.Diagnostics
	c := m.(*ApiClient)

	name := d.Get("name").(string)
	alias := d.Get("alias").(string)
	email := d.Get("email").(string)
	username := d.Get("username").(string)

	a := aclient.User{
		Name:     name,
		Alias:    alias,
		Email:    email,
		Username: username,
	}
	res, err := c.usersclient.UpdateUserByUserName(a)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("matched_count", res.MatchedCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("modified_count", res.ModifiedCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("upserted_count", res.UpsertedCount); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: resourceUsersUpdate finished successfully", d.Id())
	return diags
}

func resourceUsersDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	log.Printf("[DEBUG] %s: Beginning resourceUsersDelete", d.Id())
	var diags diag.Diagnostics
	c := m.(*ApiClient)
	name := d.Get("name").(string)
	del, err := c.usersclient.DeleteUserByUserName(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deleted_count", del.DeletedCount); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[DEBUG] %s: resourceUsersDelete finished successfully", d.Id())
	return diags
}

func flattenUsers(usersList *[]aclient.User) []interface{} {
	if usersList != nil {
		users := make([]interface{}, len(*usersList))
		for i, user := range *usersList {
			al := make(map[string]interface{})

			al["_id"] = user.ID
			al["name"] = user.Name
			al["alias"] = user.Alias
			al["email"] = user.Email
			al["username"] = user.Username

			users[i] = al
		}
		return users
	}
	return make([]interface{}, 0)
}
