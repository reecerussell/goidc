package dynamo

import (
	"fmt"
	"os"
)

func ClientsTableName() string {
	return getVariableOrPanic("CLIENTS_TABLE_NAME")
}

func UsersTableName() string {
	return getVariableOrPanic("USERS_TABLE_NAME")
}

func getVariableOrPanic(name string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	panic(fmt.Sprintf("the environment variable '%s' has not been configured", name))
}
