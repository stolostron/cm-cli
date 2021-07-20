{{- define "ocpImage" }}
  {{- $release := splitList ":" .clusterPool.ocpImage }}
  {{- if gt (len $release) 1 }}
    {{- $release = index $release 1 | replace "_" "-" | lower }}
    {{- $release = (print $release "-" .clusterPool.name ) }}
{{- $release }}
  {{- end }}
{{- end }}