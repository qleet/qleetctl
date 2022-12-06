package install

import "fmt"

const (
	WorkloadControllerManifestPath = "/tmp/qleet-workload-controller.yaml"
)

// WorkloadControllerManifest returns a yaml manifest for the workload controller
// with the namespace included.
func WorkloadControllerManifest() string {
	return fmt.Sprintf(`---
apiVersion: v1
kind: Secret
metadata:
  name: workload-controller-config
  namespace: %[1]s
type: Opaque
stringData:
  API_SERVER: http://threeport-api-server
  MSG_BROKER_HOST: threeport-message-broker
  MSG_BROKER_PORT: "4222"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: threeport-workload-controller
  namespace: %[1]s
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: threeport-workload-controller
  template:
    metadata:
      labels:
        app.kubernetes.io/name: threeport-workload-controller
    spec:
      containers:
      - name: workload-controller
        image: lander2k2/threeport-workload-controller:latest
        imagePullPolicy: IfNotPresent
        envFrom:
          - secretRef:
              name: workload-controller-config
`, ThreeportControlPlaneNs)
}
