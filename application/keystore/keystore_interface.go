package keystore

import (
	"context"
	"youtubelist/domain/config/constant"
)

type IKeyStore interface {
	GetSecret(ctx context.Context, secretName constant.SecretKeyName) ([]byte, error)
}
