// +build integration_test

package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	cos "github.com/aos-dev/go-service-cos"
	ps "github.com/aos-dev/go-storage/v3/pairs"
	"github.com/aos-dev/go-storage/v3/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for oss")

	store, err := cos.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_COS_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_COS_NAME")),
		ps.WithLocation(os.Getenv("STORAGE_COS_LOCATION")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
