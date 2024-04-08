//go:build generate

//go:generate go run build.go
//go:generate gofmt -w ./internal

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"github.com/caarlos0/log"
	"github.com/dyammarcano/version"
)

func main() {
	log.Infof("Applying gofmt to internal directory")

	ver, err := version.NewVersion()
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	if err = ver.Generate(); err != nil {
		log.Errorf(err.Error())
		return
	}
}
