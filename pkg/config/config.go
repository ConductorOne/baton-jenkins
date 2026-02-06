package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	UsernameField = field.StringField(
		"username",
		field.WithDescription("Username of administrator used to connect to the Jenkins API"),
		field.WithRequired(true),
	)
	PasswordField = field.StringField(
		"password",
		field.WithDescription("Application password used to connect to the Jenkins API"),
	)
	BaseUrlField = field.StringField(
		"base-url",
		field.WithDescription("Jenkins"),
		field.WithDefaultValue("http://localhost:8080"),
		field.WithRequired(true),
	)
	TokenField = field.StringField(
		"token",
		field.WithDescription("HTTP access tokens in Jenkins"),
	)

	ConfigurationFields = []field.SchemaField{
		UsernameField,
		PasswordField,
		BaseUrlField,
		TokenField,
	}

	ConfigRelationships = []field.SchemaFieldRelationship{
		field.FieldsMutuallyExclusive(TokenField, PasswordField),
		field.FieldsAtLeastOneUsed(TokenField, PasswordField),
	}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(ConfigurationFields, field.WithConstraints(ConfigRelationships...))
