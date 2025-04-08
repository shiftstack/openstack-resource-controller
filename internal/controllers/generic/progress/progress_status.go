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
	"github.com/k-orc/openstack-resource-controller/internal/logging"
	orcerrors "github.com/k-orc/openstack-resource-controller/internal/util/errors"
)

type ReconcileStatus interface {
	ProgressMessages() []string

	EphemeralError() error
	TerminalError() *orcerrors.TerminalError

	Requeue() time.Duration
	Return(logr.Logger) (ctrl.Result, error)

	IsSetNotAvailable() bool
	NeedsReschedule() (bool, error)

	WithProgressMessage(...string) ReconcileStatus
	WithRequeue(time.Duration) ReconcileStatus
	WithSetNotAvailable() ReconcileStatus
	WithError(error) ReconcileStatus

	WithReconcileStatus(ReconcileStatus) ReconcileStatus
}

type reconcileStatus struct {
	messages        []string
	requeue         time.Duration
	setNotAvailable bool

	ephemeralError error
	terminalError  *orcerrors.TerminalError
}

var _ ReconcileStatus = &reconcileStatus{}

func NewReconcileStatus() ReconcileStatus {
	return (*reconcileStatus)(nil)
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

func (r *reconcileStatus) IsError() bool {
	return r.ephemeralError != nil || r.terminalError != nil
}

func (r *reconcileStatus) Return(log logr.Logger) (ctrl.Result, error) {
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
	if err == nil {
		return r
	}

	if r == nil {
		r = &reconcileStatus{}
	}

	if terminalError := (*orcerrors.TerminalError)(nil); errors.As(err, &terminalError) {
		r.terminalError = terminalError
	} else {
		r.ephemeralError = errors.Join(r.ephemeralError, err)
	}

	return r
}

func (r *reconcileStatus) WithReconcileStatus(o ReconcileStatus) ReconcileStatus {
	if r == nil {
		return o
	}

	r.WithProgressMessage(o.ProgressMessages()...).
		WithRequeue(o.Requeue()).
		WithError(o.EphemeralError()).
		WithError(o.TerminalError())
	if o.IsSetNotAvailable() {
		r.WithSetNotAvailable()
	}

	return r
}

type WaitingOnEvent int

const (
	WaitingOnCreation WaitingOnEvent = iota
	WaitingOnUpdate
	WaitingOnReady
	WaitingOnDeletion
)

func WaitingOnObject(r ReconcileStatus, kind, name string, waitingOn WaitingOnEvent) ReconcileStatus {
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
	r.WithProgressMessage(fmt.Sprintf("Waiting for %s/%s to be %s", kind, name, outcome))
	return r
}

func WaitingOnFinalizer(r ReconcileStatus, finalizer string) ReconcileStatus {
	r.WithProgressMessage(fmt.Sprintf("Waiting for finalizer %s to be removed", finalizer))
	r.WithProgressMessage()
	return r
}

func WaitingOnOpenStack(r ReconcileStatus, waitingOn WaitingOnEvent, pollingPeriod time.Duration) ReconcileStatus {
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

	r.WithProgressMessage(fmt.Sprintf("Waiting for OpenStack resource to be %s", outcome)).
		WithRequeue(pollingPeriod)
	return r
}

func NeedsRefresh(r ReconcileStatus) ReconcileStatus {
	r.WithProgressMessage("Resource status will be refreshed")
	return r
}
