package common

import "path/filepath"

func GetFileFormatFromName(name string) string {
	return filepath.Ext(name)
}
