package data

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strconv"
	"tyk/tyk/bootstrap/constants"
)

type AppArguments struct {
	DashboardHost                 string
	DashboardPort                 int32
	DashBoardLicense              string
	TykAdminSecret                string
	CurrentOrgName                string
	TykAdminPassword              string
	Cname                         string
	TykAdminFirstName             string
	TykAdminLastName              string
	TykAdminEmailAddress          string
	UserAuth                      string
	OrgId                         string
	CatalogId                     string
	DashboardUrl                  string
	DashboardProto                string
	TykPodNamespace               string
	DashboardSvc                  string
	DashboardInsecureSkipVerify   bool
	IsDashboardEnabled            bool
	OperatorSecretEnabled         bool
	OperatorSecretName            string
	EnterprisePortalSecretEnabled bool
	EnterprisePortalSecretName    string
	BootstrapPortal               bool
	DashboardDeploymentName       string
	ReleaseName                   string
}

var AppConfig = AppArguments{
	DashboardPort:        3000,
	TykAdminSecret:       "12345",
	CurrentOrgName:       "TYKTYK",
	Cname:                "tykCName",
	TykAdminPassword:     "123456",
	TykAdminFirstName:    "firstName",
	TykAdminEmailAddress: "tyk@tyk.io",
	TykAdminLastName:     "lastName",
}

func InitAppDataPreDelete() error {
	AppConfig.OperatorSecretName = os.Getenv(constants.OperatorSecretNameEnvVar)
	AppConfig.EnterprisePortalSecretName = os.Getenv(constants.EnterprisePortalSecretNameEnvVar)
	AppConfig.TykPodNamespace = os.Getenv(constants.TykPodNamespaceEnvVar)
	return nil
}

func InitAppDataPostInstall() error {
	AppConfig.TykAdminFirstName = os.Getenv(constants.TykAdminFirstNameEnvVar)
	AppConfig.TykAdminLastName = os.Getenv(constants.TykAdminLastNameEnvVar)
	AppConfig.TykAdminEmailAddress = os.Getenv(constants.TykAdminEmailEnvVar)
	AppConfig.TykAdminPassword = os.Getenv(constants.TykAdminPasswordEnvVar)
	AppConfig.TykPodNamespace = os.Getenv(constants.TykPodNamespaceEnvVar)
	AppConfig.DashboardProto = os.Getenv(constants.TykDashboardProtoEnvVar)

	AppConfig.DashBoardLicense = os.Getenv(constants.TykDbLicensekeyEnvVar)
	AppConfig.TykAdminSecret = os.Getenv(constants.TykAdminSecretEnvVar)
	AppConfig.CurrentOrgName = os.Getenv(constants.TykOrgNameEnvVar)
	AppConfig.Cname = os.Getenv(constants.TykOrgCnameEnvVar)
	AppConfig.ReleaseName = os.Getenv(constants.ReleaseNameEnvVar)

	var err error

	dashEnabledRaw := os.Getenv(constants.DashboardEnabledEnvVar)
	if dashEnabledRaw != "" {
		AppConfig.IsDashboardEnabled, err = strconv.ParseBool(os.Getenv(constants.DashboardEnabledEnvVar))
		if err != nil {
			return err
		}
	}

	if AppConfig.IsDashboardEnabled {
		if err := discoverDashboardSvc(); err != nil {
			return err
		}
		AppConfig.DashboardUrl = fmt.Sprintf("%s://%s.%s.svc.cluster.local:%d",
			AppConfig.DashboardProto,
			AppConfig.DashboardSvc,
			AppConfig.TykPodNamespace,
			AppConfig.DashboardPort,
		)
	}

	operatorSecretEnabledRaw := os.Getenv(constants.OperatorSecretEnabledEnvVar)
	if operatorSecretEnabledRaw != "" {
		AppConfig.OperatorSecretEnabled, err = strconv.ParseBool(operatorSecretEnabledRaw)
		if err != nil {
			return err
		}
	}
	AppConfig.OperatorSecretName = os.Getenv(constants.OperatorSecretNameEnvVar)

	enterprisePortalSecretEnabledRaw := os.Getenv(constants.EnterprisePortalSecretEnabledEnvVar)
	if enterprisePortalSecretEnabledRaw != "" {
		AppConfig.EnterprisePortalSecretEnabled, err = strconv.ParseBool(enterprisePortalSecretEnabledRaw)
		if err != nil {
			return err
		}
	}
	AppConfig.EnterprisePortalSecretName = os.Getenv(constants.EnterprisePortalSecretNameEnvVar)

	bootstrapPortalBoolRaw := os.Getenv(constants.BootstrapPortalEnvVar)
	if bootstrapPortalBoolRaw != "" {
		AppConfig.BootstrapPortal, err = strconv.ParseBool(bootstrapPortalBoolRaw)
		if err != nil {
			return err
		}
	}
	AppConfig.DashboardDeploymentName = os.Getenv(constants.TykDashboardDeployEnvVar)

	dashboardInsecureSkipVerifyRaw := os.Getenv(constants.TykDashboardInsecureSkipVerify)
	if dashboardInsecureSkipVerifyRaw != "" {
		AppConfig.DashboardInsecureSkipVerify, err = strconv.ParseBool(dashboardInsecureSkipVerifyRaw)
		if err != nil {
			return err
		}
	}

	return nil
}

// discoverDashboardSvc lists Service objects with constants.TykBootstrapReleaseLabel label that has
// constants.TykBootstrapDashboardSvcLabel value and gets this Service's metadata name, and port and
// updates DashboardSvc and DashboardPort fields.
func discoverDashboardSvc() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	ls := metav1.LabelSelector{MatchLabels: map[string]string{
		constants.TykBootstrapLabel: constants.TykBootstrapDashboardSvcLabel,
	}}
	if AppConfig.ReleaseName != "" {
		ls.MatchLabels[constants.TykBootstrapReleaseLabel] = AppConfig.ReleaseName
	}

	l := labels.Set(ls.MatchLabels).String()

	services, err := c.
		CoreV1().
		Services(AppConfig.TykPodNamespace).
		List(context.TODO(), metav1.ListOptions{LabelSelector: l})
	if err != nil {
		return err
	}

	if len(services.Items) == 0 {
		return fmt.Errorf("failed to find services with label %v\n", l)
	}

	if len(services.Items) > 1 {
		fmt.Printf("[WARNING] Found multiple services with label %v\n", l)
	}

	service := services.Items[0]
	if len(service.Spec.Ports) == 0 {
		return fmt.Errorf("svc/%v/%v has no open ports\n", service.Name, service.Namespace)
	}
	if len(service.Spec.Ports) > 1 {
		fmt.Printf("[WARNING] Found multiple open ports in svc/%v/%v\n", service.Name, service.Namespace)
	}

	if len(service.Spec.Ports) > 0 {
		AppConfig.DashboardPort = service.Spec.Ports[0].Port
	}
	AppConfig.DashboardSvc = service.Name

	return nil
}