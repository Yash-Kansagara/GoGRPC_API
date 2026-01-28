package tokendb

import (
	"context"
	"log"

	"github.com/Yash-Kansagara/GoGRPC_API/pkg/utils"
	"github.com/allegro/bigcache/v3"
)

var Cache *bigcache.BigCache

func GetToken(token string) ([]byte, bool) {
	data, err := Cache.Get(token)
	if err != nil {
		log.Println("Failed to get token", token, err)
		return nil, false
	}
	return data, true
}

func AddToken(token string) error {
	log.Println("Adding token", token)
	return Cache.Set(token, []byte(token))
}

func RemoveToken(token string) error {
	return Cache.Delete(token)
}

func Init() {
	var err error
	duration := utils.GetDefaultRefreshTokenExpiry()
	Cache, err = bigcache.New(context.Background(), bigcache.DefaultConfig(duration))
	if err != nil {
		log.Fatal("Failed to initialize bigcache for refresh tokens")
	}
	log.Println("Bigcache initialized for refresh tokens")
}
