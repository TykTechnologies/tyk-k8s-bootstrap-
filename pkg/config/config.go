package config

import (
	"tyk/tyk/bootstrap/tyk/api"

	"github.com/kelseyhightower/envconfig"
)

const prefix = "TYK_K8SBOOTSTRAP"

type Config struct {
	// Log sets the level of the logrus logger.
	// The default is `info`
	Log string `default:"info"`

	// InsecureSkipVerify enables InsecureSkipVerify options in http request sent to Tyk - might be useful
	// for Tyk Dashboard with self-signed certs.
	InsecureSkipVerify bool
	// BootstrapDashboard controls bootstrapping Tyk Dashboard or not.
	BootstrapDashboard bool
	// BootstrapPortal controls bootstrapping Tyk Classic Portal or not.
	BootstrapPortal bool

	// OperatorKubernetesSecretName corresponds to the Kubernetes secret name that will be created for Tyk Operator.
	// Set it to an empty string to disable bootstrapping Kubernetes secret for Tyk Operator.
	OperatorKubernetesSecretName string
	// OperatorLicenseKey corresponds to Operator license key
	OperatorLicenseKey string
	// DevPortalKubernetesSecretName corresponds to the Kubernetes secret name that will be created for
	// Tyk Developer Portal. Set it to an empty to string to disable bootstrapping Kubernetes
	// secret for Tyk Developer Portal.
	DevPortalKubernetesSecretName string
	// K8s consists of configurations for Kubernetes services of Tyk.
	K8s K8sConf
	// Tyk consists of configurations for Tyk components such as Tyk Dashboard Admin information
	// or Tyk Portal configurations.
	Tyk TykConf
}

type K8sConf struct {
	// DashboardSvcUrl corresponds to the URL of Tyk Dashboard.
	DashboardSvcUrl string
	// DashboardSvcProto corresponds to Tyk Dashboard Service Protocol (either http or https).
	DashboardSvcProto string
	// ReleaseNamespace corresponds to the namespace where Tyk is deployed via Helm Chart.
	ReleaseNamespace string
	// DashboardDeploymentName corresponds to the name of the Tyk Dashboard Deployment, which is being used
	// to restart Dashboard pod after bootstrapping. By default, it discovers Dashboard Deployment name.
	// If the environment variable is populated, the discovery will not be triggered.
	DashboardDeploymentName string
}

type TykAdmin struct {
	// Secret corresponds to the secret that will be used in Admin APIs.
	Secret string
	// FirstName corresponds to the first name of the admin being created.
	FirstName string
	// LastName corresponds to the last name of the admin being created.
	LastName string
	// EmailAddress corresponds to the email address of the admin being created.
	EmailAddress string
	// Password corresponds to the password of the admin being created.
	Password string

	// Auth corresponds to Tyk Dashboard API Access Credentials of the admin user, and it will be used
	// in Authorization header of the HTTP requests that will be sent to Tyk for bootstrapping.
	// Also, if bootstrapping Tyk Operator Secret is enabled Auth corresponds to TykAuth field in the
	// Kubernetes secret of Tyk Operator.
	// Set it if Tyk Dashboard is already bootstrapped.
	Auth string
}

type TykOrg struct {
	// Name corresponds to the name for your organization that is going to be bootstrapped in Tyk
	Name string
	// Cname corresponds to the Organisation CNAME which is going to bind the Portal to.
	Cname string

	// ID corresponds to the organisation ID that is being created.
	ID string

	// Hybrid includes details of hybrid organisation while using MDCB Control Plane
	Hybrid *HybridConf
}

type TykConf struct {
	// Admin consists of configurations for Tyk Dashboard Admin.
	Admin TykAdmin
	// Org consists of configurations for the organisation that is going to be created in Tyk Dashboard.
	Org TykOrg

	// DashboardLicense corresponds to the license key of Tyk Dashboard.
	DashboardLicense string
}

type HybridConf struct {
	// Enabled specified if the Hybrid organisation is enabled or not
	Enabled bool
	// KeyEvent corresponds to `key_event` of the event options which enables key events such as updates and deletes,
	// to be propagated to the various instance zones.
	KeyEvent *api.EventConfig
	// HashedKeyEvent corresponds to `hashed_key_event` of the event options which enables key events such as updates
	// and deletes, to be propagated to the various instance zones.
	HashedKeyEvent *api.EventConfig `json:",omitempty"`
}

func NewConfig() (*Config, error) {
	conf := &Config{}
	if err := envconfig.Process(prefix, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
