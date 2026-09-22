package api

import (
	"fmt"
	"mime"
	"net/http"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/utils/restutils"
)

// SubmitForm implements restapi.ServerInterface.
func (c *Controller) SubmitForm(w http.ResponseWriter, r *http.Request) {
	mediaType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mediaType != "application/x-www-form-urlencoded" {
		c.writeError(w, r, fmt.Errorf("%w: expected an application/x-www-form-urlencoded body", errx.ErrInvalidArgument))
		return
	}

	if err := r.ParseForm(); err != nil {
		c.writeError(w, r, fmt.Errorf("%w: invalid form body: %w", errx.ErrInvalidArgument, err))
		return
	}

	restutils.Success(w, http.StatusOK, &restapi.FormResponse{Fields: r.PostForm})
}
