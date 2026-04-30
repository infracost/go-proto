package kms

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Key struct {
	resource.Resource     `tree:"-"`
	CustomerMasterKeySpec value.Value[KeySpec] `tree:"customer_master_key_spec"`
	EnableKeyRotation     value.Bool           `tree:"enable_key_rotation"`
}

type KeySpec uint32

const (
	KeySpecUnknown          KeySpec = iota
	KeySpecSymmetricDefault
	KeySpecRSA2048
	KeySpecRSA3072
	KeySpecRSA4096
	KeySpecECCNistP256
	KeySpecECCNistP384
	KeySpecECCNistP521
	KeySpecECCSecgP256K1
	KeySpecHMAC224
	KeySpecHMAC256
	KeySpecHMAC384
	KeySpecHMAC512
)
