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
	DirectoryTypeUnknown DirectoryType = iota
	DirectoryTypeSimpleAD
	DirectoryTypeADConnector
	DirectoryTypeMicrosoftAD
)

type DirectorySize uint32

const (
	DirectorySizeUnknown DirectorySize = iota
	DirectorySizeLarge
	DirectorySizeSmall
)
