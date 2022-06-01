package icon

import (
	"runtime"

	_ "embed"
)

//go:embed icon.ico
var iconWindows []byte

//go:embed icon.png
var icon []byte

func GetIcon() []byte {
	if runtime.GOOS == "windows" {
		return iconWindows
	}
	return icon
}
