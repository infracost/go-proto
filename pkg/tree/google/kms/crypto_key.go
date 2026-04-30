package kms

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CryptoKey struct {
	resource.Resource `tree:"-"`
	Algorithm         value.String                  `tree:"algorithm"`
	ProtectionLevel   value.Value[ProtectionLevel]  `tree:"protection_level"`
	RotationPeriod    value.String                  `tree:"rotation_period"`
}

type ProtectionLevel uint32

const (
	ProtectionLevelUnknown     ProtectionLevel = iota
	ProtectionLevelSoftware
	ProtectionLevelHSM
	ProtectionLevelExternal
	ProtectionLevelExternalVPC
)
