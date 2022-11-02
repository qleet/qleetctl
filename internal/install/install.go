package install

const (
	APIDepsManifestPath   = "/tmp/qleet-api-deps-manifest.yaml"
	APIServerManifestPath = "/tmp/qleet-api-server-manifest.yaml"
	APIDepsManifest       = `---
apiVersion: v1
kind: Namespace
metadata:
  name: threeport-api
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: postgres-config
  namespace: threeport-api
  labels:
    app: postgres
data:
  POSTGRES_DB: threeport_api
  POSTGRES_USER: tp_rest_api
  POSTGRES_PASSWORD: tp-rest-api-pwd
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: threeport-api
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:14.3-alpine
          imagePullPolicy: "IfNotPresent"
          ports:
            - containerPort: 5432
          envFrom:
            - configMapRef:
                name: postgres-config
          volumeMounts:
            - mountPath: /var/lib/postgresql/data
              name: postgredb
            - mountPath: /docker-entrypoint-initdb.d
              name: postgresql-initdb
      volumes:
        - name: postgredb
          emptyDir: {}
        - name: postgresql-initdb
          configMap:
            name: postgres-config-data
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: threeport-api
  labels:
    app: postgres
spec:
  ports:
   - port: 5432
  selector:
   app: postgres
---
# Source: nats/templates/pdb.yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: nats-js
  namespace: threeport-api
  labels:
    helm.sh/chart: nats-0.18.2
    app.kubernetes.io/name: nats
    app.kubernetes.io/instance: nats-js
    app.kubernetes.io/version: "2.9.3"
    app.kubernetes.io/managed-by: Helm
spec:
  maxUnavailable: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: nats
      app.kubernetes.io/instance: nats-js
---
# Source: nats/templates/rbac.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nats-js
  namespace: threeport-api
  labels:
    helm.sh/chart: nats-0.18.2
    app.kubernetes.io/name: nats
    app.kubernetes.io/instance: nats-js
    app.kubernetes.io/version: "2.9.3"
    app.kubernetes.io/managed-by: Helm
---
# Source: nats/templates/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nats-js-config
  namespace: threeport-api
  labels:
    helm.sh/chart: nats-0.18.2
    app.kubernetes.io/name: nats
    app.kubernetes.io/instance: nats-js
    app.kubernetes.io/version: "2.9.3"
    app.kubernetes.io/managed-by: Helm
data:
  nats.conf: |
    # NATS Clients Port
    port: 4222

    # PID file shared with configuration reloader.
    pid_file: "/var/run/nats/nats.pid"

    ###############
    #             #
    # Monitoring  #
    #             #
    ###############
    http: 8222
    server_name:$POD_NAME
    ###################################
    #                                 #
    # NATS JetStream                  #
    #                                 #
    ###################################
    jetstream {
      max_mem: 30Mi
    }
    lame_duck_grace_period: 10s
    lame_duck_duration: 30s
---
# Source: nats/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: nats-js
  namespace: threeport-api
  labels:
    helm.sh/chart: nats-0.18.2
    app.kubernetes.io/name: nats
    app.kubernetes.io/instance: nats-js
    app.kubernetes.io/version: "2.9.3"
    app.kubernetes.io/managed-by: Helm
spec:
  selector:
    app.kubernetes.io/name: nats
    app.kubernetes.io/instance: nats-js
  clusterIP: None
  publishNotReadyAddresses: true
  ports:
  - name: client
    port: 4222
    appProtocol: tcp
  - name: cluster
    port: 6222
    appProtocol: tcp
  - name: monitor
    port: 8222
    appProtocol: http
  - name: metrics
    port: 7777
    appProtocol: http
  - name: leafnodes
    port: 7422
    appProtocol: tcp
  - name: gateways
    port: 7522
    appProtocol: tcp
---
# Source: nats/templates/nats-box.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nats-js-box
  namespace: threeport-api
  labels:
    app: nats-js-box
    chart: nats-0.18.2
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nats-js-box
  template:
    metadata:
      labels:
        app: nats-js-box
    spec:
      volumes:
      containers:
      - name: nats-box
        image: natsio/nats-box:0.13.2
        imagePullPolicy: IfNotPresent
        resources:
          null
        env:
        - name: NATS_URL
          value: nats-js
        command:
        - "tail"
        - "-f"
        - "/dev/null"
        volumeMounts:
---
# Source: nats/templates/statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: nats-js
  namespace: threeport-api
  labels:
    helm.sh/chart: nats-0.18.2
    app.kubernetes.io/name: nats
    app.kubernetes.io/instance: nats-js
    app.kubernetes.io/version: "2.9.3"
    app.kubernetes.io/managed-by: Helm
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: nats
      app.kubernetes.io/instance: nats-js
  replicas: 1
  serviceName: nats-js

  podManagementPolicy: Parallel

  template:
    metadata:
      annotations:
        prometheus.io/path: /metrics
        prometheus.io/port: "7777"
        prometheus.io/scrape: "true"
        checksum/config: 3b398e973c292bf8c2eb90d62acb846274c0489643aad560d8c4aed123f20ce7
      labels:
        app.kubernetes.io/name: nats
        app.kubernetes.io/instance: nats-js
    spec:
      # Common volumes for the containers.
      volumes:
      - name: config-volume
        configMap:
          name: nats-js-config

      # Local volume shared with the reloader.
      - name: pid
        emptyDir: {}

      #################
      #               #
      #  TLS Volumes  #
      #               #
      #################

      serviceAccountName: nats-js

      # Required to be able to HUP signal and apply config
      # reload to the server without restarting the pod.
      shareProcessNamespace: true

      #################
      #               #
      #  NATS Server  #
      #               #
      #################
      terminationGracePeriodSeconds: 60
      containers:
      - name: nats
        image: nats:2.9.3-alpine
        imagePullPolicy: IfNotPresent
        resources:
          {}
        ports:
        - containerPort: 4222
          name: client
        - containerPort: 6222
          name: cluster
        - containerPort: 8222
          name: monitor

        command:
        - "nats-server"
        - "--config"
        - "/etc/nats-config/nats.conf"

        # Required to be able to define an environment variable
        # that refers to other environment variables.  This env var
        # is later used as part of the configuration file.
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: SERVER_NAME
          value: $(POD_NAME)
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: CLUSTER_ADVERTISE
          value: $(POD_NAME).nats-js.$(POD_NAMESPACE).svc.cluster.local
        volumeMounts:
        - name: config-volume
          mountPath: /etc/nats-config
        - name: pid
          mountPath: /var/run/nats
        

        #######################
        #                     #
        # Healthcheck Probes  #
        #                     #
        #######################
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: 8222
          initialDelaySeconds: 10
          periodSeconds: 30
          successThreshold: 1
          timeoutSeconds: 5
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: 8222
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 5
        startupProbe:
          # for NATS server versions >=2.7.1, /healthz will be enabled
          # startup probe checks that the JS server is enabled, is current with the meta leader,
          # and that all streams and consumers assigned to this JS server are current
          failureThreshold: 30
          httpGet:
            path: /healthz
            port: 8222
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 5

        # Gracefully stop NATS Server on pod deletion or image upgrade.
        #
        lifecycle:
          preStop:
            exec:
              # send the lame duck shutdown signal to trigger a graceful shutdown
              # nats-server will ignore the TERM signal it receives after this
              #
              command:
              - "nats-server"
              - "-sl=ldm=/var/run/nats/nats.pid"

      #################################
      #                               #
      #  NATS Configuration Reloader  #
      #                               #
      #################################
      - name: reloader
        image: natsio/nats-server-config-reloader:0.7.4
        imagePullPolicy: IfNotPresent
        resources:
          null
        command:
        - "nats-server-config-reloader"
        - "-pid"
        - "/var/run/nats/nats.pid"
        - "-config"
        - "/etc/nats-config/nats.conf"
        volumeMounts:
        - name: config-volume
          mountPath: /etc/nats-config
        - name: pid
          mountPath: /var/run/nats
        

      ##############################
      #                            #
      #  NATS Prometheus Exporter  #
      #                            #
      ##############################
      - name: metrics
        image: natsio/prometheus-nats-exporter:0.10.0
        imagePullPolicy: IfNotPresent
        resources:
          {}
        args:
        - -connz
        - -routez
        - -subz
        - -varz
        - -prefix=nats
        - -use_internal_server_id
        - -jsz=all
        - http://localhost:8222/
        ports:
        - containerPort: 7777
          name: metrics

  volumeClaimTemplates:
---
# Source: nats/templates/tests/test-request-reply.yaml
apiVersion: v1
kind: Pod
metadata:
  name: "nats-js-test-request-reply"
  labels:
    chart: nats-0.18.2
    app: nats-js-test-request-reply
  annotations:
    "helm.sh/hook": test
spec:
  containers:
  - name: nats-box
    image: synadia/nats-box
    env:
    - name: NATS_HOST
      value: nats-js
    command:
    - /bin/sh
    - -ec
    - |
      nats reply -s nats://$NATS_HOST:4222 'name.>' --command "echo 1" &
    - |
      "&&"
    - |
      name=$(nats request -s nats://$NATS_HOST:4222 name.test '' 2>/dev/null)
    - |
      "&&"
    - |
      [ $name = test ]

  restartPolicy: Never
`
	APIServerManifest = `---
apiVersion: v1
kind: Secret
metadata:
  name: db-config
  namespace: threeport-api
stringData:
  env: |
    DB_HOST=postgres
    DB_USER=tp_rest_api
    DB_PASSWORD=tp-rest-api-pwd
    DB_NAME=threeport_api
    DB_PORT=5432
    DB_SSL_MODE=disable
    NATS_HOST=nats-js
    NATS_PORT=4222
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: threeport-api-server
  namespace: threeport-api
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: threeport-api-server
  template:
    metadata:
      labels:
        app.kubernetes.io/name: threeport-api-server
    spec:
      containers:
      - name: api-server
        image: lander2k2/threeport-rest-api:latest
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 1323
          name: http
          protocol: TCP
        volumeMounts:
        - name: db-config
          mountPath: "/etc/threeport/"
      volumes:
      - name: db-config
        secret:
          secretName: db-config
---
apiVersion: v1
kind: Service
metadata:
  name: threeport-api-server
  namespace: threeport-api
spec:
  selector:
    app.kubernetes.io/name: threeport-api-server
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 1323
`
)
