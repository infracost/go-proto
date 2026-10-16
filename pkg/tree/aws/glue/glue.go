package glue

type Glue struct {
	Crawlers               []Crawler               `tree:"crawlers"`
	CatalogDatabases       []CatalogDatabase       `tree:"catalog_databases"`
	Jobs                   []Job                   `tree:"jobs"`
	DevEndpoints           []DevEndpoint           `tree:"dev_endpoints"`
	MLTransforms           []MLTransform           `tree:"ml_transforms"`
	CatalogTableOptimizers []CatalogTableOptimizer `tree:"catalog_table_optimizers"`
}
