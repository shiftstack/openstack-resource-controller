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

package progress

import (
	"errors"
	"fmt"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/go-logr/logr"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
)

type ReconcileStatus = *reconcileStatus

type reconcileStatus struct {
	messages        []string
	requeue         time.Duration
	setNotAvailable bool

	ephemeralError error
	terminalError  *orcerrors.TerminalError
}

func NewReconcileStatus() ReconcileStatus {
	return nil
}

func WrapError(err error) ReconcileStatus {
	return NewReconcileStatus().WithError(err)
}

func (r *reconcileStatus) ProgressMessages() []string {
	if r == nil {
		return nil
	}

	return r.messages
}

func (r *reconcileStatus) EphemeralError() error {
	if r == nil {
		return nil
	}

	return r.ephemeralError
}

func (r *reconcileStatus) TerminalError() *orcerrors.TerminalError {
	if r == nil {
		return nil
	}

	return r.terminalError
}

func (r *reconcileStatus) Requeue() time.Duration {
	if r == nil {
		return 0
	}

	return r.requeue
}

func (r *reconcileStatus) IsSetNotAvailable() bool {
	if r == nil {
		return false
	}

	return r.setNotAvailable
}

func (r *reconcileStatus) NeedsReschedule() (bool, error) {
	if r == nil {
		return false, nil
	}

	err := errors.Join(r.ephemeralError, r.terminalError)
	return len(r.messages) > 0 || err != nil, err
}

func (r *reconcileStatus) Return(log logr.Logger) (ctrl.Result, error) {
	if r == nil {
		return ctrl.Result{}, nil
	}

	if r.terminalError != nil {
		err := errors.Join(r.ephemeralError, r.terminalError)
		log.V(logging.Info).Info("not scheduling further reconciles for terminal error", "err", err.Error())
		return ctrl.Result{}, nil
	}

	if r.ephemeralError != nil {
		return ctrl.Result{}, r.ephemeralError
	}

	return ctrl.Result{RequeueAfter: r.requeue}, nil
}

func (r *reconcileStatus) WithProgressMessage(msgs ...string) ReconcileStatus {
	if len(msgs) == 0 {
		return r
	}

	if r == nil {
		r = &reconcileStatus{}
	}

	r.messages = append(r.messages, msgs...)
	return r
}

func (r *reconcileStatus) WithRequeue(requeue time.Duration) ReconcileStatus {
	if requeue == 0 {
		return r
	}

	if r == nil {
		r = &reconcileStatus{}
	}

	if r.requeue == 0 || requeue < r.requeue {
		r.requeue = requeue
	}
	return r
}

func (r *reconcileStatus) WithSetNotAvailable() ReconcileStatus {
	if r == nil {
		r = &reconcileStatus{}
	}

	r.setNotAvailable = true
	return r
}

func (r *reconcileStatus) WithError(err error) ReconcileStatus {
	fmt.Printf("\n\n\nerr: %+v\nerrT: %T\nerr == nil: %v", err, err, err == nil)

	if t, ok := err.(*orcerrors.TerminalError); ok {
		return r.WithTerminalError(t)
	}

	if err == nil {
		return r
	}

	if r == nil {
		r = &reconcileStatus{}
	}

	r.ephemeralError = errors.Join(r.ephemeralError, err)
	return r
}

func (r *reconcileStatus) WithTerminalError(err *orcerrors.TerminalError) ReconcileStatus {
	if err == nil {
		return r
	}

	if r == nil {
		r = &reconcileStatus{}
	}

	fmt.Printf("\n\n\nSetting terminal error from '%+v' to '%+v', err = '%+v'", r.terminalError, err, err)
	r.terminalError = err
	return r
}

func (r *reconcileStatus) WithReconcileStatus(o ReconcileStatus) ReconcileStatus {
	if r == nil {
		return o
	}

	if o.IsSetNotAvailable() {
		r = r.WithSetNotAvailable()
	}

	return r.WithProgressMessage(o.ProgressMessages()...).
		WithRequeue(o.Requeue()).
		WithError(o.EphemeralError()).
		WithError(o.TerminalError())
}

type WaitingOnEvent int

const (
	WaitingOnCreation WaitingOnEvent = iota
	WaitingOnUpdate
	WaitingOnReady
	WaitingOnDeletion
)

func WaitingOnObject(kind, name string, waitingOn WaitingOnEvent) ReconcileStatus {
	return NewReconcileStatus().WaitingOnObject(kind, name, waitingOn)
}

func (r *reconcileStatus) WaitingOnObject(kind, name string, waitingOn WaitingOnEvent) ReconcileStatus {
	var outcome string
	switch waitingOn {
	case WaitingOnCreation:
		outcome = "created"
	case WaitingOnUpdate:
		outcome = "updated"
	case WaitingOnReady:
		outcome = "ready"
	case WaitingOnDeletion:
		outcome = "deleted"
	}
	return r.WithProgressMessage(fmt.Sprintf("Waiting for %s/%s to be %s", kind, name, outcome))
}

func WaitingOnFinalizer(finalizer string) ReconcileStatus {
	return NewReconcileStatus().WaitingOnFinalizer(finalizer)
}

func (r *reconcileStatus) WaitingOnFinalizer(finalizer string) ReconcileStatus {
	return r.WithProgressMessage(fmt.Sprintf("Waiting for finalizer %s to be removed", finalizer))
}

func WaitingOnOpenStack(waitingOn WaitingOnEvent, pollingPeriod time.Duration) ReconcileStatus {
	return NewReconcileStatus().WaitingOnOpenStack(waitingOn, pollingPeriod)
}

func (r *reconcileStatus) WaitingOnOpenStack(waitingOn WaitingOnEvent, pollingPeriod time.Duration) ReconcileStatus {
	var outcome string
	switch waitingOn {
	case WaitingOnCreation:
		outcome = "created externally"
	case WaitingOnUpdate:
		outcome = "updated"
	case WaitingOnReady:
		outcome = "ready"
	case WaitingOnDeletion:
		outcome = "deleted"
	}

	return r.WithProgressMessage(fmt.Sprintf("Waiting for OpenStack resource to be %s", outcome)).
		WithRequeue(pollingPeriod)
}

func NeedsRefresh() ReconcileStatus {
	return NewReconcileStatus().NeedsRefresh()
}

func (r *reconcileStatus) NeedsRefresh() ReconcileStatus {
	return r.WithProgressMessage("Resource status will be refreshed")
}
