package filehttp

import "file/dto"

// externalFileInfo preserves the existing JSON shape while preventing a
// deployment-specific filesystem path from crossing the File HTTP boundary.
func externalFileInfo(info *dto.FileInfo) *dto.FileInfo {
	if info == nil {
		return nil
	}
	out := *info
	out.AbsolutePath = ""
	return &out
}
