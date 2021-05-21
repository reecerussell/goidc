package dynamo

import (
	"context"

	"github.com/reecerussell/goidc"
)

func ClientsTableName(ctx context.Context) string {
	return goidc.StageVariable(ctx, "CLIENTS_TABLE_NAME")
}

func UsersTableName(ctx context.Context) string {
	return goidc.StageVariable(ctx, "USERS_TABLE_NAME")
}
