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
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range c.ManagedDisks {
		c.ManagedDisks[i].Relationships = ManagedDiskRelationships{}
	}

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

	// Snapshots that don't set disk_size_gb directly inherit the size of the
	// managed disk identified by source_disk_id. Matches the legacy
	// internal/tree/azure/compute/service.go behaviour and drives the
	// snapshot storage cost.
	for i := range c.Snapshots {
		snap := &c.Snapshots[i]
		if snap.DiskSizeGB.Value() != 0 {
			continue
		}
		if disk, ok := diskByID[snap.SourceDiskID.Value()]; ok {
			snap.DiskSizeGB = snap.DiskSizeGB.WithValue(disk.DiskSizeGB.Value())
		}
	}

	// Images derive storage from either a source VM or from referenced
	// managed_disks (when os_disk.managed_disk_id / data_disk.managed_disk_id
	// are used instead of explicit size_gb). Resolve those references here so
	// the pricing layer can read OSDiskStorageGB + DataDisks[].SizeGB directly.
	vmByID := map[string]*VirtualMachine{}
	for i := range c.VirtualMachines {
		vmByID[c.VirtualMachines[i].ID] = &c.VirtualMachines[i]
	}
	for i := range c.Images {
		img := &c.Images[i]

		// Source VM: sum the VM's OS + data disks. Each entry checks the
		// resource's own DiskSizeGB first, then falls back to the linked
		// managed_disk's DiskSizeGB.
		if srcID := img.SourceVirtualMachineID.Value(); srcID != "" {
			if vm, ok := vmByID[srcID]; ok {
				img.OSDiskStorageGB = img.OSDiskStorageGB.WithValue(virtualMachineOSDiskSize(vm, diskByID))
				img.DataDisks = nil
				for _, dd := range virtualMachineDataDiskSizes(vm, diskByID) {
					img.DataDisks = append(img.DataDisks, ImageDataDisk{
						SizeGB: img.OSDiskStorageGB.WithValue(dd),
					})
				}
			}
			continue
		}

		// Direct disk references on the image itself.
		if img.OSDiskStorageGB.Value() == 0 && img.OSDiskID.Value() != "" {
			if disk, ok := diskByID[img.OSDiskID.Value()]; ok {
				img.OSDiskStorageGB = img.OSDiskStorageGB.WithValue(float64(disk.DiskSizeGB.Value()))
			}
		}
		for j := range img.DataDisks {
			dd := &img.DataDisks[j]
			if dd.SizeGB.Value() == 0 && dd.ID.Value() != "" {
				if disk, ok := diskByID[dd.ID.Value()]; ok {
					dd.SizeGB = dd.SizeGB.WithValue(float64(disk.DiskSizeGB.Value()))
				}
			}
		}
	}
}

// virtualMachineOSDiskSize returns the OS disk size for a VirtualMachine,
// preferring the inline disk_size_gb on storage_os_disk/os_disk over the
// referenced managed_disk's DiskSizeGB. Falls back to 128 (Azure's documented
// default for unspecified OS disks) only when the VM has no managed_disk
// reference either — matches the legacy tree's heuristic.
func virtualMachineOSDiskSize(vm *VirtualMachine, diskByID map[string]*ManagedDisk) float64 {
	if vm.StorageOSDisk != nil {
		if v := vm.StorageOSDisk.DiskSizeGB.Value(); v > 0 {
			return float64(v)
		}
		if disk, ok := diskByID[vm.StorageOSDisk.ManagedDiskID.Value()]; ok && disk.DiskSizeGB.Value() > 0 {
			return float64(disk.DiskSizeGB.Value())
		}
	}
	if vm.OSDisk != nil {
		if v := vm.OSDisk.DiskSizeGB.Value(); v > 0 {
			return float64(v)
		}
		if disk, ok := diskByID[vm.OSDisk.ManagedDiskID.Value()]; ok && disk.DiskSizeGB.Value() > 0 {
			return float64(disk.DiskSizeGB.Value())
		}
	}
	return 128
}

// virtualMachineDataDiskSizes returns the size of each data disk attached to
// the given VM (preferring inline size; falling back to referenced
// managed_disk). Zero-size entries are dropped.
func virtualMachineDataDiskSizes(vm *VirtualMachine, diskByID map[string]*ManagedDisk) []float64 {
	var out []float64
	for _, dd := range vm.StorageDataDisks {
		if dd == nil {
			continue
		}
		size := float64(dd.DiskSizeGB.Value())
		if size == 0 {
			if disk, ok := diskByID[dd.ManagedDiskID.Value()]; ok {
				size = float64(disk.DiskSizeGB.Value())
			}
		}
		if size > 0 {
			out = append(out, size)
		}
	}
	return out
}
