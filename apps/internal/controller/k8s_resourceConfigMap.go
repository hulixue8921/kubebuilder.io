package controller

import (
	"bytes"
	"text/template"

	coreV1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	appv1 "kubebuilder.io/apps/api/v1"
)

const (
	Eshost = "192.168.1.1"
)

type EsInfo struct {
	Log       string
	Type      string
	EsHost    string
	LogFormat string
}

func CreateConfigMap(object *appv1.DeployObject) *coreV1.ConfigMap {
	req := object.DeepCopy()
	es := EsInfo{
		Log:       req.Spec.AppLogDir + "/*.log",
		Type:      req.Name,
		EsHost:    Eshost,
		LogFormat: req.Spec.LogFormat,
	}

	logFile := `
	input {
		file {
			path => ["{{.Log}}" ]
			type => "{{.Type}}"
			start_position => "beginning"
			codec => multiline {
				pattern => "{{.LogFormat}}"
				negate => true
				what => "previous"
			}
		}
	}
	filter {
	}
	output {
		elasticsearch {
			 hosts => [ "{{.EsHost}}" ]
			 index => "{{.Type}}_%{+YYYY.MM.dd}"
		} 
	}`

	tmpl, _ := template.New("").Parse(logFile)
	var buf bytes.Buffer
	tmpl.Execute(&buf, es)

	return &coreV1.ConfigMap{
		TypeMeta: meta.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		Data: map[string]string{
			"logstash.yml": "path.config: /usr/share/logstash/conf.d/*.conf",
			"log.conf":     buf.String(),
		},
	}

}
