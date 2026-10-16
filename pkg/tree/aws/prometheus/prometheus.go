package prometheus

import "github.com/infracost/go-proto/pkg/tree/resource"

type Prometheus struct {
	Workspaces []Workspace `tree:"workspaces"`
	Scrapers   []Scraper   `tree:"scrapers"`
}

type Workspace struct {
	resource.Resource `tree:"-"`
}

type Scraper struct {
	resource.Resource `tree:"-"`
}
