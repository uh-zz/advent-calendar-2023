//go:generate tinygo build -o plugin.wasm -scheduler=none -target=wasi --no-debug plugin.go
//go:build tinygo.wasm

package main

import (
	"fmt"

	"github.com/aquasecurity/trivy/pkg/module/api"
	"github.com/aquasecurity/trivy/pkg/module/serialize"
	"github.com/aquasecurity/trivy/pkg/module/wasm"
)

func main() {
	wasm.RegisterModule(Plugin{})
}

type Plugin struct{}

func (p Plugin) Version() int {
	return 1
}

func (p Plugin) Name() string {
	return "plugin"
}

func (Plugin) PostScanSpec() serialize.PostScanSpec {
	return serialize.PostScanSpec{
		Action: api.ActionUpdate, // Update severity
		IDs:    []string{"CVE-2023-25652"},
	}
}

func (Plugin) PostScan(results serialize.Results) (serialize.Results, error) {
	for i, result := range results {
		for j, vuln := range result.Vulnerabilities {
			if vuln.VulnerabilityID != "CVE-2023-25652" {
				continue
			}

			results[i].Vulnerabilities[j].Severity = "CRITICAL"
			results[i].Vulnerabilities[j].Title = fmt.Sprintf("[ALERT] please alignment to security-team. detail: %s", results[i].Vulnerabilities[j].Title)
		}
	}

	return results, nil
}
