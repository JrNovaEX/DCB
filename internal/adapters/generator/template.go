package generator

// composeTemplate is the master docker-compose.yml template.
// It is compiled once at startup via template.Must().
//
// Template data type: *domain.ProjectConfig (enriched with ResolvedHealthcheck).
const composeTemplate = `version: "{{.Version}}"

services:
{{- range $name, $svc := .Services}}
  {{$name}}:
{{- if $svc.Image}}
    image: {{$svc.Image}}
{{- end}}
{{- if $svc.Build}}
    build:
{{- if $svc.Build.Dockerfile}}
      context: {{$svc.Build.Context}}
      dockerfile: {{$svc.Build.Dockerfile}}
{{- if $svc.Build.Args}}
      args:
{{- range $k, $v := $svc.Build.Args}}
        {{$k}}: {{$v}}
{{- end}}
{{- end}}
{{- else}}
      context: {{$svc.Build.Context}}
{{- end}}
{{- end}}
{{- if $svc.Port}}
    ports:
      - "{{$svc.Port.Int}}:{{$svc.Port.Int}}"
{{- end}}
{{- if $svc.EnvFile}}
    env_file:
      - {{$svc.EnvFile}}
{{- end}}
{{- if $svc.Env}}
    environment:
{{- range $k, $v := $svc.Env}}
      {{$k}}: {{$v}}
{{- end}}
{{- end}}
{{- if $svc.DependsOn}}
    depends_on:
{{- range $svc.DependsOn}}
      {{.}}:
        condition: service_healthy
{{- end}}
{{- end}}
{{- if $svc.ResolvedHealthcheck}}
    healthcheck:
      test: {{formatTest $svc.ResolvedHealthcheck.Test}}
      interval: {{$svc.ResolvedHealthcheck.Interval}}
      timeout: {{$svc.ResolvedHealthcheck.Timeout}}
      retries: {{$svc.ResolvedHealthcheck.Retries}}
      start_period: {{$svc.ResolvedHealthcheck.StartPeriod}}
{{- end}}
{{- if $svc.Volume}}
    volumes:
      - {{$svc.Volume}}
{{- end}}
    restart: {{restartPolicy $svc.Restart}}
    networks:
      - {{$.Project}}_default
{{- end}}

{{- if .Volumes}}

volumes:
{{- range $name, $_ := .Volumes}}
  {{$name}}:
{{- end}}
{{- end}}

networks:
  {{.Project}}_default:
    driver: bridge
`
