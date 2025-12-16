package instatus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the Instatus Terraform provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("INSTATUS_API_KEY", nil),
				Description: "The API key for Instatus API authentication",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"instatus_component": resourceComponent(),
			"instatus_page":      resourcePage(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)

	var diags diag.Diagnostics

	if apiKey == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing API Key",
			Detail:   "API key must be provided via the api_key provider argument or INSTATUS_API_KEY environment variable",
		})
	}

	if diags.HasError() {
		return nil, diags
	}

	client := NewClient(apiKey)

	return client, diags
}
