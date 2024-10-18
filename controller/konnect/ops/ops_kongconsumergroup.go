package ops

import (
	"context"
	"errors"
	"fmt"

	sdkkonnectcomp "github.com/Kong/sdk-konnect-go/models/components"
	sdkkonnectops "github.com/Kong/sdk-konnect-go/models/operations"
	sdkkonnecterrs "github.com/Kong/sdk-konnect-go/models/sdkerrors"
	"github.com/samber/lo"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	configurationv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
)

func createConsumerGroup(
	ctx context.Context,
	sdk ConsumerGroupSDK,
	group *configurationv1beta1.KongConsumerGroup,
) error {
	if group.GetControlPlaneID() == "" {
		return fmt.Errorf("can't create %T %s without a Konnect ControlPlane ID", group, client.ObjectKeyFromObject(group))
	}

	resp, err := sdk.CreateConsumerGroup(ctx,
		group.Status.Konnect.ControlPlaneID,
		kongConsumerGroupToSDKConsumerGroupInput(group),
	)
	// Can't adopt it as it will cause conflicts between the controller
	// that created that entity and already manages it.
	// TODO: implement entity adoption https://github.com/Kong/gateway-operator/issues/460
	if errWrap := wrapErrIfKonnectOpFailed(err, CreateOp, group); errWrap != nil {
		return errWrap
	}

	if resp == nil || resp.ConsumerGroup == nil || resp.ConsumerGroup.ID == nil || *resp.ConsumerGroup.ID == "" {
		return fmt.Errorf("failed creating %s: %w", group.GetTypeName(), ErrNilResponse)
	}

	id := *resp.ConsumerGroup.ID
	group.SetKonnectID(id)

	return nil
}

// updateConsumerGroup updates a KongConsumerGroup in Konnect.
// The KongConsumerGroup is assumed to have a Konnect ID set in status.
// It returns an error if the KongConsumerGroup does not have a ControlPlaneRef.
func updateConsumerGroup(
	ctx context.Context,
	sdk ConsumerGroupSDK,
	group *configurationv1beta1.KongConsumerGroup,
) error {
	cpID := group.GetControlPlaneID()
	if cpID == "" {
		return fmt.Errorf("can't update %T %s without a Konnect ControlPlane ID", group, client.ObjectKeyFromObject(group))
	}

	_, err := sdk.UpsertConsumerGroup(ctx,
		sdkkonnectops.UpsertConsumerGroupRequest{
			ControlPlaneID:  cpID,
			ConsumerGroupID: group.GetKonnectStatus().GetKonnectID(),
			ConsumerGroup:   kongConsumerGroupToSDKConsumerGroupInput(group),
		},
	)

	// Can't adopt it as it will cause conflicts between the controller
	// that created that entity and already manages it.
	// TODO: implement entity adoption https://github.com/Kong/gateway-operator/issues/460
	if errWrap := wrapErrIfKonnectOpFailed(err, UpdateOp, group); errWrap != nil {
		return errWrap
	}

	return nil
}

// deleteConsumerGroup deletes a KongConsumerGroup in Konnect.
// The KongConsumerGroup is assumed to have a Konnect ID set in status.
// It returns an error if the operation fails.
func deleteConsumerGroup(
	ctx context.Context,
	sdk ConsumerGroupSDK,
	consumer *configurationv1beta1.KongConsumerGroup,
) error {
	id := consumer.Status.Konnect.GetKonnectID()
	_, err := sdk.DeleteConsumerGroup(ctx, consumer.Status.Konnect.ControlPlaneID, id)
	if errWrap := wrapErrIfKonnectOpFailed(err, DeleteOp, consumer); errWrap != nil {
		// Consumer delete operation returns an SDKError instead of a NotFoundError.
		var sdkError *sdkkonnecterrs.SDKError
		if errors.As(errWrap, &sdkError) {
			if sdkError.StatusCode == 404 {
				ctrllog.FromContext(ctx).
					Info("entity not found in Konnect, skipping delete",
						"op", DeleteOp, "type", consumer.GetTypeName(), "id", id,
					)
				return nil
			}
			return FailedKonnectOpError[configurationv1beta1.KongConsumerGroup]{
				Op:  DeleteOp,
				Err: sdkError,
			}
		}
		return FailedKonnectOpError[configurationv1beta1.KongConsumerGroup]{
			Op:  DeleteOp,
			Err: errWrap,
		}
	}

	return nil
}

func kongConsumerGroupToSDKConsumerGroupInput(
	group *configurationv1beta1.KongConsumerGroup,
) sdkkonnectcomp.ConsumerGroupInput {
	return sdkkonnectcomp.ConsumerGroupInput{
		Tags: GenerateTagsForObject(group),
		Name: group.Spec.Name,
	}
}

// getConsumerGroupForUID lists consumer groups in Konnect with given k8s uid as its tag.
func getConsumerGroupForUID(
	ctx context.Context,
	sdk ConsumerGroupSDK,
	cg *configurationv1beta1.KongConsumerGroup,
) (string, error) {
	cpID := cg.GetControlPlaneID()

	reqList := sdkkonnectops.ListConsumerGroupRequest{
		// NOTE: only filter on object's UID.
		// Other fields like name might have changed in the meantime but that's OK.
		// Those will be enforced via subsequent updates.
		ControlPlaneID: cpID,
		Tags:           lo.ToPtr(UIDLabelForObject(cg)),
	}

	resp, err := sdk.ListConsumerGroup(ctx, reqList)
	if err != nil {
		return "", fmt.Errorf("failed listing %s: %w", cg.GetTypeName(), err)
	}
	if resp == nil || resp.Object == nil {
		return "", fmt.Errorf("failed listing %s: %w", cg.GetTypeName(), ErrNilResponse)
	}

	return getMatchingEntryFromListResponseData(sliceToEntityWithIDSlice(resp.Object.Data), cg)
}
