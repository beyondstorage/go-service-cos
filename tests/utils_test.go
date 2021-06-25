package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	cos "github.com/beyondstorage/go-service-cos/v2"
	ps "github.com/beyondstorage/go-storage/v4/pairs"
	"github.com/beyondstorage/go-storage/v4/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for oss")

	store, err := cos.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_COS_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_COS_NAME")),
		ps.WithLocation(os.Getenv("STORAGE_COS_LOCATION")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		cos.WithStorageFeatures(cos.StorageFeatures{
			VirtualDir: true,
		}),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
