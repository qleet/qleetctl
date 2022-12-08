package provider

import (
	"fmt"

	"github.com/qleet/qleetctl/internal/install"
)

const (
	QleetKindConfigPath = "/tmp/qleet-kind-config.yaml"
)

func GetQleetKindClusterName(qleetOSInstanceName string) string {
	return fmt.Sprintf("qleet-os-%s", qleetOSInstanceName)
}

func KindConfig(qleetOSInstanceName string) string {
	return fmt.Sprintf(`kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: %[1]s
nodes:
- role: control-plane
- role: worker
  extraPortMappings:
    - containerPort: %[2]s
      hostPort: %[2]s
      protocol: TCP
`, GetQleetKindClusterName(qleetOSInstanceName), install.QleetOSAPIPort)
}
