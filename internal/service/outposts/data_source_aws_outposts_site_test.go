package aws

import (
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/outposts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func TestAccAWSOutpostsSiteDataSource_Id(t *testing.T) {
	dataSourceName := "data.aws_outposts_site.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSOutpostsSites(t) },
		ErrorCheck:   acctest.ErrorCheck(t, outposts.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSOutpostsSiteDataSourceConfigId(),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrAccountID(dataSourceName, "account_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestMatchResourceAttr(dataSourceName, "id", regexp.MustCompile(`^os-.+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "name", regexp.MustCompile(`^.+$`)),
				),
			},
		},
	})
}

func TestAccAWSOutpostsSiteDataSource_Name(t *testing.T) {
	sourceDataSourceName := "data.aws_outposts_site.source"
	dataSourceName := "data.aws_outposts_site.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSOutpostsSites(t) },
		ErrorCheck:   acctest.ErrorCheck(t, outposts.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSOutpostsSiteDataSourceConfigName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "account_id", sourceDataSourceName, "account_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", sourceDataSourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id", sourceDataSourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", sourceDataSourceName, "name"),
				),
			},
		},
	})
}

func testAccAWSOutpostsSiteDataSourceConfigId() string {
	return `
data "aws_outposts_sites" "test" {}

data "aws_outposts_site" "test" {
  id = tolist(data.aws_outposts_sites.test.ids)[0]
}
`
}

func testAccAWSOutpostsSiteDataSourceConfigName() string {
	return `
data "aws_outposts_sites" "test" {}

data "aws_outposts_site" "source" {
  id = tolist(data.aws_outposts_sites.test.ids)[0]
}

data "aws_outposts_site" "test" {
  name = data.aws_outposts_site.source.name
}
`
}