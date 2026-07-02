package s3

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type BucketPolicy struct {
	resource.Resource       `tree:"-"`
	BucketName              value.String `tree:"bucket"`
	DeniesInsecureTransport value.Bool   `tree:"denies_insecure_transport"`

	// SSLDeny holds the SSL-deny facts extracted from the policy document at parse
	// time. The final DeniesInsecureTransport decision is deferred to PostProcess,
	// which — once the owning bucket has been linked — matches the deny statement's
	// Resource ARNs against the bucket's identities (its resolved name AND its
	// synthetic instance token), so coverage is detected regardless of how either
	// side resolved. Transient: never serialized, and nil when the converter did
	// not populate it.
	SSLDeny *SSLDenyInfo `tree:"-"`
}

// SSLDenyInfo captures whether a bucket policy contains a valid non-SSL deny
// statement and which resources those statements cover, so that the
// bucket-coverage decision can be made later against the linked bucket rather
// than by string-matching whatever the policy's `bucket` argument happened to
// resolve to.
type SSLDenyInfo struct {
	// Present is true when the policy has at least one statement with Effect
	// "Deny", an action of "s3:*" or "*", and a Bool aws:SecureTransport=false
	// condition. This fact is independent of which bucket the policy protects.
	Present bool
	// Resources are the raw Resource ARNs listed across all such deny statements.
	Resources []string
	// Opaque is true when the policy document could not be parsed as JSON but was a
	// synthetic value (e.g. `policy = var.something`); such policies are given the
	// benefit of the doubt because their contents are genuinely unknowable.
	Opaque bool
}
