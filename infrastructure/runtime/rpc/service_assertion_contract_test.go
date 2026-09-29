package rpc

import (
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestMetadataContractCopiesAndNormalizesInputs(t *testing.T) {
	privileged := []string{" ProfileID "}
	signed := []string{" X-Request-ID ", "profileid"}
	unsigned := []string{" /grpc.health.v1.Health/Check "}

	contract := newMetadataContract(
		privileged,
		signed,
		unsigned,
	)

	privileged[0] = "other"
	signed[0] = "other"
	unsigned[0] = "/other"

	if !contract.hasPrivilegedMetadata(
		metadata.Pairs("profileid", "42"),
	) {
		t.Fatal("contract did not retain normalized privileged key")
	}

	if !contract.allowsUnsigned(
		"/grpc.health.v1.Health/Check",
	) {
		t.Fatal("contract did not retain unsigned method")
	}
}

func TestServiceAssertionFieldsAreAlwaysSigned(t *testing.T) {
	contract := newMetadataContract(nil, nil, nil)

	for _, key := range serviceAssertionMetadataKeys() {
		found := false
		for _, signed := range contract.signedKeys {
			if signed == key {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf(
				"assertion key %q is not signed",
				key,
			)
		}
	}
}
