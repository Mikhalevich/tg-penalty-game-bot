package shotimageprovider

import (
	"embed"
)

//go:embed assets/*
var assetsFS embed.FS

type ShotImageProvider struct {
}

func New() *ShotImageProvider {
	return &ShotImageProvider{}
}
