package gitea

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
	"strings"

	giteaapi "code.gitea.io/sdk/gitea"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGiteaOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceGiteaOrganizationCreate,
		Read:   resourceGiteaOrganizationRead,
		Update: resourceGiteaOrganizationUpdate,
		Delete: resourceGiteaOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGiteaOrganizationImportState,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
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
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{"public", "limited", "private"}, false),
			},
			"website": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGiteaOrganizationSetToState(d *schema.ResourceData, org *giteaapi.Organization) error {
	if err := d.Set("name", org.UserName); err != nil {
		return err
	}
	if err := d.Set("full_name", org.FullName); err != nil {
		return err
	}
	if err := d.Set("description", org.Description); err != nil {
		return err
	}
	if err := d.Set("website", org.Website); err != nil {
		return err
	}
	return nil
}

func resourceGiteaOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	options := giteaapi.CreateOrgOption{
		UserName:    d.Get("name").(string),
		FullName:    d.Get("full_name").(string),
		Description: d.Get("description").(string),
		Website:     d.Get("website").(string),
		Location:    d.Get("location").(string),
		Visibility:  d.Get("visibility").(string),
	}

	log.Printf("[DEBUG] create organisation %q", options.UserName)

	org, err := client.CreateOrg(options)

	if err != nil {
		return fmt.Errorf("unable to create organization: %v", err)
	}
	log.Printf("[DEBUG] organization created: %v", org)
	d.SetId(fmt.Sprintf("%d", org.ID))
	return resourceGiteaOrganizationRead(d, meta)
}

func resourceGiteaOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	name := d.Get("name").(string)
	log.Printf("[DEBUG] read organization %q %s", d.Id(), name)
	org, err := client.GetOrg(name)
	if err != nil {
		return fmt.Errorf("unable to retrieve organization %s", name)
	}
	log.Printf("[DEBUG] organization find: %v", org)
	return resourceGiteaOrganizationSetToState(d, org)

}

func resourceGiteaOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	log.Printf("[DEBUG] update organization %s", d.Id())

	name := d.Get("name").(string)
	edit := giteaapi.EditOrgOption{
		FullName:    d.Get("full_name").(string),
		Description: d.Get("description").(string),
		Website:     d.Get("website").(string),
		Location:    d.Get("location").(string),
		Visibility:  d.Get("visibility").(string),
	}
	err := client.EditOrg(name, edit)
	if err != nil {
		return fmt.Errorf("unable to edit organization: %s", name)
	}

	return resourceGiteaOrganizationRead(d, meta)
}

func resourceGiteaOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	name := d.Get("name").(string)
	log.Printf("[DEBUG] delete organization: %s", name)
	return client.DeleteOrg(name)
}

func resourceGiteaOrganizationImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {id}/{name}", d.Id())
	}
	_ = d.Set("name", parts[1])
	d.SetId(parts[0])

	return []*schema.ResourceData{d}, nil
}
