package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"link-service/internal/service"
)

type processLinksRequest struct {
	Links []string `json:"links"`
}

func ProcessLinks(srv *service.Service, requestTimeout time.Duration, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		_ = ctx

		var reqLinks processLinksRequest
		err := json.NewDecoder(r.Body).Decode(&reqLinks)
		if err != nil {
			http.Error(w, "cannot decode body", http.StatusBadRequest)
			logger.Warn("ProcessLinks: cannot decode body", zap.Error(err))
			return
		}

		err = srv.Process(reqLinks.Links)
		if err != nil {
			http.Error(w, "failed to process links", http.StatusInternalServerError)
			logger.Error("ProcessLinks: failed to processing links", zap.Error(err))
			return
		}

		logger.Info("ProcessLinks: successfully process links")
		// TODO: write success response
	}
}

//pdf := gofpdf.New("P", "mm", "A4", "")
//pdf.AddPage()
//pdf.SetFont("Arial", "", 14)
//
//for _, recID := range ids {
//rec := storage.Get(recID)
//for link, status := range rec.Results {
//pdf.Cell(0, 10, fmt.Sprintf("%s: %s", link, status))
//pdf.Ln(10)
//}
//}
//
//pdf.Output(w)
