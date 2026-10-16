package networkmanager

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Attachment struct {
	resource.Resource `tree:"-"`
	Type              value.Value[AttachmentType] `tree:"type"`
}

type AttachmentType uint32

const (
	AttachmentTypeUnknown AttachmentType = iota
	AttachmentTypeVPC
	AttachmentTypeConnect
	AttachmentTypeSiteToSiteVPN
)
