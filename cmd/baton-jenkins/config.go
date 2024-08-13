package main

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	username = field.StringField("username", field.WithDescription("Username of administrator used to connect to the Jenkins API"), field.WithRequired(true))
	password = field.StringField("password", field.WithDescription("Application password used to connect to the Jenkins API"))
	baseUrl  = field.StringField("base-url", field.WithDescription("Jenkins"), field.WithDefaultValue("http://localhost:8080"), field.WithRequired(true))
	token    = field.StringField("token", field.WithDescription("HTTP access tokens in Jenkins"))
)

var relationships = []field.SchemaFieldRelationship{
	field.FieldsMutuallyExclusive(token, password),
}

var configuration = field.NewConfiguration([]field.SchemaField{username, password, baseUrl, token}, relationships...)
