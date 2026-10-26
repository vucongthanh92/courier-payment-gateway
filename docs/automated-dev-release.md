# Automated development release

Merges to `main` run tests, publish an immutable `linux/arm64` image to the
private `courier/payment-gateway` Amazon ECR repository, and use the Courier
Dev Release Bot GitHub App to update the image tag in `courier-common`.

The bot creates and automatically merges a development release PR. Argo CD
then deploys the exact SHA-tagged image from central GitOps; it never deploys
`latest`.

Required repository secrets:

- `COURIER_RELEASE_APP_ID`
- `COURIER_RELEASE_APP_PRIVATE_KEY`
