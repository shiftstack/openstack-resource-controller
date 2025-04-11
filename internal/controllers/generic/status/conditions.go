/*
Copyright 2024 The ORC Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package status

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyconfigv1 "k8s.io/client-go/applyconfigurations/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/applyconfigs"
)

type WithConditionsApplyConfiguration[T any] interface {
	WithConditions(...*applyconfigv1.ConditionApplyConfiguration) T
}

func SetCommonConditions[T any](
	orcObject orcv1alpha1.ObjectWithConditions,
	applyConfig WithConditionsApplyConfiguration[T],
	availableStatus metav1.ConditionStatus,
	reconcileStatus progress.ReconcileStatus,
	now metav1.Time,
) {
	availableCondition := applyconfigv1.Condition().
		WithType(orcv1alpha1.ConditionAvailable).
		WithStatus(availableStatus).
		WithObservedGeneration(orcObject.GetGeneration())
	progressingCondition := applyconfigv1.Condition().
		WithType(orcv1alpha1.ConditionProgressing).
		WithObservedGeneration(orcObject.GetGeneration())

	// We are Progressing iff we are anticipating being reconciled again. This
	// means one of:
	// - err contains a non-terminal error, so we expect an error backoff
	// - reconcileStatus does not indicate that we are waiting on some condition

	fmt.Printf("\n\n\nreconcileStatus: %+v", reconcileStatus)
	fmt.Printf("\n\n\nfoo 1")
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); !needsReschedule {
		fmt.Printf("\n\n\nfoo 2")
		progressingCondition.
			WithStatus(metav1.ConditionFalse).
			WithReason(orcv1alpha1.ConditionReasonSuccess).
			WithMessage("OpenStack resource is up to date")
	} else {
		fmt.Printf("\n\n\nfoo 3")
		if err := reconcileStatus.EphemeralError(); err != nil {
			fmt.Printf("\n\n\nfoo 4")
			progressingCondition.
				WithStatus(metav1.ConditionTrue).
				WithReason(orcv1alpha1.ConditionReasonTransientError).
				WithMessage(err.Error())
		} else if err := reconcileStatus.TerminalError(); err != nil {
			fmt.Printf("\n\n\nfoo 5")
			progressingCondition.
				WithStatus(metav1.ConditionFalse).
				WithReason(err.Reason).
				WithMessage(err.Message)
		} else {
			fmt.Printf("\n\n\nfoo 6")
			progressingCondition.
				WithStatus(metav1.ConditionTrue).
				WithReason(orcv1alpha1.ConditionReasonProgressing).
				WithMessage(strings.Join(reconcileStatus.ProgressMessages(), "\n"))
		}
	}

	if availableStatus == metav1.ConditionTrue {
		availableCondition.
			WithReason(orcv1alpha1.ConditionReasonSuccess).
			WithMessage("OpenStack resource is available")
	} else {
		// Copy reason and message from progressing
		availableCondition.
			WithReason(*progressingCondition.Reason).
			WithMessage(*progressingCondition.Message)
	}

	// Maintain condition timestamps if they haven't changed
	// This also ensures that we don't generate an update event if nothing has changed
	for _, condition := range []*applyconfigv1.ConditionApplyConfiguration{availableCondition, progressingCondition} {
		previous := meta.FindStatusCondition(orcObject.GetConditions(), *condition.Type)
		if previous != nil && applyconfigs.ConditionsEqual(previous, condition) {
			condition.WithLastTransitionTime(previous.LastTransitionTime)
		} else {
			condition.WithLastTransitionTime(now)
		}
	}

	applyConfig.WithConditions(availableCondition, progressingCondition)
}
