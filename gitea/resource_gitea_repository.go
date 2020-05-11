package gitea

import (
	"fmt"
	"log"

	giteaapi "code.gitea.io/sdk/gitea"
	"github.com/hashicorp/terraform/helper/schema"
)

type EditRepoOptionHelper struct {
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Website string `json:"website,omitempty"`
	Private bool `json:"private,omitempty"`
	HasIssues bool `json:"has_issues,omitempty"`
	HasWiki bool `json:"has_wiki,omitempty"`
	DefaultBranch string `json:"default_branch,omitempty"`
	HasPullRequests bool `json:"has_pull_requests,omitempty"`
	IgnoreWhitespaceConflicts bool `json:"ignore_whitespace_conflicts,omitempty"`
	AllowMerge bool `json:"allow_merge_commits,omitempty"`
	AllowRebase bool `json:"allow_rebase,omitempty"`
	AllowRebaseMerge bool `json:"allow_rebase_explicit,omitempty"`
	AllowSquash bool `json:"allow_squash_merge,omitempty"`
	Archived bool `json:"archived,omitempty"`
}
func resourceGiteaRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceGiteaRepositoryCreate,
		Read:   resourceGiteaRepositoryRead,
		Update: resourceGiteaRepositoryUpdate,
		Delete: resourceGiteaRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"issue_labels": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"auto_init": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"git_ignores": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"license": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"readme": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"website": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"has_issues": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"has_wiki": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"default_branch": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"has_pull_requests": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"ignore_whitespace_conflicts": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"allow_merge": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"allow_rebase": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"allow_rebase_merge": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"allow_squash": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"archived": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceGiteaRepositorySetToState(d *schema.ResourceData, repo *giteaapi.Repository) {
	d.SetId(fmt.Sprintf("%d", repo.ID))
	d.Set("owner", repo.Owner.UserName)
	d.Set("name", repo.Name)
	d.Set("description", repo.Description)
	d.Set("is_private", repo.Private)
	d.Set("website", repo.Website)
	d.Set("has_issues", repo.HasIssues)
	d.Set("has_wiki", repo.HasWiki)
	d.Set("default_branch", repo.DefaultBranch)
	d.Set("has_pull_requests", repo.HasPullRequests)
	d.Set("ignore_whitespace_conflicts", repo.IgnoreWhitespaceConflicts)
	d.Set("allow_merge", repo.AllowMerge)
	d.Set("allow_rebase", repo.AllowRebase)
	d.Set("allow_rebase_merge", repo.AllowRebaseMerge)
	d.Set("allow_squash", repo.AllowSquash)
	d.Set("archived", repo.Archived)
}

func resourceGiteaRepositoryEditOptions(d *schema.ResourceData) giteaapi.EditRepoOption {
	edit := EditRepoOptionHelper{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Private:     d.Get("is_private").(bool),
		Website:	 d.Get("website").(string),
		HasIssues:   d.Get("has_issues").(bool),
		HasWiki:     d.Get("has_wiki").(bool),
		DefaultBranch: d.Get("default_branch").(string),
		HasPullRequests: d.Get("has_pull_requests").(bool),
		IgnoreWhitespaceConflicts:  d.Get("ignore_whitespace_conflicts").(bool),
		AllowMerge: d.Get("allow_merge").(bool),
		AllowRebase: d.Get("allow_rebase").(bool),
		AllowRebaseMerge: d.Get("allow_rebase_merge").(bool),
		AllowSquash: d.Get("allow_squash").(bool),
		Archived: d.Get("archived").(bool),
	}

	return giteaapi.EditRepoOption{
		Name:        &edit.Name,
		Description: &edit.Description,
		Private:     &edit.Private,
		Website:	 &edit.Website,
		HasIssues:   &edit.HasIssues,
		HasWiki:     &edit.HasWiki,
		DefaultBranch: &edit.DefaultBranch,
		HasPullRequests: &edit.HasPullRequests,
		IgnoreWhitespaceConflicts:  &edit.IgnoreWhitespaceConflicts,
		AllowMerge:  &edit.AllowMerge,
		AllowRebase:  &edit.AllowRebase,
		AllowRebaseMerge:  &edit.AllowRebaseMerge,
		AllowSquash:  &edit.AllowSquash,
		Archived:  &edit.Archived,
	}
}


func resourceGiteaRepositoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	owner := d.Get("owner").(string)
	// need to manage partial state as some properties can only be set on edit
	d.Partial(true)
	options := giteaapi.CreateRepoOption{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Private:     d.Get("is_private").(bool),
		AutoInit:    d.Get("auto_init").(bool),
		Gitignores:  d.Get("git_ignores").(string),
		License:     d.Get("license").(string),
		Readme:      d.Get("readme").(string),
	}

	log.Printf("[DEBUG] create repository %s", options.Name)

	repository, err := client.AdminCreateRepo(owner, options)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Repository created: %v", repository)
	d.SetId(fmt.Sprintf("%d", repository.ID))

	// Everything complete
	d.Partial(false)
	return nil
}

func resourceGiteaRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	owner := d.Get("owner").(string)
	name := d.Get("name").(string)
	log.Printf("[DEBUG] read repository %q %s %s", d.Id(), owner, name)
	repo, err := client.GetRepo(owner, name)
	if err != nil {
		return fmt.Errorf("unable to retrieve repository %s %s", owner, name)
	}
	log.Printf("[DEBUG] repository find: %v", repo)
	resourceGiteaRepositorySetToState(d, repo)
	return nil
}

func resourceGiteaRepositoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	owner := d.Get("owner").(string)
	name := d.Get("name").(string)
	log.Printf("[DEBUG] update repository %s", d.Id())

	edit := resourceGiteaRepositoryEditOptions(d)

	_, err := client.EditRepo(owner, name, edit)
	if err != nil {
		return fmt.Errorf("unable to edit organization: %s", name)
	}

	return resourceGiteaRepositoryRead(d, meta)
}

func resourceGiteaRepositoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	owner := d.Get("owner").(string)
	name := d.Get("name").(string)
	log.Printf("[DEBUG] delete repository: %s %s", owner, name)
	return client.DeleteRepo(owner, name)
}
