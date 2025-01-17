package globalaccelerator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/globalaccelerator"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfglobalaccelerator "github.com/hashicorp/terraform-provider-aws/internal/service/globalaccelerator"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccGlobalAcceleratorCustomRoutingEndpointGroup_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v globalaccelerator.CustomRoutingEndpointGroup
	resourceName := "aws_globalaccelerator_custom_routing_endpoint_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, globalaccelerator.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomRoutingEndpointGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRoutingEndpointGroupConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCustomRoutingEndpointGroupExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "destination_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "destination_configuration.0.from_port", "443"),
					resource.TestCheckResourceAttr(resourceName, "destination_configuration.0.protocols.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "destination_configuration.0.protocols.*", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "destination_configuration.0.to_port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "endpoint_configuration.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint_group_region"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGlobalAcceleratorCustomRoutingEndpointGroup_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v globalaccelerator.CustomRoutingEndpointGroup
	resourceName := "aws_globalaccelerator_custom_routing_endpoint_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, globalaccelerator.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomRoutingEndpointGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRoutingEndpointGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomRoutingEndpointGroupExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfglobalaccelerator.ResourceCustomRoutingEndpointGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccGlobalAcceleratorCustomRoutingEndpointGroup_endpointConfiguration(t *testing.T) {
	ctx := acctest.Context(t)
	var v globalaccelerator.CustomRoutingEndpointGroup
	resourceName := "aws_globalaccelerator_custom_routing_endpoint_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, globalaccelerator.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomRoutingEndpointGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRoutingEndpointGroupConfig_endpointConfiguration(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCustomRoutingEndpointGroupExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "destination_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "destination_configuration.0.from_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "destination_configuration.0.protocols.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "destination_configuration.0.protocols.*", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "destination_configuration.0.to_port", "8081"),
					resource.TestCheckResourceAttr(resourceName, "endpoint_configuration.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint_configuration.0.endpoint_id"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint_group_region"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCustomRoutingEndpointGroupExists(ctx context.Context, n string, v *globalaccelerator.CustomRoutingEndpointGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).GlobalAcceleratorConn()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Global Accelerator Custom Routing Endpoint Group ID is set")
		}

		output, err := tfglobalaccelerator.FindCustomRoutingEndpointGroupByARN(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckCustomRoutingEndpointGroupDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).GlobalAcceleratorConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_globalaccelerator_custom_routing_endpoint_group" {
				continue
			}

			_, err := tfglobalaccelerator.FindCustomRoutingEndpointGroupByARN(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Global Accelerator Custom Routing Endpoint Group %s still exists", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCustomRoutingEndpointGroupConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_globalaccelerator_custom_routing_accelerator" "test" {
  name = %[1]q
}

resource "aws_globalaccelerator_custom_routing_listener" "test" {
  accelerator_arn = aws_globalaccelerator_custom_routing_accelerator.test.id

  port_range {
    from_port = 443
    to_port   = 443
  }
}

resource "aws_globalaccelerator_custom_routing_endpoint_group" "test" {
  listener_arn = aws_globalaccelerator_custom_routing_listener.test.id

  destination_configuration {
    from_port = 443
    to_port   = 8443
    protocols = ["TCP"]
  }
}
`, rName)
}

func testAccCustomRoutingEndpointGroupConfig_endpointConfiguration(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_globalaccelerator_custom_routing_accelerator" "test" {
  name = %[1]q
}

resource "aws_globalaccelerator_custom_routing_listener" "test" {
  accelerator_arn = aws_globalaccelerator_custom_routing_accelerator.test.id

  port_range {
    from_port = 1
    to_port   = 65534
  }
}

resource "aws_globalaccelerator_custom_routing_endpoint_group" "test" {
  listener_arn = aws_globalaccelerator_custom_routing_listener.test.id

  destination_configuration {
    from_port = 8080
    to_port   = 8081
    protocols = ["TCP"]
  }

  endpoint_configuration {
    endpoint_id = aws_subnet.test.id
  }

  depends_on = [aws_internet_gateway.test]
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/28"

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName))
}
