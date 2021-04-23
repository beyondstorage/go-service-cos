package cos

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/aos-dev/go-storage/v3/pkg/headers"
	"github.com/aos-dev/go-storage/v3/pkg/iowrap"
	. "github.com/aos-dev/go-storage/v3/types"
)

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	o = s.newObject(false)
	o.Mode = ModeRead
	o.ID = s.getAbsPath(path)
	o.Path = path
	return o
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	_, err = s.object.Delete(ctx, rp)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	input := &objectPageStatus{
		maxKeys: 200,
		prefix:  s.getAbsPath(path),
	}

	var nextFn NextObjectFunc

	switch {
	case opt.ListMode.IsDir():
		input.delimiter = "/"
		nextFn = s.nextObjectPageByDir
	case opt.ListMode.IsPrefix():
		nextFn = s.nextObjectPageByPrefix
	default:
		return nil, fmt.Errorf("invalid list mode")
	}

	return NewObjectIterator(ctx, nextFn, input), nil
}

func (s *Storage) metadata(ctx context.Context, opt pairStorageMetadata) (meta *StorageMeta, err error) {
	meta = NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	return
}

func (s *Storage) nextObjectPageByDir(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	output, _, err := s.bucket.Get(ctx, &cos.BucketGetOptions{
		Prefix:    input.prefix,
		Delimiter: input.delimiter,
		Marker:    input.marker,
		MaxKeys:   input.maxKeys,
	})
	if err != nil {
		return err
	}

	for _, v := range output.CommonPrefixes {
		o := s.newObject(true)
		o.ID = v
		o.Path = s.getRelPath(v)
		o.Mode |= ModeDir

		page.Data = append(page.Data, o)
	}

	for _, v := range output.Contents {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if !output.IsTruncated {
		return IterateDone
	}

	input.marker = output.NextMarker
	return nil
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	output, _, err := s.bucket.Get(ctx, &cos.BucketGetOptions{
		Prefix:  input.prefix,
		Marker:  input.marker,
		MaxKeys: input.maxKeys,
	})
	if err != nil {
		return err
	}

	for _, v := range output.Contents {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if !output.IsTruncated {
		return IterateDone
	}

	input.marker = output.NextMarker
	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)

	getOptions := &cos.ObjectGetOptions{}
	// SSE-C
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		getOptions.XCosSSECustomerAglo, getOptions.XCosSSECustomerKey, getOptions.XCosSSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return 0, err
		}
	}
	resp, err := s.object.Get(ctx, rp, getOptions)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	rc := resp.Body
	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}

	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	rp := s.getAbsPath(path)

	headOptions := &cos.ObjectHeadOptions{}
	// SSE-C
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		headOptions.XCosSSECustomerAglo, headOptions.XCosSSECustomerKey, headOptions.XCosSSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return
		}
	}
	output, err := s.object.Head(ctx, rp, headOptions)
	if err != nil {
		return nil, err
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.Mode |= ModeRead

	o.SetContentLength(output.ContentLength)

	// COS uses RFC1123 format in HEAD
	//
	// > Last-Modified: Fri, 09 Aug 2019 10:20:56 GMT
	//
	// ref: https://cloud.tencent.com/document/product/436/7745
	if v := output.Header.Get(headers.LastModified); v != "" {
		lastModified, err := time.Parse(time.RFC1123, v)
		if err != nil {
			return nil, err
		}
		o.SetLastModified(lastModified)
	}

	if v := output.Header.Get(headers.ContentType); v != "" {
		o.SetContentType(v)
	}

	if v := output.Header.Get(headers.ETag); v != "" {
		o.SetEtag(v)
	}

	sm := make(map[string]string)
	if v := output.Header.Get(storageClassHeader); v != "" {
		sm[MetadataStorageClass] = v
	}
	if v := output.Header.Get(serverSideEncryptionHeader); v != "" {
		sm[MetadataServerSideEncryption] = v
	}
	if v := output.Header.Get(serverSideEncryptionCosKmsKeyIdHeader); v != "" {
		sm[MetadataServerSideEncryptionCosKmsKeyID] = v
	}
	if v := output.Header.Get(serverSideEncryptionCustomerAlgorithmHeader); v != "" {
		sm[MetadataServerSideEncryptionCustomerAlgorithm] = v
	}
	if v := output.Header.Get(serverSideEncryptionCustomerKeyMd5Header); v != "" {
		sm[MetadataServerSideEncryptionCustomerKeyMd5] = v
	}
	o.SetServiceMetadata(sm)

	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	rp := s.getAbsPath(path)

	putOptions := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: int(size),
		},
	}
	if opt.HasContentMd5 {
		putOptions.ContentMD5 = opt.ContentMd5
	}
	if opt.HasStorageClass {
		putOptions.XCosStorageClass = opt.StorageClass
	}
	// SSE-C
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		putOptions.XCosSSECustomerAglo, putOptions.XCosSSECustomerKey, putOptions.XCosSSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return
		}
	}
	// SSE-COS or SSE-KMS
	if opt.HasServerSideEncryption {
		putOptions.XCosServerSideEncryption = opt.ServerSideEncryption
		if putOptions.XCosServerSideEncryption == ServerSideEncryptionCosKms {
			putOptions.XOptionHeader = &http.Header{}
			if opt.HasServerSideEncryptionCosKmsKeyID {
				putOptions.XOptionHeader.Set("x-cos-server-side-encryption-cos-kms-key-id", opt.ServerSideEncryptionCosKmsKeyID)
			}
			if opt.HasServerSideEncryptionContext {
				putOptions.XOptionHeader.Set("x-cos-server-side-encryption-context", opt.ServerSideEncryptionContext)
			}
		}
	}
	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	_, err = s.object.Put(ctx, rp, r, putOptions)
	if err != nil {
		return 0, err
	}
	return
}
