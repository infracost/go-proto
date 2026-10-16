package securitylake

import "github.com/infracost/go-proto/pkg/tree/resource"

type SecurityLake struct {
	DataLakes []DataLake `tree:"data_lakes"`
}

type DataLake struct {
	resource.Resource `tree:"-"`
}
