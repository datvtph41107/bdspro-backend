package reference

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

type Generator struct{}

func New() *Generator { return &Generator{} }
func (*Generator) NewOrderReference() (string, error) {
	var b [8]byte
	if _, e := rand.Read(b[:]); e != nil {
		return "", fmt.Errorf("generate order reference: %w", e)
	}
	return "QHP-" + strings.ToUpper(hex.EncodeToString(b[:])), nil
}
