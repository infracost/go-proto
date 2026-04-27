package fsx

type FSx struct {
	FileSystems []FileSystem `tree:"file_systems"`
}
