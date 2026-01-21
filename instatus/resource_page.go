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
			"workspace_slug": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subdomain/slug for the status page",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Workspace ID returned by the Instatus API",
			},
			"logoUrl": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of the logo for the status page",
			},
			"faviconUrl": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of the favicon for the status page",
			},
			"googleAnalytics": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Google Analytics tracking ID for the status page",
			},
			"customDomain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom domain for the status page",
			},
		},
	}
}

func resourcePageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	page := &Page{
		Email:            d.Get("email").(string),
		Name:             d.Get("name").(string),
		Subdomain:        d.Get("workspace_slug").(string),
		logo_url:         d.Get("logoUrl").(string),
		favicon_url:      d.Get("faviconUrl").(string),
		google_analytics: d.Get("googleAnalytics").(string),
		custom_domain:    d.Get("customDomain").(string),
	}

	created, err := client.CreateStatusPage(page)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating status page: %w", err))
	}

	d.SetId(created.ID)
	if err := d.Set("workspace_slug", created.WorkspaceSlug); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("workspace_id", created.WorkspaceID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)
	var diags diag.Diagnostics

	pageID := d.Id()

	page, err := client.GetStatusPage(pageID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading status page: %w", err))
	}

	if err := d.Set("name", page.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("workspace_slug", page.WorkspaceSlug); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("workspace_id", page.WorkspaceID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	pageID := d.Id()

	page := &PageUpdate{
		Email: d.Get("email").(string),
		Name:  d.Get("name").(string),
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
	client := meta.(*Client)
	pageID := d.Id()
	workspaceID := d.Get("workspace_id").(string)

	err := client.DeleteStatusPage(pageID, workspaceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting status page: %w", err))
	}

	d.SetId("")
	return nil
}
