{{- define "creality2mqtt.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "creality2mqtt.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s" .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "creality2mqtt.labels" -}}
app.kubernetes.io/name: {{ include "creality2mqtt.name" . }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "creality2mqtt.selectorLabels" -}}
app.kubernetes.io/name: {{ include "creality2mqtt.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- /* Resolve image reference with sensible defaults.
	If .Values.image.tag is empty, fall back to .Chart.AppVersion.
	Returns full image string: <repo>:<tag> */ -}}
{{- define "creality2mqtt.image" -}}
{{- $repo := .Values.image.repository -}}
{{- $tag := default .Chart.AppVersion .Values.image.tag -}}
{{- printf "%s:%s" $repo $tag -}}
{{- end -}}
