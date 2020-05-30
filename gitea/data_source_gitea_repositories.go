package gitea

import (
	"fmt"
	"log"

	giteaapi "code.gitea.io/sdk/gitea"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGiteaRepositories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGiteaRepositoriesRead,

		Schema: map[string]*schema.Schema{
			"owner": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repositories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceGiteaRepositoriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	owner := d.Get("owner").(string)
	id := schema.HashString(owner)
	repos, err := client.ListUserRepos(owner)
	if err != nil {
		return fmt.Errorf("unable to retrieve repositories for %s", owner)
	}

	log.Printf("[DEBUG] repositories find: %v", repos)
	d.Set("repositories", flattenGiteaRepositories(repos))
	d.SetId(fmt.Sprintf("%d", id))

	return nil
}


func flattenGiteaRepositories(repos []*giteaapi.Repository) []interface{} {
	repoList := []interface{}{}

	for _, repo := range repos {
		log.Printf("[DEBUG] repository flatten : %s", repo.Name)
		values := map[string]interface{}{
			"id":          repo.ID,
			"name":    	repo.Name,
			"description": repo.Description,
			"full_name":    repo.FullName,
			"private":    repo.Private,
			"fork":    repo.Fork,
			"empty":    repo.Empty,
			"mirror":  repo.Mirror,
			"size":  repo.Size,
			"html_url":  repo.HTMLURL,
			"ssh_url":  repo.SSHURL,
			"clone_url":  repo.CloneURL,
			"website":     repo.Website,
			"stars":  repo.Stars,
			"forks":  repo.Forks,
			"watchers":  repo.Watchers,
			"open_issue_count":  repo.OpenIssues,
			"default_branch":  repo.DefaultBranch,
			"created":  fmt.Sprintf("%v", repo.Created),
			"updated":  fmt.Sprintf("%v", repo.Updated),
			"permission_admin":  repo.Permissions.Admin,
			"permission_push":  repo.Permissions.Push,
			"permission_pull":  repo.Permissions.Pull,
		}

		if repo.Parent != nil {
			values["parent_username"] = repo.Parent
			values["parent_name"] = repo.Parent
		}

		repoList = append(repoList, values)

	}
	return repoList
}
