package gitea

import (
	"fmt"
	"log"
	"strconv"

	giteaapi "code.gitea.io/sdk/gitea"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGiteaRepositoryHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceGiteaRepositoryHookCreate,
		Read:   resourceGiteaRepositoryHookRead,
		Update: resourceGiteaRepositoryHookUpdate,
		Delete: resourceGiteaRepositoryHookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"owner": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
			},
			"repository": {
				Type:      schema.TypeString,
				Required:  true,
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

func resourceGiteaRepositoryHookSetToState(d *schema.ResourceData, hook *giteaapi.Hook) error {
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

func resourceGiteaRepositoryHookCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	owner := d.Get("owner").(string)
	repository := d.Get("repository").(string)

	object, err := resourceGiteaRepositoryHookCreateObject(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] create repo hook: %s %s %v", owner, repository, object)

	hook, err := client.CreateRepoHook(owner, repository, object)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] repo hook created %v", hook)
	d.SetId(strconv.FormatInt(hook.ID, 10))

	return resourceGiteaRepositoryHookRead(d, meta)
}

func resourceGiteaRepositoryHookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	log.Printf("[DEBUG] repo hook informations: %s", d.Id())
	hookId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return unconvertibleIdErr(d.Id(), err)
	}

	owner := d.Get("owner").(string)
	repository := d.Get("repository").(string)
	log.Printf("[DEBUG] read repo hook %q", hookId)

	hook, err := client.GetRepoHook(owner, repository, hookId)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] repo hook find %v", hook)
	return resourceGiteaRepositoryHookSetToState(d, hook)
}

func resourceGiteaRepositoryHookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	hookId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return unconvertibleIdErr(d.Id(), err)
	}

	owner := d.Get("owner").(string)
	repository := d.Get("repository").(string)

	object, err := resourceGiteaRepositoryHookUpdateObject(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] edit repository hook: %s %s %q %v", owner, repository, hookId, object)
	err = client.EditRepoHook(owner, repository, hookId, object)
	if err != nil {
	 	return err
	}

	return resourceGiteaRepositoryHookRead(d, meta)
}

func resourceGiteaRepositoryHookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return unconvertibleIdErr(d.Id(), err)
	}
	owner := d.Get("owner").(string)
	repository := d.Get("repository").(string)
	log.Printf("[DEBUG] delete repo hook: %d %s %s", id, owner, repository)
	return client.DeleteRepoHook(owner, repository, id)
}

func resourceGiteaRepositoryHookCreateObject(d *schema.ResourceData) (giteaapi.CreateHookOption, error) {
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

func resourceGiteaRepositoryHookUpdateObject(d *schema.ResourceData) (giteaapi.EditHookOption, error) {
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