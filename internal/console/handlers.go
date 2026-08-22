package console

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"maskhub/internal/audit"
	"maskhub/internal/classify"
	"maskhub/internal/detect"
	"maskhub/internal/mask"
	"maskhub/internal/policy"
	"maskhub/internal/quota"
	"maskhub/internal/source"
	"maskhub/internal/store"
)

var (
	errNotFound = &httpError{status: http.StatusNotFound, message: "not found"}
	errBadBody  = &httpError{status: http.StatusBadRequest, message: "bad request body"}
)

type httpError struct {
	status  int
	message string
}

func (e *httpError) Error() string {
	return e.message
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, err error) {
	if he, ok := err.(*httpError); ok {
		status = he.status
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func chiURL(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func decodeBody(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func itoa(v int) string {
	return strconv.Itoa(v)
}

func itoa64(v int64) string {
	return strconv.FormatInt(v, 10)
}

func (a *API) Stats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.state.Summarize())
}

func (a *API) ListSources(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, source.List(a.state))
}

func (a *API) CreateSource(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantID string `json:"tenant_id"`
		Name     string `json:"name"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, errBadBody)
		return
	}
	s, err := source.Create(a.state, body.TenantID, body.Name)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	_ = audit.Record(a.state, "source_create", s.ID, body.Name)
	writeJSON(w, http.StatusCreated, s)
}

func (a *API) ListSourceFields(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	writeJSON(w, http.StatusOK, detect.Fields(a.state, id))
}

func (a *API) SourceStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, source.ListStats(a.state))
}

func (a *API) ScanField(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	var body struct {
		Name string `json:"name"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, errBadBody)
		return
	}
	f, err := detect.Scan(a.state, id, body.Name)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	_ = audit.Record(a.state, "field_scan", f.ID, body.Name)
	writeJSON(w, http.StatusCreated, f)
}

func (a *API) ListPolicies(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, policy.List(a.state))
}

func (a *API) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantID string `json:"tenant_id"`
		Name     string `json:"name"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, errBadBody)
		return
	}
	p, err := policy.Create(a.state, body.TenantID, body.Name)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	_ = audit.Record(a.state, "policy_create", p.ID, body.Name)
	writeJSON(w, http.StatusCreated, p)
}

func (a *API) AddRule(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	var body struct {
		Category string `json:"category"`
		Pattern  string `json:"pattern"`
		Replacer string `json:"replacer"`
		Order    int    `json:"order"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, errBadBody)
		return
	}
	rl := &store.Rule{
		ID:       newID(),
		PolicyID: id,
		Category: body.Category,
		Pattern:  body.Pattern,
		Replacer: body.Replacer,
		Order:    body.Order,
	}
	if err := a.state.PutRule(rl); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusCreated, rl)
}

func (a *API) PolicyRules(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	writeJSON(w, http.StatusOK, policy.Rules(a.state, id))
}

func (a *API) PolicyCoverage(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"policy_id": id,
		"rules":     policy.RuleCount(a.state, id),
		"categories": policy.CategoryCoverage(a.state, id),
	})
}

func (a *API) PublishPolicy(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	p, err := policy.Publish(a.state, id)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	_ = audit.Record(a.state, "policy_publish", id, "v"+itoa(p.Version))
	writeJSON(w, http.StatusOK, p)
}

func (a *API) RollbackPolicy(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	p, err := policy.Rollback(a.state, id)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	_ = audit.Record(a.state, "policy_rollback", id, "v"+itoa(p.Version))
	writeJSON(w, http.StatusOK, p)
}

func (a *API) StartMasking(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantID string `json:"tenant_id"`
		SourceID string `json:"source_id"`
		PolicyID string `json:"policy_id"`
		Rows     int    `json:"rows"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, errBadBody)
		return
	}
	j, err := mask.Execute(a.state, body.TenantID, body.SourceID, body.PolicyID, body.Rows)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusCreated, j)
}

func (a *API) MaskingPreflight(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	if err := mask.Preflight(a.state, id); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (a *API) ListMasking(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.state.Jobs())
}

func (a *API) GetMasking(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	j, ok := a.state.Job(id)
	if !ok {
		writeErr(w, http.StatusNotFound, mask.ErrNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"job":     j,
		"results": a.state.ResultsForJob(id),
	})
}

func (a *API) MaskingStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, mask.ListStats(a.state))
}

func (a *API) FieldDetail(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	f, ok := a.state.Field(id)
	if !ok {
		writeErr(w, http.StatusNotFound, detect.ErrNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"field":          f,
		"classification": f.Category,
		"status":         f.Status,
	})
}

func (a *API) ApplyRow(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	var body struct {
		RowID    string `json:"row_id"`
		Category string `json:"category"`
		Value    string `json:"value"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, errBadBody)
		return
	}
	res, err := mask.Apply(a.state, id, body.RowID, body.Category, body.Value)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (a *API) BatchRows(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	var body struct {
		Rows []mask.RowInput `json:"rows"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, errBadBody)
		return
	}
	results, err := mask.Batch(a.state, id, body.Rows)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusCreated, results)
}

func (a *API) CommitResult(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ResultID string `json:"result_id"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, errBadBody)
		return
	}
	res, ok := a.state.Result(body.ResultID)
	if !ok {
		writeErr(w, http.StatusNotFound, mask.ErrNotFound)
		return
	}
	if err := mask.Commit(a.state, res); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (a *API) ValidatePolicy(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	if err := policy.ValidateDraft(a.state, id); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "valid"})
}

func (a *API) TenantQuota(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant")
	writeJSON(w, http.StatusOK, map[string]any{
		"tenant":    tenantID,
		"used":      quota.Used(a.state, tenantID),
		"remaining": quota.Remaining(a.state, tenantID),
		"full":      quota.Full(a.state, tenantID),
	})
}

func (a *API) ClassifyFieldValue(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, errBadBody)
		return
	}
	category, err := classify.ApplyModel(a.state, body.Name)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	version, _ := classify.SnapshotVersion(a.state)
	writeJSON(w, http.StatusOK, map[string]any{
		"name":         body.Name,
		"category":     category,
		"model_version": version,
	})
}

func (a *API) FinishMasking(w http.ResponseWriter, r *http.Request) {
	id := chiURL(r, "id")
	j, err := mask.Finish(a.state, id)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, j)
}

func (a *API) ListModels(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.state.Models())
}

func (a *API) UpdateModel(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Snapshot string `json:"snapshot"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, errBadBody)
		return
	}
	m, err := classify.Update(a.state, body.Snapshot)
	if err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

func (a *API) ListAudit(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.state.RecentAudit(100))
}
