package directoryservice

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Directory struct {
	resource.Resource `tree:"-"`
	Type              value.Value[DirectoryType] `tree:"type"`
	Size              value.Value[DirectorySize] `tree:"size"`
}

type DirectoryType uint32

const (
	DirectoryTypeSimpleAD    DirectoryType = 0
	DirectoryTypeADConnector DirectoryType = 1
	DirectoryTypeMicrosoftAD DirectoryType = 2
)

type DirectorySize uint32

const (
	DirectorySizeLarge DirectorySize = 0
	DirectorySizeSmall DirectorySize = 1
)
