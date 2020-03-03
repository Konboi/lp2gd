// refs: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// +build tools
package tools

import (
	_ "github.com/a8m/envsubst/cmd/envsubst"
)
