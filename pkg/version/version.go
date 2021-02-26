// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package version

import (
	"fmt"
	"runtime"
)

// Version is set manually (Makefile)
var Version = "v0.8.0"

// TemplatedVersion _
var TemplatedVersion = fmt.Sprintf("Patrasche %s, %s %s %s", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
