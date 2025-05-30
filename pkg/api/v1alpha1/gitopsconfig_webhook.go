package v1alpha1

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var gitopsconfiglog = logf.Log.WithName("gitopsconfig-resource")

func (r *GitOpsConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(r).
		Complete()
}

// Change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-anywhere-eks-amazonaws-com-v1alpha1-gitopsconfig,mutating=false,failurePolicy=fail,sideEffects=None,groups=anywhere.eks.amazonaws.com,resources=gitopsconfigs,verbs=create;update,versions=v1alpha1,name=validation.gitopsconfig.anywhere.amazonaws.com,admissionReviewVersions={v1,v1beta1}

var _ webhook.CustomValidator = &GitOpsConfig{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *GitOpsConfig) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	gitopsConfig, ok := obj.(*GitOpsConfig)
	if !ok {
		return nil, fmt.Errorf("expected a GitOpsConfig but got %T", obj)
	}

	gitopsconfiglog.Info("validate create", "name", gitopsConfig.Name)

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *GitOpsConfig) ValidateUpdate(_ context.Context, old, obj runtime.Object) (admission.Warnings, error) {
	gitopsConfig, ok := obj.(*GitOpsConfig)
	if !ok {
		return nil, fmt.Errorf("expected a GitOpsConfig but got %T", obj)
	}

	gitopsconfiglog.Info("validate update", "name", gitopsConfig.Name)

	oldGitOpsConfig, ok := old.(*GitOpsConfig)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a GitOpsConfig but got a %T", old))
	}

	var allErrs field.ErrorList

	allErrs = append(allErrs, validateImmutableGitOpsFields(gitopsConfig, oldGitOpsConfig)...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(GroupVersion.WithKind(GitOpsConfigKind).GroupKind(), gitopsConfig.Name, allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *GitOpsConfig) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	gitopsConfig, ok := obj.(*GitOpsConfig)
	if !ok {
		return nil, fmt.Errorf("expected a GitOpsConfig but got %T", obj)
	}

	gitopsconfiglog.Info("validate delete", "name", gitopsConfig.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

func validateImmutableGitOpsFields(new, old *GitOpsConfig) field.ErrorList {
	var allErrs field.ErrorList

	if !new.Spec.Equal(&old.Spec) {
		allErrs = append(
			allErrs,
			field.Forbidden(field.NewPath(GitOpsConfigKind), "config is immutable"),
		)
	}

	return allErrs
}
