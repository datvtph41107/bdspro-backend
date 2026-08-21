package trust

import (
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestMetadataContractCopiesAndNormalizesInputs(t *testing.T) {
	privileged := []string{" ProfileID "}
	signed := []string{" X-Request-ID ", "profileid"}
	unsigned := []string{" /grpc.health.v1.Health/Check "}

	contract := NewMetadataContract(privileged, signed, unsigned)
	privileged[0] = "other"
	signed[0] = "other"
	unsigned[0] = "/other"

	if !HasPrivilegedMetadata(metadata.Pairs("profileid", "42"), contract) {
		t.Fatal("contract did not retain normalized privileged key")
	}
	if !contract.allowsUnsignedMethod("/grpc.health.v1.Health/Check") {
		t.Fatal("contract did not retain unsigned method")
	}
}

func TestTrustAssertionFieldsAreAlwaysSigned(t *testing.T) {
	contract := NewMetadataContract(nil, nil, nil)
	for _, key := range assertionMetadataKeys {
		found := false
		for _, signed := range contract.signedKeys {
			if signed == key {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("assertion key %q is not signed", key)
		}
	}
}
