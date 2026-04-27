package directoryservice

type DirectoryService struct {
	Directories []Directory `tree:"directories"`
}
