package helper

import "github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/models/config"

func GetUsersByChain(users []config.UserConfig, chain string) []config.UserConfig {
	var filteredUsers []config.UserConfig

	for _, user := range users {
		if user.Chain == chain {
			filteredUsers = append(filteredUsers, user)
		}
	}

	return filteredUsers
}
