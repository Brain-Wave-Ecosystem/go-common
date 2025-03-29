package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/rand"
	"google.golang.org/grpc/metadata"
)

func GenerateSlug(str string) (slug string) {
	slug = strings.ToLower(str)
	slug = strings.ReplaceAll(slug, " ", "-")

	rand.Seed(uint64(time.Now().UnixNano()))
	suffix := rand.Intn(9000) + 1000

	return fmt.Sprintf("%s-%d", slug, suffix)
}

func ExtractFromMD(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata not provided")
	}

	values := md.Get(key)
	if len(values) == 0 {
		return "", errors.New("key not found")
	}

	return values[0], nil
}
