package fs

// Operation -
type Operation struct {
	Error chan error
}

// MoveOperation - Move source object to target object. Copy source to target, delete the source.
type MoveOperation struct {
	*Operation

	Source string
	Target string
}

func newMoveOp(sourcePath, targetPath string) MoveOperation {
	return MoveOperation{
		Source: sourcePath,
		Target: targetPath,
		Operation: &Operation{
			Error: make(chan error),
		},
	}
}

// CopyOperation - Copy source object to target.
type CopyOperation struct {
	*Operation

	Source string
	Target string
}

// PutOperation - Copy source file to target.
type PutOperation struct {
	*Operation

	Length int64

	Source string
	Target string
}

func newPutOp(sourcePath string, targetPath string, length int64) PutOperation {
	return PutOperation{
		Source: sourcePath,
		Target: targetPath,
		Length: int64(length),
		Operation: &Operation{
			Error: make(chan error),
		},
	}
}
