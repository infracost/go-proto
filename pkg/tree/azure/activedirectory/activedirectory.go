package activedirectory

type ActiveDirectory struct {
	DomainServices           []DomainService           `tree:"domain_services"`
	DomainServiceReplicaSets []DomainServiceReplicaSet `tree:"domain_service_replica_sets"`
}

func (ad *ActiveDirectory) PostProcess() {
	// propagate SKU from domain services to replica sets
	for i, rs := range ad.DomainServiceReplicaSets {
		for _, ds := range ad.DomainServices {
			if rs.DomainServiceID.Value() == ds.ID {
				ad.DomainServiceReplicaSets[i].SKU = ds.SKU
				break
			}
		}
	}
}
