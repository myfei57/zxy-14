package console

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.New("pages").Parse(`
{{define "layout"}}<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>{{.Title}}</title>
<style>
body{font-family:sans-serif;margin:24px;color:#222}
nav a{margin-right:12px}
table{border-collapse:collapse;margin-top:12px}
td,th{border:1px solid #ccc;padding:6px 10px;font-size:13px}
</style></head><body>
<nav><a href="/">总览</a><a href="/console/sources">数据源</a><a href="/console/policies">策略</a><a href="/console/masking">脱敏</a><a href="/console/audit">审计</a></nav>
<h2>{{.Title}}</h2>{{.Body}}</body></html>{{end}}
{{define "index"}}{{template "layout" .}}{{end}}
{{define "sources"}}{{template "layout" .}}{{end}}
{{define "policies"}}{{template "layout" .}}{{end}}
{{define "masking"}}{{template "layout" .}}{{end}}
{{define "audit"}}{{template "layout" .}}{{end}}
`))

type pageData struct {
	Title string
	Body  template.HTML
}

func (a *API) render(w http.ResponseWriter, name, title string, body template.HTML) {
	data := pageData{Title: title, Body: body}
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *API) IndexPage(w http.ResponseWriter, r *http.Request) {
	s := a.state.Summarize()
	body := template.HTML(`<p>数据源 ` + itoa(s.Sources) + ` 个，字段 ` + itoa(s.Fields) + ` 个（已分类 ` + itoa(s.Classified) + `），策略 ` + itoa(s.Policies) + ` 个（生效 ` + itoa(s.ActivePolicies) + `），脱敏任务 ` + itoa(s.Jobs) + ` 个。</p>`)
	a.render(w, "index", "总览", body)
}

func (a *API) SourcesPage(w http.ResponseWriter, r *http.Request) {
	rows := ""
	for _, s := range a.state.Sources() {
		rows += `<tr><td>` + s.Name + `</td><td>` + s.TenantID + `</td><td>` + s.Status + `</td></tr>`
	}
	body := template.HTML(`<table><tr><th>数据源</th><th>租户</th><th>状态</th></tr>` + rows + `</table>`)
	a.render(w, "sources", "数据源", body)
}

func (a *API) PoliciesPage(w http.ResponseWriter, r *http.Request) {
	rows := ""
	for _, p := range a.state.Policies() {
		rows += `<tr><td>` + p.Name + `</td><td>` + p.Status + `</td><td>v` + itoa(p.Version) + `</td></tr>`
	}
	body := template.HTML(`<table><tr><th>策略</th><th>状态</th><th>版本</th></tr>` + rows + `</table>`)
	a.render(w, "policies", "脱敏策略", body)
}

func (a *API) MaskingPage(w http.ResponseWriter, r *http.Request) {
	rows := ""
	for _, j := range a.state.Jobs() {
		rows += `<tr><td>` + j.ID[:8] + `</td><td>` + j.Status + `</td><td>` + itoa64(j.QuotaUsed) + `</td></tr>`
	}
	body := template.HTML(`<table><tr><th>任务</th><th>状态</th><th>配额</th></tr>` + rows + `</table>`)
	a.render(w, "masking", "脱敏执行", body)
}

func (a *API) AuditPage(w http.ResponseWriter, r *http.Request) {
	rows := ""
	for _, e := range a.state.RecentAudit(50) {
		rows += `<tr><td>` + e.Action + `</td><td>` + e.Target + `</td><td>` + e.Detail + `</td><td>` + e.At.Format("2006-01-02 15:04:05") + `</td></tr>`
	}
	body := template.HTML(`<table><tr><th>动作</th><th>对象</th><th>详情</th><th>时间</th></tr>` + rows + `</table>`)
	a.render(w, "audit", "审计日志", body)
}
