package gfs

import (
	"context"
	"nav-green-download/pkg/server"
)

func ExecuteScript(startStr, endStr string) error {
	gfs := server.NewGFSDownloader()
	defer gfs.Stop(context.TODO())

	gfs.Download(context.TODO())

	return nil
}
