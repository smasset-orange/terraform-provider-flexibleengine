package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/layer3/routers"
)

func TestAccNetworkingV2Router_basic(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("flexibleengine_networking_router_v2.router_1", &router),
				),
			},
			{
				Config: testAccNetworkingV2Router_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"flexibleengine_networking_router_v2.router_1", "name", "router_2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_update_external_gw(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFloatingIP(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_update_external_gw_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("flexibleengine_networking_router_v2.router_1", &router),
				),
			},
			{
				Config: testAccNetworkingV2Router_update_external_gw_2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"flexibleengine_networking_router_v2.router_1", "external_gateway", OS_EXTGW_ID),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_timeout(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("flexibleengine_networking_router_v2.router_1", &router),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_networking_router_v2" {
			continue
		}

		_, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2RouterExists(n string, router *routers.Router) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
		}

		found, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Router not found")
		}

		*router = *found

		return nil
	}
}

const testAccNetworkingV2Router_basic = `
resource "flexibleengine_networking_router_v2" "router_1" {
	name = "router_1"
	admin_state_up = "true"
	distributed = "false"
}
`

const testAccNetworkingV2Router_update = `
resource "flexibleengine_networking_router_v2" "router_1" {
	name = "router_2"
	admin_state_up = "true"
	distributed = "false"
}
`

const testAccNetworkingV2Router_update_external_gw_1 = `
resource "flexibleengine_networking_router_v2" "router_1" {
	name = "router"
	admin_state_up = "true"
	distributed = "false"
}
`

var testAccNetworkingV2Router_update_external_gw_2 = fmt.Sprintf(`
resource "flexibleengine_networking_router_v2" "router_1" {
	name = "router"
	admin_state_up = "true"
	distributed = "false"
	external_gateway = "%s"
}
`, OS_EXTGW_ID)

const testAccNetworkingV2Router_timeout = `
resource "flexibleengine_networking_router_v2" "router_1" {
	name = "router_1"
	admin_state_up = "true"
	distributed = "false"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
