/*
Copyright 2024, Tax Administration of The Netherlands.
Licensed under the EUPL 1.2.
See LICENSE.md for details.
*/

package v1alpha1

import (
	"context"
	"fmt"

	"github.com/belastingdienst/opr-paas/api/v1alpha1"
	"github.com/belastingdienst/opr-paas/internal/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func setRequestLogger(ctx context.Context, obj client.Object) (context.Context, *zerolog.Logger) {
	logger := log.With().
		Any("webhook", obj.GetObjectKind().GroupVersionKind()).
		Dict("object", zerolog.Dict().
			Str("name", obj.GetName()).
			Str("namespace", obj.GetNamespace()),
		).
		Str("requestId", uuid.NewString()).
		Logger()

	return logger.WithContext(ctx), &logger
}

// SetupPaasWebhookWithManager registers the webhook for Paas in the manager.
func SetupPaasWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&v1alpha1.Paas{}).
		WithValidator(&PaasCustomValidator{client: mgr.GetClient()}).
		Complete()
}

// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-cpet-belastingdienst-nl-v1alpha1-paas,mutating=false,failurePolicy=fail,sideEffects=None,groups=cpet.belastingdienst.nl,resources=paas,verbs=create;update,versions=v1alpha1,name=vpaas-v1alpha1.kb.io,admissionReviewVersions=v1

// PaasCustomValidator struct is responsible for validating the Paas resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
// +kubebuilder:object:generate=false
type PaasCustomValidator struct {
	client client.Client
}

var _ webhook.CustomValidator = &PaasCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Paas.
func (v *PaasCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	paas, ok := obj.(*v1alpha1.Paas)
	if !ok {
		return nil, fmt.Errorf("expected a Paas object but got %T", obj)
	}
	_, logger := setRequestLogger(ctx, paas)
	logger.Info().Msg("starting validation webhook for creation")

	var allErrs field.ErrorList
	conf := config.GetConfig()
	if err := v.validateCaps(conf, paas.Spec.Capabilities); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: v1alpha1.GroupVersion.Group, Kind: "Paas"},
		paas.Name,
		allErrs,
	)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Paas.
func (v *PaasCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	paas, ok := newObj.(*v1alpha1.Paas)
	if !ok {
		return nil, fmt.Errorf("expected a Paas object for the newObj but got %T", newObj)
	}
	_, logger := setRequestLogger(ctx, paas)
	logger.Info().Msg("starting validation webhook for update")

	var allErrs field.ErrorList
	conf := config.GetConfig()
	if err := v.validateCaps(conf, paas.Spec.Capabilities); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: v1alpha1.GroupVersion.Group, Kind: "Paas"},
		paas.Name,
		allErrs,
	)
}

// TODO(portly-halicore-76): determine whether this can be left out
// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Paas.
func (v *PaasCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	paas, ok := obj.(*v1alpha1.Paas)
	if !ok {
		return nil, fmt.Errorf("expected a Paas object but got %T", obj)
	}
	_, logger := setRequestLogger(ctx, paas)
	logger.Info().Msg("starting validation webhook for deletion")

	// TODO(portly-halicore-76): fill in your validation logic upon object deletion.

	return nil, nil
}

// validateCaps returns an error if any of the passed capabilities is not configured.
func (v *PaasCustomValidator) validateCaps(conf v1alpha1.PaasConfigSpec, caps v1alpha1.PaasCapabilities) []*field.Error {
	errs := make([]*field.Error, 0)

	for name := range caps {
		if _, ok := conf.Capabilities[name]; !ok {
			errs = append(errs, field.Invalid(
				field.NewPath("spec").Child("capabilities"),
				name,
				"capability not configured",
			))
		}
	}

	return errs
}
