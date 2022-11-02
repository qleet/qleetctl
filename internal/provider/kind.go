package provider

import "fmt"

const (
	QleetKindClusterName = "qleet-os"
	QleetKindConfigPath  = "/tmp/qleet-kind-config.yaml"
)

func KindConfig() string {
	return fmt.Sprintf(`kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: %s
nodes:
- role: control-plane
- role: worker
- role: worker
`, QleetKindClusterName)
}
