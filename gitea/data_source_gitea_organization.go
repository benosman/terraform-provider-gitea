package gitea

import (
	"fmt"
	"log"
	"strings"

	giteaapi "code.gitea.io/sdk/gitea"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGiteaOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGiteaOrganizationRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fullname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"avatar_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"website": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceGiteaOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	username := strings.ToLower(d.Get("username").(string))
	log.Printf("[DEBUG] read organization %q %s", d.Id(), username)
	org, err := client.GetOrg(username)
	if err != nil {
		return fmt.Errorf("unable to retrieve organization %s", username)
	}
	log.Printf("[DEBUG] organization find: %v", org)

	d.SetId(fmt.Sprintf("%d", org.ID))
	d.Set("username", org.UserName)
	d.Set("fullname", org.FullName)
	d.Set("avatar_url", org.AvatarURL)
	d.Set("description", org.Description)
	d.Set("website", org.Website)
	d.Set("location", org.Location)
	d.Set("visibility", org.Visibility)
	return nil
}
