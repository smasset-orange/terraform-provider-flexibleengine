package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccRtsSoftwareConfigV1_importBasic(t *testing.T) {
	resourceName := "flexibleengine_rts_software_config_v1.config_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRtsSoftwareConfigV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareConfigV1_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
