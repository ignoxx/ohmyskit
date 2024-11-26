package handlers

import (
	"AABBCCDD/app/types"

	"github.com/ignoxx/ohmyskit/kit"
)

func HandleAuthentication(kit *kit.Kit) (kit.Auth, error) {
	return types.AuthUser{}, nil
}
