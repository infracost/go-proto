package sagemaker

type SageMaker struct {
	NotebookInstances         []NotebookInstance         `tree:"notebook_instances"`
	EndpointConfigurations    []EndpointConfiguration    `tree:"endpoint_configurations"`
	Endpoints                 []Endpoint                 `tree:"endpoints"`
	Apps                      []App                      `tree:"apps"`
	Spaces                    []Space                    `tree:"spaces"`
	FeatureGroups             []FeatureGroup             `tree:"feature_groups"`
	MLflowTrackingServers     []MLflowTrackingServer     `tree:"mlflow_tracking_servers"`
	DataQualityJobDefinitions []DataQualityJobDefinition `tree:"data_quality_job_definitions"`
}

func (s *SageMaker) PostProcess() {
	for i := range s.Endpoints {
		s.Endpoints[i].Relationships.EndpointConfiguration = nil
	}

	// Match on config ID, or name for references by name.
	for i, endpoint := range s.Endpoints {
		ref := endpoint.EndpointConfigName.Value()
		if ref == "" {
			continue
		}
		for j := range s.EndpointConfigurations {
			c := &s.EndpointConfigurations[j]
			if c.ID == ref || (!c.Name.IsEmpty() && c.Name.Value() == ref) {
				s.Endpoints[i].Relationships.EndpointConfiguration = c
				break
			}
		}
	}
}
