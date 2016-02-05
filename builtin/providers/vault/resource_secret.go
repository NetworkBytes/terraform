package vault

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/vault/api"
)

func resourceVaultSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceVaultSecretCreate,
		// Yay for PUT
		Update: resourceVaultSecretCreate,
		Read:   resourceVaultSecretRead,
		Delete: resourceVaultSecretDelete,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ttl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"data": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceVaultSecretCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	data := d.Get("data").(map[string]interface{})
	if ttl := d.Get("ttl").(string); ttl != "" {
		data["ttl"] = ttl
	}

	_, err := client.Logical().Write(d.Get("path").(string), data)
	if err != nil {
		return err
	}

	d.SetId(d.Get("path").(string))
	return nil
}

func resourceVaultSecretRead(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*api.Client)

	return nil
}

func isSecretNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "bad token")
}

func resourceVaultSecretDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	_, err := client.Logical().Delete(d.Get("path").(string))
	if err != nil {
		return err
	}

	return nil
}
