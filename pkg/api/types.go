package api

type Config struct {
	Bucket    string
	Location  string
	AccessKey string
	SecretKey string
}

const (
	// GET operation
	OpsGetFileStatus = "GETFILESTATUS"
	OpsListStatus    = "LISTSTATUS"
	// Get Content Summary of a Directory
	OpsGetContentSummary = "GETCONTENTSUMMARY"
	// Get File Checksum
	OpsGetFileChecksum = "GETFILECHECKSUM"

	// DELETE operation
	OpsDelete = "DELETE"

	// PUT operation
	// Create and Write to a File
	OpsFileCreate = "CREATE"
	// Make a Directory
	OpsDirCreate = "MKDIRS"
	// Rename a File/Directory
	OpsRename = "RENAME"
	// Set Permission
	OpsSetPermission = "SETPERMISSION"
	// Set Access or Modification Time
	OpsSetTimes = "SETTIMES"
)
