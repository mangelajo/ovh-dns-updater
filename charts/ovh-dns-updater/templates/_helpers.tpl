{{/*
Expand the name of the chart.
*/}}
{{- define "ovh-dns-updater.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "ovh-dns-updater.fullname" -}}
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
{{- define "ovh-dns-updater.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "ovh-dns-updater.labels" -}}
helm.sh/chart: {{ include "ovh-dns-updater.chart" . }}
{{ include "ovh-dns-updater.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "ovh-dns-updater.selectorLabels" -}}
app.kubernetes.io/name: {{ include "ovh-dns-updater.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "ovh-dns-updater.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "ovh-dns-updater.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Determine the secret name to use
*/}}
{{- define "ovh-dns-updater.secretName" -}}
{{- if .Values.config.ovhCredentials.create }}
{{- include "ovh-dns-updater.fullname" . }}
{{- else }}
{{- default (include "ovh-dns-updater.fullname" .) .Values.config.ovhCredentials.existingSecret }}
{{- end }}
{{- end }}
