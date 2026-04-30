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
