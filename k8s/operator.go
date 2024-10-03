package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tykModePro            = "pro"
	tykModeKey            = "TYK_MODE"
	tykAuthKey            = "TYK_AUTH"
	tykOrgKey             = "TYK_ORG"
	tykURLKey             = "TYK_URL"
	tykOperatorLicenseKey = "TYK_OPERATOR_LICENSEKEY"
)

// BootstrapTykOperatorSecret bootstrap a Kubernetes Secret utilized by Tyk Operator.
// If the system has the secret created already, it deletes the existing one and recreates
// a secret for Tyk Operator.
func (c *Client) BootstrapTykOperatorSecret() error {
	if err := c.deleteSecret(c.appArgs.OperatorKubernetesSecretName, true); err != nil {
		return err
	}

	secretData := map[string][]byte{
		tykAuthKey:            []byte(c.appArgs.Tyk.Admin.Auth),
		tykOrgKey:             []byte(c.appArgs.Tyk.Org.ID),
		tykModeKey:            []byte(tykModePro),
		tykURLKey:             []byte(c.appArgs.K8s.DashboardSvcUrl),
		tykOperatorLicenseKey: []byte(c.appArgs.OperatorLicenseKey),
	}

	objectMeta := metav1.ObjectMeta{Name: c.appArgs.OperatorKubernetesSecretName}

	secret := v1.Secret{
		ObjectMeta: objectMeta,
		Data:       secretData,
	}

	_, err := c.clientSet.
		CoreV1().
		Secrets(c.appArgs.K8s.ReleaseNamespace).
		Create(context.TODO(), &secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
