package v1

import (
	"fmt"

	"github.com/caicloud/clientset/listerfactory"
	releaseapi "github.com/caicloud/clientset/pkg/apis/release/v1alpha1"
	"github.com/caicloud/clientset/util/event"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	ssetErrorEventCases = []event.EventCase{
		// match all FailedCreate
		{EventType: corev1.EventTypeWarning, Reason: event.FailedCreatePodReason},
	}
)

func JudgeStatefulSet(factory listerfactory.ListerFactory, obj runtime.Object) (releaseapi.ResourceStatus, error) {
	statefulset, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		return releaseapi.ResourceStatusFrom(""), fmt.Errorf("unknown type for statefulset: %s", obj.GetObjectKind().GroupVersionKind().String())
	}
	if statefulset == nil {
		return releaseapi.ResourceStatusFrom(""), fmt.Errorf("statefulset can not be nil")
	}

	lr, err := newLongRunning(factory, statefulset)
	if err != nil {
		return releaseapi.ResourceStatusFrom(""), err
	}
	return lr.Judge()
}

type statefulsetLongRunning struct {
	statefulset     *appsv1.StatefulSet
	updatedRevision string
}

func newStatefulSetLongRunning(statefulset *appsv1.StatefulSet) LongRunning {
	return &statefulsetLongRunning{
		statefulset: statefulset,
	}
}

func (d *statefulsetLongRunning) PredictUpdatedRevision(factory listerfactory.ListerFactory, events []*corev1.Event) (*releaseapi.ResourceStatus, error) {
	statefulset := d.statefulset
	d.updatedRevision = statefulset.Status.UpdateRevision
	if d.updatedRevision == "" {
		return nil, ErrUpdatedRevisionNotExists
	}

	return nil, nil
}

func (d *statefulsetLongRunning) IsUpdatedPod(pod *corev1.Pod) bool {
	if d.updatedRevision == "" {
		return false
	}

	return getLabel(pod, appsv1.StatefulSetRevisionLabel) == d.updatedRevision
}

func (d *statefulsetLongRunning) PredictEvents(events []*corev1.Event) (*releaseapi.ResourceStatus, *corev1.Event) {
	lastEvent := getLatestEventFor(d.statefulset.GroupVersionKind().Kind, d.statefulset, events)
	for _, c := range ssetErrorEventCases {
		if c.Match(lastEvent) {
			return &releaseapi.ResourceStatus{
				Phase:   releaseapi.ResourceFailed,
				Reason:  lastEvent.Reason,
				Message: lastEvent.Message,
			}, lastEvent
		}
	}
	return nil, lastEvent
}

func (d *statefulsetLongRunning) DesiredReplics() int32 {
	// statefulset has no desired replicas, its value should always be 0
	return *d.statefulset.Spec.Replicas
}
