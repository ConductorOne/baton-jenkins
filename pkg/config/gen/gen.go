package main

import (
	cfg "github.com/conductorone/baton-jenkins/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("jenkins", cfg.Config)
}
