package gitea

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func hookConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"url": {
					Type:     schema.TypeString,
					Required: true,
				},
				"content_type": {
					Type:     schema.TypeString,
					Optional: true,
					Default: "json",
				},
				"secret": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
					ForceNew:  true,
				},
			},
		},
	}
}