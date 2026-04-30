package compute

import "github.com/infracost/go-proto/pkg/tree/value"

type Compute struct {
	Addresses                []Address              `tree:"addresses"`
	Disks                    []Disk                 `tree:"disks"`
	RegionDisks              []RegionDisk           `tree:"region_disks"`
	ForwardingRules          []ForwardingRule       `tree:"forwarding_rules"`
	Snapshots                []Snapshot             `tree:"snapshots"`
	Images                   []Image                `tree:"images"`
	Instances                []Instance             `tree:"instances"`
	ExternalVPNGateways      []ExternalVPNGateway   `tree:"external_vpn_gateways"`
	MachineImages            []MachineImage         `tree:"machine_images"`
	InstanceGroupManagers    []InstanceGroupManager `tree:"instance_group_managers"`
	PerInstanceConfigs       []PerInstanceConfig    `tree:"per_instance_configs"`
	InstanceTemplates        []InstanceTemplate     `tree:"instance_templates"`
	RouterNATs               []RouterNAT            `tree:"router_nats"`
	TargetGRPCProxies        []TargetGRPCProxy      `tree:"target_grpc_proxies"`
	TargetHTTPProxies        []TargetGRPCProxy      `tree:"target_http_proxies"`
	TargetHTTPSProxies       []TargetGRPCProxy      `tree:"target_https_proxies"`
	TargetSSLProxies         []TargetGRPCProxy      `tree:"target_ssl_proxies"`
	TargetTCPProxies         []TargetGRPCProxy      `tree:"target_tcp_proxies"`
	RegionTargetHTTPProxies  []TargetGRPCProxy      `tree:"region_target_http_proxies"`
	RegionTargetHTTPSProxies []TargetGRPCProxy      `tree:"region_target_https_proxies"`
	HAVPNGateways            []HAVPNGateway         `tree:"ha_vpn_gateways"`
	VPNTunnels               []VPNTunnel            `tree:"vpn_tunnels"`
	VPNGateways              []VPNGateway           `tree:"vpn_gateways"`
	NodeTemplates            []NodeTemplate         `tree:"node_templates"`
	DataprocClusters         []DataprocCluster      `tree:"dataproc_clusters"`
	DataflowJobs             []DataflowJob          `tree:"dataflow_jobs"`
	RegionInstanceTemplates  []RegionInstanceTemplate `tree:"region_instance_templates"`
	ResourcePolicies         []ResourcePolicy       `tree:"resource_policies"`
}

func (c *Compute) PostProcess() {
	// Link InstanceGroupManagers to InstanceTemplates by TemplateRef matching template Name or ID
	for i, igm := range c.InstanceGroupManagers {
		if igm.TemplateRef.IsEmpty() {
			continue
		}
		for j := range c.InstanceTemplates {
			t := &c.InstanceTemplates[j]
			if igm.TemplateRef.Equal(t.ID) || igm.TemplateRef.Value() == t.Name.Value() {
				c.InstanceGroupManagers[i].Relationships.InstanceTemplate = t
				break
			}
		}
	}

	// Resolve disk storage from source images/snapshots
	for i := range c.Disks {
		if c.Disks[i].StorageGB.IsDefaultOrEmpty() {
			c.Disks[i].StorageGB = c.findDiskStorage(&c.Disks[i])
		}
	}

	// Resolve snapshot storage from source disks
	for i := range c.Snapshots {
		if c.Snapshots[i].StorageGB.IsDefaultOrEmpty() && !c.Snapshots[i].SourceDisk.IsEmpty() {
			for _, disk := range c.Disks {
				if c.Snapshots[i].SourceDisk.Value() == disk.ID ||
					c.Snapshots[i].SourceDisk.Value() == disk.Name.Value() ||
					c.Snapshots[i].SourceDisk.Value() == disk.SelfLink.Value() {
					c.Snapshots[i].StorageGB = disk.StorageGB
					break
				}
			}
		}
	}

	// Resolve image storage from source disks/images/snapshots
	for i := range c.Images {
		if c.Images[i].StorageGB.IsDefaultOrEmpty() {
			c.Images[i].StorageGB = c.findImageStorage(&c.Images[i])
		}
	}

	// Mark disks as attached when referenced by instances
	for _, inst := range c.Instances {
		if !inst.AttachedDisk.IsEmpty() {
			for j := range c.Disks {
				if inst.AttachedDisk.Value() == c.Disks[j].ID ||
					inst.AttachedDisk.Value() == c.Disks[j].Name.Value() ||
					inst.AttachedDisk.Value() == c.Disks[j].SelfLink.Value() {
					c.Disks[j].IsAttached = true
				}
			}
		}
	}
}

func (c *Compute) findDiskStorage(disk *Disk) value.Double {
	if !disk.SourceImage.IsEmpty() {
		for _, img := range c.Images {
			if disk.SourceImage.Value() == img.ID ||
				disk.SourceImage.Value() == img.Name.Value() {
				if !img.StorageGB.IsDefaultOrEmpty() {
					return img.StorageGB
				}
				return c.findImageStorage(&img)
			}
		}
	}
	if !disk.SourceSnapshot.IsEmpty() {
		for _, snap := range c.Snapshots {
			if disk.SourceSnapshot.Value() == snap.ID ||
				disk.SourceSnapshot.Value() == snap.Name.Value() {
				return snap.StorageGB
			}
		}
	}
	return disk.StorageGB
}

func (c *Compute) findImageStorage(img *Image) value.Double {
	if !img.SourceDisk.IsEmpty() {
		for _, disk := range c.Disks {
			if img.SourceDisk.Value() == disk.ID ||
				img.SourceDisk.Value() == disk.Name.Value() ||
				img.SourceDisk.Value() == disk.SelfLink.Value() {
				return disk.StorageGB
			}
		}
	}
	if !img.SourceImage.IsEmpty() {
		for _, src := range c.Images {
			if img.SourceImage.Value() == src.ID ||
				img.SourceImage.Value() == src.Name.Value() {
				if !src.StorageGB.IsDefaultOrEmpty() {
					return src.StorageGB
				}
			}
		}
	}
	if !img.SourceSnapshot.IsEmpty() {
		for _, snap := range c.Snapshots {
			if img.SourceSnapshot.Value() == snap.ID ||
				img.SourceSnapshot.Value() == snap.Name.Value() {
				return snap.StorageGB
			}
		}
	}
	return img.StorageGB
}
