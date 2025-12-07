package instatus

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComponent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComponentCreate,
		ReadContext:   resourceComponentRead,
		UpdateContext: resourceComponentUpdate,
		DeleteContext: resourceComponentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the component",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description of the component",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "OPERATIONAL",
				Description: "The status of the component (OPERATIONAL, UNDERMAINTENANCE, DEGRADEDPERFORMANCE, PARTIALOUTAGE, MAJOROUTAGE)",
			},
			"show_uptime": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to show uptime for this component",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The order of the component (managed in UI)",
			},
			"grouped": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this component belongs to a group",
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the parent group (if grouped is true)",
			},
			"group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the parent group (read-only)",
			},
			"archived": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the component is archived",
			},
			"unique_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique email address for this component",
			},
		},
	}
}

func resourceComponentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	component := &Component{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Status:      d.Get("status").(string),
		ShowUptime:  d.Get("show_uptime").(bool),
		Grouped:     d.Get("grouped").(bool),
		Archived:    d.Get("archived").(bool),
	}

	// Only set order if explicitly provided
	if order, ok := d.GetOk("order"); ok {
		component.Order = order.(int)
	}

	// Handle group_id if provided
	if groupID, ok := d.GetOk("group_id"); ok {
		component.GroupID = groupID.(string)
	}

	created, err := client.CreateComponent(component)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating component: %w", err))
	}

	d.SetId(created.ID)

	return resourceComponentRead(ctx, d, meta)
}

func resourceComponentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)
	var diags diag.Diagnostics

	component, err := client.GetComponent(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading component: %w", err))
	}

	if err := d.Set("name", component.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", component.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status", component.Status); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("show_uptime", component.ShowUptime); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("order", component.Order); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("archived", component.Archived); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("unique_email", component.UniqueEmail); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_id", component.GroupIDRead); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_name", component.GroupName); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceComponentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	component := &Component{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Status:      d.Get("status").(string),
		ShowUptime:  d.Get("show_uptime").(bool),
		Archived:    d.Get("archived").(bool),
	}

	// Only set order if explicitly provided (don't overwrite UI-managed order)
	if order, ok := d.GetOk("order"); ok {
		component.Order = order.(int)
	}

	// Handle group_id for updates (uses groupId field)
	if groupID, ok := d.GetOk("group_id"); ok {
		component.Grouped = true
		component.GroupIDRead = groupID.(string)
	} else {
		component.Grouped = false
	}

	_, err := client.UpdateComponent(d.Id(), component)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating component: %w", err))
	}

	return resourceComponentRead(ctx, d, meta)
}

func resourceComponentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)
	var diags diag.Diagnostics

	err := client.DeleteComponent(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting component: %w", err))
	}

	d.SetId("")

	return diags
}
