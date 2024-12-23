{{/*
Expand the name of the chart.
*/}}
{{- define "centaurus.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "centaurus.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "centaurus.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "centaurus.labels" -}}
helm.sh/chart: {{ include "centaurus.chart" . }}
{{ include "centaurus.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "centaurus.selectorLabels" -}}
app.kubernetes.io/name: {{ include "centaurus.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Generate certificates
*/}}
{{- define "helm-chart.gen-certs" -}}
{{- $ca := genCA "centaurus-ca" 1825 -}}
{{- $host := default "centaurus.local" .Values.tls.host }}
{{- $cert := genSignedCert $host nil (list $host) 1825 $ca -}}
tls.crt: {{ $cert.Cert | b64enc }}
tls.key: {{ $cert.Key | b64enc }}
{{- end -}}


{{/*
Default TLS secret name
*/}}
{{- define "helm-chart.default-tls-secret-name" -}}
centaurus-tls-secret
{{- end }}

{{- define "helm-chart.centaurus-tls-secret-name" -}}
{{ .Values.tls.secretName | default (include "helm-chart.default-tls-secret-name" .) }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "centaurus.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "centaurus.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}