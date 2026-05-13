package compute

type Compute struct {
	ManagedDisks                   []ManagedDisk                   `tree:"managed_disks"`
	VirtualMachines                []VirtualMachine                `tree:"virtual_machines"`
	WindowsVirtualMachines         []WindowsVirtualMachine         `tree:"windows_virtual_machines"`
	LinuxVirtualMachines           []LinuxVirtualMachine           `tree:"linux_virtual_machines"`
	Images                         []Image                         `tree:"images"`
	Snapshots                      []Snapshot                      `tree:"snapshots"`
	VirtualMachineScaleSets        []VirtualMachineScaleSet        `tree:"virtual_machine_scale_sets"`
	LinuxVirtualMachineScaleSets   []LinuxVirtualMachineScaleSet   `tree:"linux_virtual_machine_scale_sets"`
	WindowsVirtualMachineScaleSets []WindowsVirtualMachineScaleSet `tree:"windows_virtual_machine_scale_sets"`
	BatchPools                     []BatchPool                     `tree:"batch_pools"`
	DiskAttachments                []DiskAttachment                `tree:"disk_attachments"`
}

// PostProcess links each managed disk to the VM it is attached to via
// DiskAttachment entries or via the VM's inline OS/data disk references.
func (c *Compute) PostProcess() {
	diskByID := make(map[string]*ManagedDisk)
	for i := range c.ManagedDisks {
		if id := c.ManagedDisks[i].ID; id != "" {
			diskByID[id] = &c.ManagedDisks[i]
		}
	}

	link := func(managedDiskID string, set func(*ManagedDisk)) {
		if managedDiskID == "" {
			return
		}
		if disk, ok := diskByID[managedDiskID]; ok {
			set(disk)
		}
	}

	for _, att := range c.DiskAttachments {
		vmID := att.VirtualMachineID.Value()
		link(att.ManagedDiskID.Value(), func(disk *ManagedDisk) {
			for i := range c.VirtualMachines {
				if c.VirtualMachines[i].ID == vmID {
					disk.Relationships.VirtualMachine = &c.VirtualMachines[i]
					return
				}
			}
			for i := range c.LinuxVirtualMachines {
				if c.LinuxVirtualMachines[i].ID == vmID {
					disk.Relationships.LinuxVirtualMachine = &c.LinuxVirtualMachines[i]
					return
				}
			}
			for i := range c.WindowsVirtualMachines {
				if c.WindowsVirtualMachines[i].ID == vmID {
					disk.Relationships.WindowsVirtualMachine = &c.WindowsVirtualMachines[i]
					return
				}
			}
		})
	}

	for i := range c.VirtualMachines {
		vm := &c.VirtualMachines[i]
		setVM := func(d *ManagedDisk) { d.Relationships.VirtualMachine = vm }
		if vm.OSDisk != nil {
			link(vm.OSDisk.ManagedDiskID.Value(), setVM)
		}
		if vm.StorageOSDisk != nil {
			link(vm.StorageOSDisk.ManagedDiskID.Value(), setVM)
		}
		for _, dd := range vm.StorageDataDisks {
			if dd == nil {
				continue
			}
			link(dd.ManagedDiskID.Value(), setVM)
		}
	}

	for i := range c.LinuxVirtualMachines {
		vm := &c.LinuxVirtualMachines[i]
		if vm.OSDisk != nil {
			link(vm.OSDisk.ManagedDiskID.Value(), func(d *ManagedDisk) {
				d.Relationships.LinuxVirtualMachine = vm
			})
		}
	}

	for i := range c.WindowsVirtualMachines {
		vm := &c.WindowsVirtualMachines[i]
		if vm.OSDisk != nil {
			link(vm.OSDisk.ManagedDiskID.Value(), func(d *ManagedDisk) {
				d.Relationships.WindowsVirtualMachine = vm
			})
		}
	}
}
