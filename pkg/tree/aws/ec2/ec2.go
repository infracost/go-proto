package ec2

type EC2 struct {
	AutoscalingGroups                []AutoscalingGroup                `tree:"autoscaling_groups"`
	AutoscalingSchedules             []AutoscalingSchedule             `tree:"autoscaling_schedules"`
	ClassicLoadBalancers             []ClassicLoadBalancer             `tree:"classic_load_balancers"`
	ClientVPNEndpoints               []ClientVPNEndpoint               `tree:"client_vpn_endpoints"`
	ClientVPNNetworkAssociations     []ClientVPNNetworkAssociation     `tree:"client_vpn_network_associations"`
	DefaultSecurityGroups            []DefaultSecurityGroup            `tree:"default_security_groups"`
	EBSSnapshots                     []EBSSnapshot                     `tree:"ebs_snapshots"`
	EBSSnapshotCopies                []EBSSnapshotCopy                 `tree:"ebs_snapshot_copies"`
	EBSVolumes                       []EBSVolume                       `tree:"ebs_volumes"`
	ElasticIPs                       []ElasticIP                       `tree:"elastic_ips"`
	ElasticIPAssociations            []ElasticIPAssociation            `tree:"elastic_ip_associations"`
	Hosts                            []Host                            `tree:"hosts"`
	Instances                        []Instance                        `tree:"instances"`
	InstanceStates                   []InstanceStateMapping            `tree:"instance_states"`
	LaunchConfigurations             []LaunchConfiguration             `tree:"launch_configurations"`
	LaunchTemplates                  []LaunchTemplate                  `tree:"launch_templates"`
	LBListeners                      []LBListener                      `tree:"lb_listeners"`
	LoadBalancers                    []LoadBalancer                    `tree:"load_balancers"`
	NATGateways                      []NATGateway                      `tree:"nat_gateways"`
	Subnets                          []Subnet                          `tree:"subnets"`
	TrafficMirrorSessions            []TrafficMirrorSession            `tree:"traffic_mirror_sessions"`
	TransitGateways                  []TransitGateway                  `tree:"transit_gateways"`
	TransitGatewayPeeringAttachments []TransitGatewayPeeringAttachment `tree:"transit_gateway_peering_attachments"`
	TransitGatewayVPCAttachments     []TransitGatewayVPCAttachment     `tree:"transit_gateway_vpc_attachments"`
	VPCs                             []VPC                             `tree:"vpcs"`
	VPCEndpoints                     []VPCEndpoint                     `tree:"vpc_endpoints"`
	VPNConnections                   []VPNConnection                   `tree:"vpn_connections"`
}

func (ec2 *EC2) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range ec2.Instances {
		ec2.Instances[i].Relationships = InstanceRelationships{}
	}
	for i := range ec2.AutoscalingGroups {
		ec2.AutoscalingGroups[i].Relationships = AutoscalingGroupRelationships{}
	}
	for i := range ec2.VPCEndpoints {
		ec2.VPCEndpoints[i].Relationships = VPCEndpointRelationships{}
	}
	for i := range ec2.VPCs {
		ec2.VPCs[i].Relationships.VPCEndpoints = nil
	}
	for i := range ec2.Subnets {
		ec2.Subnets[i].Relationships.NATGateways = nil
	}
	for i := range ec2.NATGateways {
		ec2.NATGateways[i].Relationships = NATGatewayRelationships{}
	}
	for i := range ec2.ElasticIPs {
		ec2.ElasticIPs[i].Relationships = ElasticIPRelationships{}
	}
	for i := range ec2.ElasticIPAssociations {
		ec2.ElasticIPAssociations[i].Relationships = ElasticIPAssociationRelationships{}
	}
	for i := range ec2.TransitGatewayVPCAttachments {
		ec2.TransitGatewayVPCAttachments[i].Relationships = TransitGatewayVPCAttachmentRelationships{}
	}
	for i := range ec2.TransitGatewayPeeringAttachments {
		ec2.TransitGatewayPeeringAttachments[i].Relationships = TransitGatewayPeeringAttachmentRelationships{}
	}
	for i := range ec2.LoadBalancers {
		for j := range ec2.LoadBalancers[i].SubnetMappings {
			ec2.LoadBalancers[i].SubnetMappings[j].Relationships = SubnetMappingRelationships{}
		}
	}

	// link instance states to instances
	for i, instanceState := range ec2.InstanceStates {
		for j := range ec2.Instances {
			if instanceState.InstanceID.Equal(ec2.Instances[j].ID) {
				ec2.Instances[j].Relationships.InstanceState = &ec2.InstanceStates[i]
				break
			}
		}
	}

	// link launch templates to instances
	for i, instance := range ec2.Instances {
		for j := range ec2.LaunchTemplates {
			lt := &ec2.LaunchTemplates[j]
			if (instance.LaunchTemplateID.IsEmpty() || !instance.LaunchTemplateID.Equal(lt.ID)) &&
				(instance.LaunchTemplateName.IsEmpty() || instance.LaunchTemplateName.Value() != lt.Name.Value()) &&
				(instance.LaunchTemplateID.IsEmpty() || instance.LaunchTemplateID.Value() != lt.Name.Value()) {
				continue
			}

			// store the launch template on the instance
			instance.Relationships.LaunchTemplate = lt

			if instance.Type.IsDefaultOrEmpty() && !lt.InstanceType.IsEmpty() {
				instance.Type = lt.InstanceType
			}

			if instance.EBSOptimized.IsDefault() {
				instance.EBSOptimized = lt.EBSOptimized
			}

			if instance.MonitoringEnabled.IsDefault() {
				instance.MonitoringEnabled = lt.MonitoringEnabled
			}

			if instance.CPUCredits.IsDefaultOrEmpty() && !lt.CPUCredits.IsEmpty() {
				instance.CPUCredits = lt.CPUCredits
			}

			if instance.Tenancy.IsDefaultOrEmpty() && !lt.Tenancy.IsEmpty() {
				instance.Tenancy = lt.Tenancy
			}

			// Iterate the original slices for deterministic order — Go map
			// iteration is randomised, and downstream consumers index into this
			// slice to name sub-resources, so unstable order produces flaky
			// output (instance-defined and LT-only devices can swap positions).
			// Maps are only used for lookups during the merge.
			//
			// TODO: LT-derived block devices currently inherit their address
			// from the launch template (e.g. `aws_launch_template.X.block_device_mapping[0]`).
			// Downstream the providers package falls back to a synthetic
			// `ebs_block_device[N]` name for these because the address doesn't
			// start with the instance's. It would be more accurate to surface
			// the actual LT-prefixed address so users can see where the device
			// is really defined — left for a follow-up since it changes
			// user-visible output.
			launchTemplateBlockDeviceMap := make(map[string]BlockDeviceMapping, len(lt.BlockDeviceMappings))
			for _, blockDevice := range lt.BlockDeviceMappings {
				launchTemplateBlockDeviceMap[blockDevice.DeviceName.Value()] = blockDevice
			}

			instanceBlockDeviceNames := make(map[string]struct{}, len(instance.BlockDeviceMappings))
			for _, blockDevice := range instance.BlockDeviceMappings {
				instanceBlockDeviceNames[blockDevice.DeviceName.Value()] = struct{}{}
			}

			blockDevices := make([]BlockDeviceMapping, 0, len(instance.BlockDeviceMappings)+len(lt.BlockDeviceMappings))

			// instance-defined first, merging any defaults from the launch template
			for _, blockDevice := range instance.BlockDeviceMappings {
				if launchTemplateBlockDevice, ok := launchTemplateBlockDeviceMap[blockDevice.DeviceName.Value()]; ok {
					if blockDevice.EBSVolume.Type.IsDefaultOrEmpty() {
						blockDevice.EBSVolume.Type = launchTemplateBlockDevice.EBSVolume.Type
					}
					if blockDevice.EBSVolume.IOPS.IsDefaultOrEmpty() {
						blockDevice.EBSVolume.IOPS = launchTemplateBlockDevice.EBSVolume.IOPS
					}
					if blockDevice.EBSVolume.SizeGB.IsDefaultOrEmpty() {
						blockDevice.EBSVolume.SizeGB = launchTemplateBlockDevice.EBSVolume.SizeGB
					}
					if blockDevice.EBSVolume.ThroughputMiBperS.IsDefaultOrEmpty() {
						blockDevice.EBSVolume.ThroughputMiBperS = launchTemplateBlockDevice.EBSVolume.ThroughputMiBperS
					}
				}
				blockDevices = append(blockDevices, blockDevice)
			}

			// LT-only block devices appended after, in launch template order
			for _, blockDevice := range lt.BlockDeviceMappings {
				if _, ok := instanceBlockDeviceNames[blockDevice.DeviceName.Value()]; !ok {
					blockDevices = append(blockDevices, blockDevice)
				}
			}

			instance.BlockDeviceMappings = blockDevices

			// save changes
			ec2.Instances[i] = instance
			break
		}
	}

	// link autoscaling schedules to autoscaling groups
	for i := range ec2.AutoscalingSchedules {
		schedule := &ec2.AutoscalingSchedules[i]
		if schedule.AutoscalingGroupName.IsEmpty() {
			continue
		}
		for j := range ec2.AutoscalingGroups {
			asg := &ec2.AutoscalingGroups[j]
			if schedule.AutoscalingGroupName.Value() == asg.Name.Value() || schedule.AutoscalingGroupName.Equal(asg.ID) {
				asg.Relationships.AutoscalingSchedules = append(asg.Relationships.AutoscalingSchedules, schedule)
				break
			}
		}
	}

	// link launch templates and launch configurations to autoscaling groups
	for i, asg := range ec2.AutoscalingGroups {
		for j := range ec2.LaunchTemplates {
			lt := &ec2.LaunchTemplates[j]
			if (!asg.LaunchTemplateID.IsEmpty() && asg.LaunchTemplateID.Equal(lt.ID)) ||
				(!asg.LaunchTemplateName.IsEmpty() && asg.LaunchTemplateName.Value() == lt.Name.Value()) ||
				(!asg.LaunchTemplateID.IsEmpty() && asg.LaunchTemplateID.Value() == lt.Name.Value()) {
				ec2.AutoscalingGroups[i].Relationships.LaunchTemplate = lt
			}
			if !asg.MixedInstanceLaunchTemplateID.IsEmpty() && asg.MixedInstanceLaunchTemplateID.Equal(lt.ID) {
				ec2.AutoscalingGroups[i].Relationships.MixedInstanceLaunchTemplate = lt
			}
		}
		for j := range ec2.LaunchConfigurations {
			lc := &ec2.LaunchConfigurations[j]
			if !asg.LaunchConfigurationName.IsEmpty() &&
				(asg.LaunchConfigurationName.Value() == lc.Name.Value() || asg.LaunchConfigurationName.Equal(lc.ID)) {
				ec2.AutoscalingGroups[i].Relationships.LaunchConfiguration = lc
			}
		}
	}

	// link VPC endpoints to VPCs
	for i, vpcEndpoint := range ec2.VPCEndpoints {
		for j := range ec2.VPCs {
			if ec2.VPCs[j].ID == vpcEndpoint.VPCID.Value() {
				ec2.VPCEndpoints[i].Relationships.VPC = &ec2.VPCs[j]
				ec2.VPCs[j].Relationships.VPCEndpoints = append(ec2.VPCs[j].Relationships.VPCEndpoints, &ec2.VPCEndpoints[i])
				break
			}
		}
	}

	// link VPC endpoints to subnets
	for i, vpcEndpoint := range ec2.VPCEndpoints {
		for j := range ec2.Subnets {
			if vpcEndpoint.SubnetIDs.Contains(ec2.Subnets[j].ID) {
				ec2.VPCEndpoints[i].Relationships.Subnets = append(ec2.VPCEndpoints[i].Relationships.Subnets, &ec2.Subnets[j])
			}
		}
	}

	// link EBS snapshots to volumes (inherit size)
	for i, snapshot := range ec2.EBSSnapshots {
		if !snapshot.VolumeID.IsEmpty() {
			for _, volume := range ec2.EBSVolumes {
				if volume.ID == snapshot.VolumeID.Value() {
					if snapshot.SizeGB.IsDefaultOrEmpty() {
						ec2.EBSSnapshots[i].SizeGB = volume.SizeGB
					}
					break
				}
			}
		}
	}

	// link EBS snapshot copies to source snapshots (inherit size)
	for i, snapshotCopy := range ec2.EBSSnapshotCopies {
		if !snapshotCopy.SourceSnapshotID.IsEmpty() {
			for _, snapshot := range ec2.EBSSnapshots {
				if snapshot.ID == snapshotCopy.SourceSnapshotID.Value() {
					if snapshotCopy.SizeGB.IsDefaultOrEmpty() {
						ec2.EBSSnapshotCopies[i].SizeGB = snapshot.SizeGB
					}
					break
				}
			}
		}
	}

	// link NAT gateways to subnets and elastic IPs
	for i := range ec2.NATGateways {
		for j := range ec2.Subnets {
			if ec2.Subnets[j].ID == ec2.NATGateways[i].SubnetID.Value() {
				ec2.Subnets[j].Relationships.NATGateways = append(ec2.Subnets[j].Relationships.NATGateways, &ec2.NATGateways[i])
				ec2.NATGateways[i].Relationships.Subnet = &ec2.Subnets[j]
				break
			}
		}
		for j := range ec2.ElasticIPs {
			if ec2.ElasticIPs[j].ID == ec2.NATGateways[i].AllocationID.Value() {
				ec2.ElasticIPs[j].Relationships.NATGateway = &ec2.NATGateways[i]
				ec2.NATGateways[i].Relationships.AllocatedElasticIP = &ec2.ElasticIPs[j]
				break
			}
		}
	}

	// link elastic IP associations to elastic IPs
	for i := range ec2.ElasticIPAssociations {
		for j := range ec2.ElasticIPs {
			if ec2.ElasticIPs[j].ID == ec2.ElasticIPAssociations[i].AllocationID.Value() {
				ec2.ElasticIPs[j].Relationships.Association = &ec2.ElasticIPAssociations[i]
				ec2.ElasticIPAssociations[i].Relationships.ElasticIP = &ec2.ElasticIPs[j]
				break
			}
		}
	}

	// link transit gateway VPC attachments to VPCs and transit gateways
	for i, attachment := range ec2.TransitGatewayVPCAttachments {
		for j := range ec2.VPCs {
			if ec2.VPCs[j].ID == attachment.VPCID.Value() {
				ec2.TransitGatewayVPCAttachments[i].Relationships.VPC = &ec2.VPCs[j]
				break
			}
		}
		for j := range ec2.TransitGateways {
			if ec2.TransitGateways[j].ID == attachment.TransitGatewayID.Value() {
				ec2.TransitGatewayVPCAttachments[i].Relationships.TransitGateway = &ec2.TransitGateways[j]
				break
			}
		}
	}

	// link transit gateway peering attachments to transit gateways
	for i, attachment := range ec2.TransitGatewayPeeringAttachments {
		for j := range ec2.TransitGateways {
			if ec2.TransitGateways[j].ID == attachment.TransitGatewayID.Value() {
				ec2.TransitGatewayPeeringAttachments[i].Relationships.TransitGateway = &ec2.TransitGateways[j]
				break
			}
		}
	}

	// link load balancer subnet mappings to elastic IPs and subnets
	for i := range ec2.LoadBalancers {
		for j, mapping := range ec2.LoadBalancers[i].SubnetMappings {
			for k := range ec2.ElasticIPs {
				if ec2.ElasticIPs[k].ID == mapping.AllocationID.Value() {
					ec2.ElasticIPs[k].Relationships.LoadBalancer = &ec2.LoadBalancers[i]
					ec2.LoadBalancers[i].SubnetMappings[j].Relationships.Allocation = &ec2.ElasticIPs[k]
					break
				}
			}
			for k := range ec2.Subnets {
				if ec2.Subnets[k].ID == mapping.SubnetID.Value() {
					ec2.LoadBalancers[i].SubnetMappings[j].Relationships.Subnet = &ec2.Subnets[k]
					break
				}
			}
		}
	}
}
