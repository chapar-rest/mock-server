package api

import (
	"fmt"
	"net/http"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
)

type sampleDocument struct {
	contentType string
	body        string
}

// sampleDocuments hold the same small todo list in every supported format, so
// a client's response viewers can be compared side by side.
var sampleDocuments = map[restapi.GetFormatParamsFormat]sampleDocument{
	restapi.GetFormatParamsFormatJson: {
		contentType: "application/json",
		body: `{"todos":[{"id":1,"title":"Try Chapar","completed":true},{"id":2,"title":"Call the gRPC API","completed":false}]}
`,
	},
	restapi.GetFormatParamsFormatXml: {
		contentType: "application/xml",
		body: `<?xml version="1.0" encoding="UTF-8"?>
<todos>
  <todo id="1" completed="true">Try Chapar</todo>
  <todo id="2" completed="false">Call the gRPC API</todo>
</todos>
`,
	},
	restapi.GetFormatParamsFormatHtml: {
		contentType: "text/html; charset=utf-8",
		body: `<!doctype html>
<html>
  <head><title>Todos</title></head>
  <body>
    <ul>
      <li><s>Try Chapar</s></li>
      <li>Call the gRPC API</li>
    </ul>
  </body>
</html>
`,
	},
	restapi.GetFormatParamsFormatText: {
		contentType: "text/plain; charset=utf-8",
		body: `[x] Try Chapar
[ ] Call the gRPC API
`,
	},
	restapi.GetFormatParamsFormatCsv: {
		contentType: "text/csv; charset=utf-8",
		body: `id,title,completed
1,Try Chapar,true
2,Call the gRPC API,false
`,
	},
}

// GetFormat implements restapi.ServerInterface.
func (c *Controller) GetFormat(w http.ResponseWriter, r *http.Request, format restapi.GetFormatParamsFormat) {
	doc, ok := sampleDocuments[format]
	if !ok {
		c.writeError(w, r, fmt.Errorf("%w: format must be one of json, xml, html, text, csv", errx.ErrInvalidArgument))
		return
	}

	w.Header().Set("Content-Type", doc.contentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(doc.body))
}
