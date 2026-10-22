package google

import (
	"github.com/infracost/go-proto/pkg/tree/google/alloydb"
	"github.com/infracost/go-proto/pkg/tree/google/artifactregistry"
	"github.com/infracost/go-proto/pkg/tree/google/bigquery"
	"github.com/infracost/go-proto/pkg/tree/google/cloudfunctions"
	"github.com/infracost/go-proto/pkg/tree/google/cloudrun"
	"github.com/infracost/go-proto/pkg/tree/google/composer"
	"github.com/infracost/go-proto/pkg/tree/google/compute"
	"github.com/infracost/go-proto/pkg/tree/google/container"
	"github.com/infracost/go-proto/pkg/tree/google/datastream"
	"github.com/infracost/go-proto/pkg/tree/google/dns"
	"github.com/infracost/go-proto/pkg/tree/google/kms"
	"github.com/infracost/go-proto/pkg/tree/google/logging"
	"github.com/infracost/go-proto/pkg/tree/google/monitoring"
	"github.com/infracost/go-proto/pkg/tree/google/pubsub"
	"github.com/infracost/go-proto/pkg/tree/google/redis"
	"github.com/infracost/go-proto/pkg/tree/google/secretmanager"
	"github.com/infracost/go-proto/pkg/tree/google/servicenetworking"
	"github.com/infracost/go-proto/pkg/tree/google/spanner"
	"github.com/infracost/go-proto/pkg/tree/google/sql"
	"github.com/infracost/go-proto/pkg/tree/google/storage"
)

type Google struct {
	AlloyDB           alloydb.AlloyDB                     `tree:"alloydb"`
	ArtifactRegistry  artifactregistry.ArtifactRegistry   `tree:"artifactregistry"`
	BigQuery          bigquery.BigQuery                   `tree:"bigquery"`
	CloudFunctions    cloudfunctions.CloudFunctions       `tree:"cloudfunctions"`
	CloudRun          cloudrun.CloudRun                   `tree:"cloudrun"`
	Composer          composer.Composer                   `tree:"composer"`
	Compute           compute.Compute                     `tree:"compute"`
	Container         container.Container                 `tree:"container"`
	Datastream        datastream.Datastream               `tree:"datastream"`
	DNS               dns.DNS                             `tree:"dns"`
	KMS               kms.KMS                             `tree:"kms"`
	Logging           logging.Logging                     `tree:"logging"`
	Monitoring        monitoring.Monitoring               `tree:"monitoring"`
	PubSub            pubsub.PubSub                       `tree:"pubsub"`
	Redis             redis.Redis                         `tree:"redis"`
	SecretManager     secretmanager.SecretManager         `tree:"secretmanager"`
	ServiceNetworking servicenetworking.ServiceNetworking `tree:"servicenetworking"`
	Spanner           spanner.Spanner                     `tree:"spanner"`
	SQL               sql.SQL                             `tree:"sql"`
	Storage           storage.Storage                     `tree:"storage"`
}

func (g *Google) PostProcess() {
	// NOTE: Service-level PostProcess() methods are invoked automatically by the
	// reflective tree walker in tree.go. Only cross-service wiring belongs here.

	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range g.Compute.Addresses {
		g.Compute.Addresses[i].Relationships.Instance = nil
	}

	// cross-service: link compute addresses to instances by NATIP
	for i, addr := range g.Compute.Addresses {
		for j := range g.Compute.Instances {
			if addr.Address.Value() == g.Compute.Instances[j].NATIP.Value() {
				g.Compute.Addresses[i].Relationships.Instance = &g.Compute.Instances[j]
				break
			}
		}
	}
}
