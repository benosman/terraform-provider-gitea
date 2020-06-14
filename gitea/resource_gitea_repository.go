package gitea

import (
	"fmt"
	"log"
	"strings"

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
		Importer: &schema.ResourceImporter{
			State: resourceGiteaRepositoryImportState,
		},

		Schema: map[string]*schema.Schema{
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private": {
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
			"fork": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"empty": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mirror": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"html_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clone_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stars": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"forks": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"watchers": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"open_issue_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission_admin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"permission_push": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"permission_pull": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceGiteaRepositorySetToState(d *schema.ResourceData, repo *giteaapi.Repository) {
	d.SetId(fmt.Sprintf("%d", repo.ID))
	d.Set("owner", repo.Owner.UserName)
	d.Set("name", repo.Name)
	d.Set("description", repo.Description)
	d.Set("full_name", repo.FullName)
	d.Set("private", repo.Private)
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
	d.Set("fork", repo.Fork)
	d.Set("empty", repo.Empty)
	d.Set("mirror", repo.Mirror)
	d.Set("size", repo.Size)
	d.Set("html_url", repo.HTMLURL)
	d.Set("ssh_url", repo.SSHURL)
	d.Set("clone_url", repo.CloneURL)
	d.Set("stars", repo.Stars)
	d.Set("forks", repo.Forks)
	d.Set("watchers", repo.Watchers)
	d.Set("open_issue_count", repo.OpenIssues)
	d.Set("created", repo.Created)
	d.Set("updated", repo.Updated)
	d.Set("permission_admin", repo.Permissions.Admin)
	d.Set("permission_push", repo.Permissions.Push)
	d.Set("permission_pull", repo.Permissions.Pull)
}

func resourceGiteaRepositoryEditOptions(d *schema.ResourceData) giteaapi.EditRepoOption {
	edit := EditRepoOptionHelper{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Private:     d.Get("private").(bool),
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
		Private:     d.Get("private").(bool),
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
	log.Printf("[DEBUG] Repository created (partial): %v", repository)
	d.SetId(fmt.Sprintf("%d", repository.ID))

	log.Printf("[DEBUG] update repository %s", d.Id())
	edit := resourceGiteaRepositoryEditOptions(d)
	_, err = client.EditRepo(owner, options.Name, edit)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Repository finalized: %v", repository)
	// Everything complete
	d.Partial(false)
	return resourceGiteaRepositoryRead(d, meta)
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

	if d.HasChange("owner") {
		log.Printf("[DEBUG] change owner of repository %s to %s", d.Id(), owner)
		o, _ := d.GetChange("owner")
		old := o.(string)
		transferOptions := giteaapi.TransferRepoOption{
			NewOwner: owner,
		}
		_, err := client.TransferRepo(old, name, transferOptions)
		if err != nil {
			 return err
		}
	}
	log.Printf("[DEBUG] update repository %s", d.Id())

	edit := resourceGiteaRepositoryEditOptions(d)

	_, err := client.EditRepo(owner, name, edit)
	if err != nil {
		return err
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

func resourceGiteaRepositoryImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {owner}/{name}", d.Id())
	}

	client := meta.(*giteaapi.Client)
	owner := parts[0]
	name := parts[1]
	repo, err := client.GetRepo(owner, name)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve repository %s %s", owner, name)
	}

	d.SetId(fmt.Sprintf("%d", repo.ID))
	resourceGiteaRepositorySetToState(d, repo)
	
	return []*schema.ResourceData{d}, nil
}