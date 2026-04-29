package web

import "embed"

//go:embed templates/*
var Assets embed.FS
