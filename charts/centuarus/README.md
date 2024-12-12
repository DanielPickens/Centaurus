# centaurus helm chart

centaurus allows you to manage Kubernetes clusters. This Helm chart simplifies the installation and configuration of centaurus.

## Installation

To install the centaurus chart using Helm, run the following command:

```bash
helm install centaurus oci://ghcr.io/centaurus/charts/centaurus -n centaurus-system --create-namespace
```

### Notes:

- **Default Setup**: By default, centaurus runs on port `8443` with self-signed certificates.
- **Namespace**: A new namespace `centaurus-system` will be created automatically if it doesn't exist.

### Using Custom TLS Certificates

To use your own TLS certificates instead of the default self-signed ones:

1. **Create a Kubernetes Secret**: Store your TLS certificate and key in a secret.

   ```bash
   kubectl create namespace centaurus-system
   kubectl -n centaurus-system create secret tls centaurus-tls-secret --cert=tls.crt --key=tls.key
   ```

2. **Install centaurus with your certificates**:

   ```bash
   helm install centaurus oci://ghcr.io/centaurus/centaurus \
     -n centaurus-system --version v0.0.4 --create-namespace \
     --set tls.secretName=centaurus-tls-secret
   ```

### Using a Custom Service Account

By default, the chart creates a service account with `admin` RBAC permissions in the release namespace. If you'd like centaurus to use an existing service account, you can disable the creation of a new one.

1. **Install centaurus with an existing service account**:

   ```bash
   helm install centaurus oci://ghcr.io/centaurus/centaurus \
     -n centaurus-system --version v0.0.4 --create-namespace \
     --set serviceAccount.create=false \
     --set serviceAccount.name=<yourServiceAccountName>
   ```

## Upgrading the Chart

To upgrade to a newer version of the chart, run the following command:

```bash
helm upgrade centaurus oci://ghcr.io/centaurus/centaurus \
  -n centaurus-system --version v0.0.4
```

## Configuration Parameters

The following are some key configuration parameters you can customize when installing the chart:

| Parameter               | Description                                                                                       | Default  |
|-------------------------|---------------------------------------------------------------------------------------------------|----------|
| `tls.secretName`         | Kubernetes secret name containing your TLS certificate and key. Must be in the `centaurus-system` namespace. | `""`     |
| `service.port`           | The HTTPS port number centaurus listens on.                                                        | `8443`   |
| `serviceAccount.create`  | Set to `false` if you want to use an existing service account.                                     | `true`   |
| `serviceAccount.name`    | Name of the service account to use (if `serviceAccount.create=false`).                            | `""`     |

For a complete list of configurable parameters, refer to the values file or documentation.
