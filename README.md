# Starting dev environment
1. Build Docker image
```bash
docker build -t zbsss/snippetbox:0.1.4  -f deploy/docker/snippetbox/Dockerfile .
```

2. Publish Docker image
```bash
docker push zbsss/snippetbox:0.1.4
```

3. Update `appVersion` in `deploy/helm/snippetbox/Chart.yaml` to match the image tag

4. Start KinD cluster
```bash
cd infra/terraform/dev
terraform apply --auto-approve
cd -
```

5. Generate TLS certificates
```bash
mkdir deploy/helm/snippetbox/tls
cd deploy/helm/snippetbox/tls
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
cd ..
```

6. Update chart dependencies
```bash
helm dependency update
helm dependency build
```

7. Install helm chart
```bash
helm upgrade --install  snippetbox-chart . --values values.yaml --values environments/dev.values.yaml
```

8. (Optional) Test Helm release
```bash
helm test snippetbox-chart
```

9. Open browser at https://localhost/user/login

# TODOs
- [ ] Create scripts for database migrations
- [ ] Clean up Github Actions
  - [ ] remove sleep, instead wait for pods to be running
  - [ ] add integration tests that make requests to localhost
    - [ ] signup
    - [ ] login
    - [ ] create snippet
    - [ ] view snippet
