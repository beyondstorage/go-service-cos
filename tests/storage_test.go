package tests

import (
	"os"
	"testing"

	tests "github.com/beyondstorage/go-integration-test/v4"
)

func TestStorage(t *testing.T) {
	if os.Getenv("STORAGE_COS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_COS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorager(t, setupTest(t))
}

// FIXME: For `CompleteMultipartUpload`, the numbers of the uploaded parts must be continuous and the part information entries in the request body must be sorted by number in ascending order
// ref: https://cloud.tencent.com/document/product/436/7742
//func TestMultiparter(t *testing.T) {
//	if os.Getenv("STORAGE_COS_INTEGRATION_TEST") != "on" {
//		t.Skipf("STORAGE_COS_INTEGRATION_TEST is not 'on', skipped")
//	}
//	tests.TestMultiparter(t, setupTest(t))
//}

func TestDir(t *testing.T) {
	if os.Getenv("STORAGE_COS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_COS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestDirer(t, setupTest(t))
}
