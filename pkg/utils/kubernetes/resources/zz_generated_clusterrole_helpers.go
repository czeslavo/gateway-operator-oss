// This file is generated by /hack/generators/kic-role-generator. DO NOT EDIT.

package resources

import (
	"fmt"

	"github.com/Masterminds/semver"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/kong/gateway-operator/internal/versions"
	"github.com/kong/gateway-operator/pkg/consts"
	"github.com/kong/gateway-operator/pkg/utils/kubernetes/resources/clusterroles"
)

// -----------------------------------------------------------------------------
// ClusterRole generator helper
// -----------------------------------------------------------------------------

// GenerateNewClusterRoleForControlPlane is a helper function that extract
// the version from the tag, and returns the ClusterRole with all the needed
// permissions.
func GenerateNewClusterRoleForControlPlane(controlplaneName string, image string) (*rbacv1.ClusterRole, error) {
	versionToUse := versions.DefaultControlPlaneVersion
	imageToUse := consts.DefaultControlPlaneImage
	var constraint *semver.Constraints

	if image != "" {
		v, err := versions.FromImage(image)
		if err != nil {
			return nil, err
		}
		supported, err := versions.IsControlPlaneImageVersionSupported(image)
		if err != nil {
			return nil, err
		}
		if supported {
			imageToUse = image
			versionToUse = v.String()
		}
	}

	semVersion, err := semver.NewVersion(versionToUse)
	if err != nil {
		return nil, err
	}

	constraint, err = semver.NewConstraint("<2.12, >=2.11")
	if err != nil {
		return nil, err
	}
	if constraint.Check(semVersion) {
		cr := clusterroles.GenerateNewClusterRoleForControlPlane_lt2_12_ge2_11(controlplaneName)
		LabelObjectAsControlPlaneManaged(cr)
		return cr, nil
	}

	constraint, err = semver.NewConstraint("<3.0, >=2.12")
	if err != nil {
		return nil, err
	}
	if constraint.Check(semVersion) {
		cr := clusterroles.GenerateNewClusterRoleForControlPlane_lt3_0_ge2_12(controlplaneName)
		LabelObjectAsControlPlaneManaged(cr)
		return cr, nil
	}

	constraint, err = semver.NewConstraint(">=3.0")
	if err != nil {
		return nil, err
	}
	if constraint.Check(semVersion) {
		cr := clusterroles.GenerateNewClusterRoleForControlPlane_ge3_0(controlplaneName)
		LabelObjectAsControlPlaneManaged(cr)
		return cr, nil
	}

	return nil, fmt.Errorf("version %s not supported", imageToUse)
}
