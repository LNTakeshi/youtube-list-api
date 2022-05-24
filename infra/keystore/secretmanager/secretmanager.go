package secretmanager

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"youtubelist/application/keystore"
	"youtubelist/domain/config"
	"youtubelist/domain/config/constant"
)

type Client struct {
	cli *secretmanager.Client
}

func NewClient(ctx context.Context) (keystore.IKeyStore, error) {
	cli, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{
		cli: cli,
	}, nil
}

func (c Client) GetSecret(ctx context.Context, secretName constant.SecretKeyName) ([]byte, error) {
	secret, err := c.cli.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", config.ProjectID, string(secretName)),
	})
	if err != nil {
		return []byte{}, err
	}

	return secret.GetPayload().GetData(), nil
}
