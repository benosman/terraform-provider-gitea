package gitea

import (
	"fmt"
	"strings"

	giteaapi "code.gitea.io/sdk/gitea"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGiteaRepository() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGiteaRepositoryRead,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"fork": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_issues": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_wiki": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_pull_requests": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ignore_whitespace_conflicts": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_merge": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_rebase": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_rebase_merge": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_squash": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"archived": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"parent_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_name": {
				Type:     schema.TypeString,
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
			"website": {
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
			"default_branch": {
				Type:     schema.TypeString,
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

func dataSourceGiteaRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)

	username := strings.ToLower(d.Get("username").(string))
	name := d.Get("name").(string)

	repository, err := client.GetRepo(username, name)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", repository.ID))
	d.Set("name", repository.Name)
	d.Set("description", repository.Description)
	d.Set("full_name", repository.FullName)
	d.Set("description", repository.Description)
	d.Set("private", repository.Private)
	d.Set("fork", repository.Fork)
	if repository.Parent != nil {
		d.Set("parent_username", repository.Parent)
		d.Set("parent_name", repository.Parent)
	}
	d.Set("empty", repository.Empty)
	d.Set("mirror", repository.Mirror)
	d.Set("size", repository.Size)
	d.Set("html_url", repository.HTMLURL)
	d.Set("ssh_url", repository.SSHURL)
	d.Set("clone_url", repository.CloneURL)
	d.Set("website", repository.Website)
	d.Set("stars", repository.Stars)
	d.Set("forks", repository.Forks)
	d.Set("has_issues", repository.HasIssues)
	d.Set("has_wiki", repository.HasWiki)
	d.Set("has_pull_requests", repository.HasPullRequests)
	d.Set("ignore_whitespace_conflicts", repository.IgnoreWhitespaceConflicts)
	d.Set("allow_merge", repository.AllowMerge)
	d.Set("allow_rebase", repository.AllowRebase)
	d.Set("allow_rebase_merge", repository.AllowRebaseMerge)
	d.Set("allow_squash", repository.AllowSquash)
	d.Set("archived", repository.Archived)
	d.Set("watchers", repository.Watchers)
	d.Set("open_issue_count", repository.OpenIssues)
	d.Set("default_branch", repository.DefaultBranch)
	d.Set("created", repository.Created)
	d.Set("updated", repository.Updated)
	d.Set("permission_admin", repository.Permissions.Admin)
	d.Set("permission_push", repository.Permissions.Push)
	d.Set("permission_pull", repository.Permissions.Pull)
	return nil
}
