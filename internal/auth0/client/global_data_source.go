package client

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	internalSchema "github.com/auth0/terraform-provider-auth0/internal/schema"
)

// NewGlobalDataSource will return a new auth0_global_client data source.
func NewGlobalDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataGlobalClient,
		Schema:      globalDataSourceSchema(),
		Description: "Retrieve a tenant's global Auth0 application client.",
	}
}

func globalDataSourceSchema() map[string]*schema.Schema {
	dataSourceSchema := internalSchema.TransformResourceToDataSource(NewResource().Schema)
	delete(dataSourceSchema, "client_secret_rotation_trigger")
	return dataSourceSchema
}

func readDataGlobalClient(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if err := readGlobalClientID(ctx, d, m); err != nil {
		return err
	}
	return readClient(ctx, d, m)
}
