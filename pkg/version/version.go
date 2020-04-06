// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package version

import (
	"fmt"
	"runtime"
)

// Version is set with build flags
var Version = "unversioned"

// TemplatedVersion _
var TemplatedVersion = fmt.Sprintf("Patrasche version %s, %s %s %s\n", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())

// Print _
func Print() {
	fmt.Println(TemplatedVersion)
}
