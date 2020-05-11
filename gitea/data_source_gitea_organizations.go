package gitea

import (
	"fmt"
	"log"

	giteaapi "code.gitea.io/sdk/gitea"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGiteaOrganizations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGiteaOrganizationsRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organizations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
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
				},
			},
		},
	}
}

func dataSourceGiteaOrganizationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	if data, ok := d.GetOk("username"); ok {
		username := data.(string)
		orgs, err := client.ListUserOrgs(username)
		if err != nil {
			return fmt.Errorf("unable to retrieve organizations for %s", username)
		}
		log.Printf("[DEBUG] organization find: %v", orgs)
		d.Set("organizations", flattenGiteaOrganizations(orgs))
		d.SetId(fmt.Sprintf("%d", orgs))
	} else {
		orgs, err := client.ListMyOrgs()
		if err != nil {
			return fmt.Errorf("unable to retrieve organizations for myself")
		}
		log.Printf("[DEBUG] organizations find: %v", orgs)
		d.Set("organizations", flattenGiteaOrganizations(orgs))
		d.SetId(fmt.Sprintf("%d", orgs))
	}

	return nil
}

func flattenGiteaOrganizations(orgs []*giteaapi.Organization) []interface{} {
	organizationsList := []interface{}{}

	for _, organization := range orgs {
		log.Printf("[DEBUG] organization flatten : %s", organization.UserName)
		values := map[string]interface{}{
			"id":          organization.ID,
			"username":    organization.UserName,
			"fullname":    organization.FullName,
			"avatar_url":  organization.AvatarURL,
			"description": organization.Description,
			"website":     organization.Website,
			"location":    organization.Location,
			"visibility":  organization.Visibility,
		}

		organizationsList = append(organizationsList, values)

	}
	return organizationsList
}
