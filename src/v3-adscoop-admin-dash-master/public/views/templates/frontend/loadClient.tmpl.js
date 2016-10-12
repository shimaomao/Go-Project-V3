{{ .Callback }}({
	Host: "http://{{ .Host }}",
	ET: {{ .ET }},
	MinTimeout: {{ .MinTimeout }},
	MaxTimeout: {{ .MaxTimeout }},
	Redir: "{{ .Redir }}",
	EnableUnloadTracking: {{ .EnableUnloadTracking }}
})