package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/halradaideh/terraform-provider-puppetca/internal/provider"
)

var (
	resources []func(p *provider.Provider) resource.Resource
)

func Resources() []func(p *provider.Provider) resource.Resource {
	return resources
}
