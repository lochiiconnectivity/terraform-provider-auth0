package client_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/auth0/terraform-provider-auth0/internal/acctest"
)

const testAccDataGlobalClientConfig = `
%v
data auth0_global_client global {
}
`

func TestAccDataGlobalClient(t *testing.T) {
	acctest.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccGlobalClientConfigWithCustomLogin,
			},
			{
				Config: fmt.Sprintf(testAccDataGlobalClientConfig, testAccGlobalClientConfigWithCustomLogin),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.auth0_global_client.global", "custom_login_page", "<html>TEST123</html>"),
					resource.TestCheckResourceAttr("data.auth0_global_client.global", "custom_login_page_on", "true"),
					resource.TestCheckResourceAttrSet("data.auth0_global_client.global", "client_id"),
					resource.TestCheckResourceAttr("data.auth0_global_client.global", "app_type", ""),
					resource.TestCheckResourceAttr("data.auth0_global_client.global", "name", "All Applications"),
				),
			},
		},
	})
}
