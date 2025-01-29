{{- range .Values.maps -}}
{{- maplist . $.Values.mapsRaw }}
{{- end -}}
