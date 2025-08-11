package _jwt

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc/metadata"
)

func MetadataFromContext(ctx context.Context) *UserMetadata {
	// Lấy metadata từ context
	md, ok := metadata.FromIncomingContext(ctx)
	metadata := UserMetadata{}
	if !ok {
		fmt.Println("No metadata found in context")
		return &metadata
	}

	// In ra các key-value trong metadata
	for key, values := range md {
		if key == "profileid" {
			if len(values) > 0 {
				profileID, err := strconv.ParseUint(values[0], 10, 64)
				if err != nil {
					fmt.Println("Error parsing profileId:", err)
					return nil
				}
				metadata.ProfileID = &profileID
			}
		}
		if key == "organizationid" {
			if len(values) > 0 {
				organizationID, err := strconv.ParseUint(values[0], 10, 64)
				if err != nil {
					fmt.Println("Error parsing profileId:", err)
					return nil
				}
				metadata.OrganizationID = &organizationID
			}
		}
		// Nếu tìm thấy "role", lưu giá trị role
		if key == "role" {
			metadata.Role = values
		}
	}

	return &metadata
}
