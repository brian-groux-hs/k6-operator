package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/grafana/k6-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/labels"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Mark k6 as finished as jobs finish
func DeleteJobs(ctx context.Context, log logr.Logger, k6 *v1alpha1.K6, r *K6Reconciler) (ctrl.Result, error) {
	selector := labels.SelectorFromSet(map[string]string{
		"app":   "k6",
		"k6_cr": k6.Name,
	})

	opts := &client.ListOptions{LabelSelector: selector, Namespace: k6.Namespace}
	jl := &batchv1.JobList{}

	if err := r.List(ctx, jl, opts); err != nil {
		log.Error(err, "Could not list jobs")
		return ctrl.Result{}, err
	}

	k6.Status.Stage = "deleted"
	for _, job := range jl.Items {
		r.Delete(ctx, &job)
	}

	return ctrl.Result{}, nil
}
