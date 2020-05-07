// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package version

import (
	"fmt"
	"runtime"
)

// Version is set manually (Makefile)
var Version = "v0.1.0"

// TemplatedVersion _
var TemplatedVersion = fmt.Sprintf("Patrasche version %s, %s %s %s\n", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
