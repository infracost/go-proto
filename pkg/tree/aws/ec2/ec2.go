package ec2

type EC2 struct {
	Instances       []Instance             `tree:"instance"`
	InstanceStates  []InstanceStateMapping `tree:"instance_state"`
	LaunchTemplates []LaunchTemplate       `tree:"launch_template"`
}

func (ec2 *EC2) PostProcess() {
	// reset instance relationships
	for i := range ec2.Instances {
		ec2.Instances[i].Relationships = InstanceRelationships{}
	}

	// link instance states to instances
	for i, instanceState := range ec2.InstanceStates {
		for j := range ec2.Instances {
			if instanceState.InstanceID.Equals(ec2.Instances[j].ID) {
				ec2.Instances[j].Relationships.InstanceState = &ec2.InstanceStates[i]
				break
			}
		}
	}

	// link launch templates to instances
	for i, instance := range ec2.Instances {
		for j := range ec2.LaunchTemplates {
			lt := &ec2.LaunchTemplates[j]
			if (instance.LaunchTemplateID.IsEmpty() || instance.LaunchTemplateID.Value() != lt.ID) &&
				(instance.LaunchTemplateName.IsEmpty() || instance.LaunchTemplateName.Value() != lt.Name.Value()) &&
				(instance.LaunchTemplateID.IsEmpty() || instance.LaunchTemplateID.Value() != lt.Name.Value()) {
				continue
			}

			// store the launch template on the instance - useful for policies to have this available
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

			launchTemplateBlockDeviceMap := make(map[string]BlockDeviceMapping)
			for _, blockDevice := range lt.BlockDeviceMappings {
				launchTemplateBlockDeviceMap[blockDevice.DeviceName.Value()] = blockDevice
			}

			instanceBlockDeviceMap := make(map[string]BlockDeviceMapping)
			for _, blockDevice := range instance.BlockDeviceMappings {
				instanceBlockDeviceMap[blockDevice.DeviceName.Value()] = blockDevice
			}

			blockDevices := []BlockDeviceMapping{}

			// merge any defaults from the launch template into existing block devices
			for name, blockDevice := range instanceBlockDeviceMap {
				if launchTemplateBlockDevice, ok := launchTemplateBlockDeviceMap[name]; ok {
					if blockDevice.EBSVolume.Type.IsDefaultOrEmpty() {
						blockDevice.EBSVolume.Type = launchTemplateBlockDevice.EBSVolume.Type
					}
					if blockDevice.EBSVolume.IOPS.IsDefaultOrEmpty() {
						blockDevice.EBSVolume.IOPS = launchTemplateBlockDevice.EBSVolume.IOPS
					}
					if blockDevice.EBSVolume.Size.IsDefaultOrEmpty() {
						blockDevice.EBSVolume.Size = launchTemplateBlockDevice.EBSVolume.Size
					}
					if blockDevice.EBSVolume.Throughput.IsDefaultOrEmpty() {
						blockDevice.EBSVolume.Throughput = launchTemplateBlockDevice.EBSVolume.Throughput
					}
				}
				blockDevices = append(blockDevices, blockDevice)
			}

			// add any new block devices from the launch template
			for name, blockDevice := range launchTemplateBlockDeviceMap {
				if _, ok := instanceBlockDeviceMap[name]; !ok {
					blockDevices = append(blockDevices, blockDevice)
				}
			}

			instance.BlockDeviceMappings = blockDevices

			// save changes
			ec2.Instances[i] = instance
			break
		}
	}
}
