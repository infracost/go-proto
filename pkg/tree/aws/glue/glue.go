package glue

type Glue struct {
	Crawlers         []Crawler         `tree:"crawlers"`
	CatalogDatabases []CatalogDatabase `tree:"catalog_databases"`
	Jobs             []Job             `tree:"jobs"`
}
