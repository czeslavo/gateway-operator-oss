package extensions

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kcfgkonnect "github.com/kong/kubernetes-configuration/api/konnect"
	konnectv1alpha1 "github.com/kong/kubernetes-configuration/api/konnect/v1alpha1"
)

// ValidateExtensions validates the extensions referenced by the provided DataPlane and returns
// a condition indicating if the extensions are supported or not.
func validateExtensions[t ExtendableT](obj t) *metav1.Condition {
	if len(obj.GetExtensions()) == 0 {
		return nil
	}

	condition := metav1.Condition{
		Status:             metav1.ConditionTrue,
		Type:               string(kcfgkonnect.AcceptedExtensionsType),
		Reason:             string(kcfgkonnect.AcceptedExtensionsReason),
		Message:            "All extensions are accepted",
		ObservedGeneration: obj.GetGeneration(),
		LastTransitionTime: metav1.Now(),
	}
	var messageBuilder strings.Builder
	for i, ext := range obj.GetExtensions() {
		if ext.Group != konnectv1alpha1.SchemeGroupVersion.Group || ext.Kind != konnectv1alpha1.KonnectExtensionKind {
			buildMessage(&messageBuilder, fmt.Sprintf("Extension %s/%s is not supported", ext.Group, ext.Kind))
			continue
		}
		for j, ext2 := range obj.GetExtensions() {
			if i != j {
				if ext.Group == ext2.Group && ext.Kind == ext2.Kind {
					message := fmt.Sprintf("Extension %s/%s is duplicated", ext.Group, ext.Kind)
					if !strings.Contains(messageBuilder.String(), message) {
						buildMessage(&messageBuilder, message)
					}
				}
			}
		}
	}
	if messageBuilder.Len() > 0 {
		condition.Status = metav1.ConditionFalse
		condition.Reason = string(kcfgkonnect.NotSupportedExtensionsReason)
		condition.Message = messageBuilder.String()
	}

	return &condition
}

func buildMessage(messageBuilder *strings.Builder, message string) {
	if messageBuilder.Len() > 0 {
		messageBuilder.WriteString(" - ")
	}
	messageBuilder.WriteString(message)
}
