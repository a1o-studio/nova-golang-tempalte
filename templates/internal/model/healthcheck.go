package model

import "github.com/a1ostudio/nova/internal/config"

type System struct {
	Environment config.Env `json:"environment" example:"prod"`
	Version     string     `json:"version" example:"1.0.0"`
} //	@name	System

type Healthcheck struct {
	Status string `json:"status" example:"available"`
	System System `json:"system"`
} //	@name	Healthcheck
