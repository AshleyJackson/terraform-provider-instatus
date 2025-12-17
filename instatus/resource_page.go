package instatus

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePageCreate,
		ReadContext:   resourcePageRead,
		UpdateContext: resourcePageUpdate,
		DeleteContext: resourcePageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Billing email address for the status page workspace",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name for the status page workspace",
			},
			"subdomain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subdomain slug used for the public status page",
			},
			"workspace_slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server-generated slug returned by the Instatus API",
			},
		},
	}
}

func resourcePageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	page := &Page{
		Email:     d.Get("email").(string),
		Name:      d.Get("name").(string),
		Subdomain: d.Get("subdomain").(string),
	}

	created, err := client.CreateStatusPage(page)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating status page: %w", err))
	}

	d.SetId(created.ID)

	if err := d.Set("workspace_slug", created.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO: The public API currently lacks a GET endpoint for pages, so keep state as-is.
	client := meta.(*Client)
	var diags diag.Diagnostics

	pageID := d.Id()

	page, err := client.GetStatusPage(pageID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading status page: %w", err))
	}

	if err := d.Set("email", page.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", page.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("subdomain", page.Subdomain); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("workspace_slug", page.ID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	pageID := d.Id()

	page := &PageUpdate{
		Email:      d.Get("email").(string),
		Name:       d.Get("name").(string),
		Subdomain:  d.Get("subdomain").(string),
		Components: d.Get("components").([]Component),
	}

	_, err := client.UpdateStatusPage(pageID, page)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating status page: %w", err))
	}

	return diag.Diagnostics{}
	// TODO: Re-enable after Instatus Has a valid API response for GET page
	// return resourcePageRead(ctx, d, meta)
}

func resourcePageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO: Note. Deleting the Status Page does not delete the workspace. This will have to be manually cleaned up via the Instatus dashboard.
	client := meta.(*Client)
	pageID := d.Id()

	err := client.DeleteStatusPage(pageID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting status page: %w", err))
	}

	d.SetId("")
	return nil
}
