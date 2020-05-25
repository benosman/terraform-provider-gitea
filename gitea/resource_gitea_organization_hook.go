package gitea

import (
	"fmt"
	"log"
	"strconv"

	giteaapi "code.gitea.io/sdk/gitea"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGiteaOrganizationHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceGiteaOrganizationHookCreate,
		Read:   resourceGiteaOrganizationHookRead,
		Update: resourceGiteaOrganizationHookUpdate,
		Delete: resourceGiteaOrganizationHookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
			},
			"type": {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "gitea",
			},
			"branch_filter": {
				Type: 		schema.TypeString,
				Optional:   true,
			},
			"config": hookConfigurationSchema(),
			"url": {
				Type:      schema.TypeString,
				Computed:  true,
			},
			"events": {
				Type:      schema.TypeSet,
				Required:  true,
				Elem:      &schema.Schema{Type: schema.TypeString},
				Set:       schema.HashString,
			},
			"active": {
				Type:      schema.TypeBool,
				Optional:  true,
				Default:   true,
			},
		},
	}
}

func resourceGiteaOrganizationHookSetToState(d *schema.ResourceData, hook *giteaapi.Hook) error {
	d.Set("type", hook.Type)
	d.Set("url", hook.URL)
	d.Set("active", hook.Active)
	d.Set("events", hook.Events)

	// Gitea does not return the secret.
	// We want to store the secret in state, so we'll write the configuration
	// secret in state from what we get from ResourceData
	if len(d.Get("config").([]interface{})) > 0 {
		currentSecret := d.Get("config").([]interface{})[0].(map[string]interface{})["secret"]

		if currentSecret != "" {
			hook.Config["secret"] = fmt.Sprintf("%v", currentSecret)
		}
	}

	d.Set("config", []interface{}{hook.Config})

	return nil
}
/*
	ID      int64             `json:"id"`
	Config  map[string]string `json:"config"`
	Events  []string          `json:"events"`
	Active  bool              `json:"active"`
	Updated time.Time         `json:"updated_at"`
	Created time.Time         `json:"created_at"`

 */

func resourceGiteaOrganizationHookCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	organization := d.Get("organization").(string)

	object, err := resourceGiteaOrganizationHookCreateObject(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] create org hook: %s %s %v", organization, object)

	hook, err := client.CreateOrgHook(organization, object)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] org hook created %v", hook)
	d.SetId(strconv.FormatInt(hook.ID, 10))

	return resourceGiteaOrganizationHookRead(d, meta)
}

func resourceGiteaOrganizationHookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	log.Printf("[DEBUG] org hook informations: %s", d.Id())
	hookId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return unconvertibleIdErr(d.Id(), err)
	}

	organization := d.Get("organization").(string)
	log.Printf("[DEBUG] read org hook %q", hookId)

	hook, err := client.GetOrgHook(organization, hookId)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] repo hook find %v", hook)
	return resourceGiteaOrganizationHookSetToState(d, hook)
}

func resourceGiteaOrganizationHookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	hookId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return unconvertibleIdErr(d.Id(), err)
	}

	organization := d.Get("organization").(string)

	object, err := resourceGiteaOrganizationHookUpdateObject(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] edit repository hook: %s %q %v", organization, hookId, object)
	err = client.EditOrgHook(organization, hookId, object)
	if err != nil {
	 	return err
	}

	return resourceGiteaOrganizationHookRead(d, meta)
}

func resourceGiteaOrganizationHookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return unconvertibleIdErr(d.Id(), err)
	}
	organization := d.Get("organization").(string)
	log.Printf("[DEBUG] delete org hook: %d %s", id, organization)
	return client.DeleteOrgHook(organization, id)
}

func resourceGiteaOrganizationHookCreateObject(d *schema.ResourceData) (giteaapi.CreateHookOption, error) {
	hookType := d.Get("type").(string)
	branchFilter := d.Get("branch_filter").(string)
	active := d.Get("active").(bool)
	var events []string

	eventSet := d.Get("events").(*schema.Set)
	for _, v := range eventSet.List() {
		events = append(events, v.(string))
	}

	config := map[string]string{}
	configList := d.Get("config").([]interface{})
	for key, value := range configList[0].(map[string]interface{}) {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		if strValue != "" {
			config[strKey] = strValue
		}
	}

	hook := giteaapi.CreateHookOption{
		Type: hookType,
		BranchFilter: branchFilter,
		Events: events,
		Config: config,
		Active: active,
	}

	return hook, nil
}

func resourceGiteaOrganizationHookUpdateObject(d *schema.ResourceData) (giteaapi.EditHookOption, error) {
	branchFilter := d.Get("branch_filter").(string)
	active := d.Get("active").(bool)
	var events []string

	eventSet := d.Get("events").(*schema.Set)
	for _, v := range eventSet.List() {
		events = append(events, v.(string))
	}

	config := map[string]string{}
	configList := d.Get("config").([]interface{})
	for key, value := range configList[0].(map[string]interface{}) {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		if strValue != "" {
			config[strKey] = strValue
		}
	}

	hook := giteaapi.EditHookOption{
		BranchFilter: branchFilter,
		Events: events,
		Config: config,
		Active: &active,
	}

	return hook, nil
}