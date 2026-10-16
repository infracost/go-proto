package directoryservice

type DirectoryService struct {
	Directories       []Directory       `tree:"directories"`
	Regions           []Region          `tree:"regions"`
	SharedDirectories []SharedDirectory `tree:"shared_directories"`
}

func (d *DirectoryService) PostProcess() {
	for i := range d.SharedDirectories {
		d.SharedDirectories[i].Relationships = SharedDirectoryRelationships{}
		for j := range d.Directories {
			id := d.SharedDirectories[i].DirectoryID.Value()
			if d.Directories[j].ID == id || (id != "" && d.Directories[j].Name.Value() == id) {
				d.SharedDirectories[i].Relationships.Directory = &d.Directories[j]
				break
			}
		}
	}
}
