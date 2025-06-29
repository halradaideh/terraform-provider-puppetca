package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/halradaideh/terraform-provider-puppetca/internal/provider"
)

var (
	dataSources []func(p *provider.Provider) datasource.DataSource
)

func DataSources() []func(p *provider.Provider) datasource.DataSource {
	return dataSources
}
