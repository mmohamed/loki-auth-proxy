# Auth Proxy for Grafana Loki

This component is based on [Loki Multi Tenant proxy](https://github.com/k8spin/loki-multi-tenant-proxy) project but with same features in addition. 
[Loki](https://github.com/grafana/loki) support multi-tenant by default by using the `X-Scope-OrgID` extracted from evry HTTP Request on his API. Meanwhile, Loki don'st provide an authentication module.

To responde to this case ([@see issue](https://github.com/grafana/loki/issues/701)), we need another compoenent to provide the authentication layer between Loki and any client (Promtail, Logstash, ...).

This component is a simple HTTP proxy (writed in golang) with HTTP Basic authentication based on users data.

## Set up

1- Set `auth_enabled: true`in your configration file.
2- Configure and deploy the proxy in front of your Loki instance, for version 2.7.x, we need to patch `Readers`, `Writes` and the `Loki Gateway` to add the proxy as Sidecar like (For this release we cant add an extra conteainer with Hel Chart values):
```yaml
spec:
  template:
    spec:
      containers:
        - name: reverse-proxy
          image: medinvention/loki-multi-tenant-proxy:1.0.1
          args:
            - "run"
            - "--port=3101"
            - "--loki-server=http://localhost:3100"
            - "--auth-config=/etc/reverse-proxy-conf/auth.yaml"
            - "--org-check"
          ports:
            - name: http
              containerPort: 3101
              protocol: TCP
          resources:
            limits:
              cpu: 250m
              memory: 200Mi
            requests:
              cpu: 50m
              memory: 40Mi
          volumeMounts:
            - name: reverse-proxy-auth-config
              mountPath: /etc/reverse-proxy-conf
      volumes:
        - name: reverse-proxy-auth-config
          secret:
            secretName: reverse-proxy-auth-config
```

Then patch services (and keep same port):
```yaml
spec:
    ports:
      - name: http-metrics
        port: 3100
        protocol: TCP
        targetPort: http
      - name: grpc
        port: 9095
        protocol: TCP
        targetPort: grpc
```

For this example, the authentication configuration file will be provided by the secret `reverse-proxy-auth-config` like :
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: reverse-proxy-auth-config
type: Opaque
StringData:
  auth.yaml: |-
    users:
      - username: user-tenant-1
        password: pass-tenant-1
        orgid: tenant-1
      - username: user-tenant-2
        password: pass-tenant-2
        orgid: tenant-2
      ...
```

*For services we need to use the merge patch command of kubectl*

Availabel options:
- `--port`: Proxy port.
- `--loki-server`: Loki server URL.
- `--auth-config`: Authentication configuration file path.
- `--org-check`: To force `X-Scope-OrgID` checking to match `orgid` of the authentication configuration file (default: false).

**/!\ A tenant can contains multiple users. and a user can be tied to a multiple tenant.**

3- Configure Clients (like Promtail), to set the username and the password :
```yaml
...
client:
  url: http://loki:3501/loki/api/v1/push
  basic_auth:
    username: user-tenant-2
    password: pass-tenant-2
...
```

**/!\ If orgCheck is activated and an account with "*" configured as OrgID, the request orgID will be sended to Loki**

---

## Build yours

If you want to build it from this repository, follow the instructions bellow:

```bash
docker build --tag loki-proxy:local . -f build/Dockerfile
# Formulti plateform 
# docker buildx build --push --platform linux/arm64,linux/amd64 --tag loki-proxy:local . -f build/Dockerfile
```
