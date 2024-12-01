package fileManager

type OutsideOfFileBoundsError struct{}

func (err OutsideOfFileBoundsError) Error() string {
	return "cannot read outside of file bounds"
}
