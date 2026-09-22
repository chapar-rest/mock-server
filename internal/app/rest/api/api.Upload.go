package api

import (
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"slices"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// maxUploadMemory is how much of a multipart body is held in memory; the
// rest spills to temporary files. The router's body limit caps the total.
const maxUploadMemory = 8 << 20

// Upload implements restapi.ServerInterface.
func (c *Controller) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
		c.writeError(w, r, fmt.Errorf("%w: expected a multipart/form-data body: %w", errx.ErrInvalidArgument, err))
		return
	}
	defer func() { _ = r.MultipartForm.RemoveAll() }()

	files := make([]restapi.UploadedFile, 0)
	for field, headers := range r.MultipartForm.File {
		for _, header := range headers {
			file, err := describeUpload(field, header)
			if err != nil {
				c.writeError(w, r, err)
				return
			}
			files = append(files, file)
		}
	}
	slices.SortFunc(files, func(a, b restapi.UploadedFile) int {
		return cmp.Or(cmp.Compare(a.Field, b.Field), cmp.Compare(a.Filename, b.Filename))
	})

	restutils.Success(w, http.StatusOK, &restapi.UploadResponse{
		Files:  files,
		Fields: r.MultipartForm.Value,
	})
}

func describeUpload(field string, header *multipart.FileHeader) (restapi.UploadedFile, error) {
	f, err := header.Open()
	if err != nil {
		return restapi.UploadedFile{}, fmt.Errorf("opening upload %q: %w", header.Filename, err) //nolint:exhaustruct // zero value on error.
	}
	defer func() { _ = f.Close() }()

	hash := sha256.New()
	size, err := io.Copy(hash, f)
	if err != nil {
		return restapi.UploadedFile{}, fmt.Errorf("reading upload %q: %w", header.Filename, err) //nolint:exhaustruct // zero value on error.
	}

	return restapi.UploadedFile{
		Field:       field,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        size,
		Sha256:      hex.EncodeToString(hash.Sum(nil)),
	}, nil
}
