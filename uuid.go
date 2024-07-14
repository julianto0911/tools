package tools

import (
	"github.com/lithammer/shortuuid/v4"
)

func ShortUUID() string {
	return shortuuid.New()
}
