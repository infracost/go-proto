package fsx

type FSx struct {
	FileSystems []FileSystem `tree:"file_systems"`
	FileCaches  []FileCache  `tree:"file_caches"`
}
