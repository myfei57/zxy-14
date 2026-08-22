package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"maskhub/internal/store"
)

// API exposes the masking control plane over HTTP.
type API struct {
	state *store.State
}

// NewAPI builds the HTTP handler for the control plane.
func NewAPI(state *store.State) http.Handler {
	api := &API{state: state}
	r := chi.NewRouter()
	r.Get("/", api.IndexPage)
	r.Get("/console/sources", api.SourcesPage)
	r.Get("/console/policies", api.PoliciesPage)
	r.Get("/console/masking", api.MaskingPage)
	r.Get("/console/audit", api.AuditPage)

	r.Get("/api/stats", api.Stats)
	r.Get("/api/sources", api.ListSources)
	r.Post("/api/sources", api.CreateSource)
	r.Get("/api/sources/{id}/fields", api.ListSourceFields)
	r.Get("/api/fields/{id}", api.FieldDetail)
	r.Get("/api/sources/stats", api.SourceStats)
	r.Post("/api/sources/{id}/scan", api.ScanField)

	r.Get("/api/policies", api.ListPolicies)
	r.Get("/api/policies/{id}/rules", api.PolicyRules)
	r.Get("/api/policies/{id}/coverage", api.PolicyCoverage)
	r.Post("/api/policies", api.CreatePolicy)
	r.Post("/api/policies/{id}/rules", api.AddRule)
	r.Post("/api/policies/{id}/publish", api.PublishPolicy)
	r.Post("/api/policies/{id}/rollback", api.RollbackPolicy)

	r.Post("/api/masking", api.StartMasking)
	r.Get("/api/masking", api.ListMasking)
	r.Get("/api/masking/stats", api.MaskingStats)
	r.Get("/api/masking/{id}", api.GetMasking)
	r.Post("/api/masking/{id}/preflight", api.MaskingPreflight)
	r.Post("/api/masking/{id}/apply", api.ApplyRow)
	r.Post("/api/masking/{id}/batch", api.BatchRows)
	r.Post("/api/masking/{id}/commit", api.CommitResult)
	r.Post("/api/masking/{id}/finish", api.FinishMasking)

	r.Get("/api/models", api.ListModels)
	r.Post("/api/models", api.UpdateModel)
	r.Post("/api/policies/{id}/validate", api.ValidatePolicy)
	r.Get("/api/quota", api.TenantQuota)
	r.Post("/api/classify", api.ClassifyFieldValue)

	r.Get("/api/audit", api.ListAudit)
	return r
}
