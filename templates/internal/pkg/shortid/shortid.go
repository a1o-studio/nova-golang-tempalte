package shortid

import (
	"fmt"

	"github.com/sqids/sqids-go"
)

const (
	alphabet = "jzqvwnbtkfylgphdsoxumcirae" //	@name	charset
	base     = int64(len(alphabet))
	length   = 6                 // 固定
	offset   = int64(10_000_000) // 偏移量，确保生成的 ID 不会太小
)

var sq *sqids.Sqids

func init() {
	sq, _ = sqids.New(sqids.Options{
		MinLength: length,
		Alphabet:  alphabet,
	})
}

func Encode(id int64) (string, error) {
	return sq.Encode([]uint64{uint64(id)})
}

// Decode 还原为原始 int64
func Decode(s string) (int64, error) {
	ids := sq.Decode(s)
	if len(ids) == 0 {
		return 0, fmt.Errorf("invalid short ID: %s", s)
	}
	return int64(ids[0]), nil
}
