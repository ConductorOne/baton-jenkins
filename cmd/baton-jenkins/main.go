package main

import (
	"context"
	"fmt"
	"os"

	cfg "github.com/conductorone/baton-jenkins/pkg/config"
	"github.com/conductorone/baton-jenkins/pkg/client"
	"github.com/conductorone/baton-jenkins/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"baton-jenkins",
		getConnector,
		cfg.Config,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilder(&connector.Connector{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, c *cfg.Jenkins) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)
	jenkinsClient := client.NewClient()
	if c.Token != "" {
		jenkinsClient.WithUser(c.Username).WithBearerToken(c.Token)
	}

	if c.Username != "" && c.Password != "" {
		jenkinsClient.WithUser(c.Username).WithPassword(c.Password)
	}

	cb, err := connector.New(ctx, c.BaseUrl, jenkinsClient)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	conn, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	return conn, nil
}
