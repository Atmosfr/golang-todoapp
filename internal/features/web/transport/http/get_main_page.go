package web_transport_http

import (
	"net/http"

	core_logger "github.com/Atmosfr/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Atmosfr/golang-todoapp/internal/core/transport/http/response"
)

func (h *WebHTTPHandler) GetMainPage(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	htmlFile, err := h.webService.GetMainPage()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get index.html for main page")
		return
	}

	responseHandler.HTMLResponse(htmlFile)
}
