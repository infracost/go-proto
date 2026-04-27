package efs

type EFS struct {
	FileSystems []FileSystem `tree:"file_systems"`
}
