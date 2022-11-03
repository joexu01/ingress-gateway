package cache

import (
	"fmt"
	"log"
	"testing"
)

func TestMemoryCacheService_Validate(t *testing.T) {
	cacheService := NewMemoryCacheService(300)

	err := cacheService.Validate("token 188")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cacheService.TokenStatusMap)
	fmt.Printf("%+v\n", *cacheService.TokenStatusMap["token 188"])
}

func TestMemoryCacheService_Revoke(t *testing.T) {
	cacheService := NewMemoryCacheService(300)

	err := cacheService.Validate("token 188")
	if err != nil {
		log.Fatal(err)
	}

	err = cacheService.Validate("token 189")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Before: %+v\n", *cacheService.TokenStatusMap["token 188"])

	_ = cacheService.Revoke("token 188")

	fmt.Printf("After: %+v\n", cacheService.TokenStatusMap["token 188"])
}

func TestMemoryCacheService_Verify(t *testing.T) {
	cacheService := NewMemoryCacheService(300)

	err := cacheService.Validate("token 188")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", *cacheService.TokenStatusMap["token 188"])

	result := cacheService.Verify("token 188")

	fmt.Printf("Token 188: %+v\n", result)

	if !result {
		log.Fatalln("result should be true")
	}
}
