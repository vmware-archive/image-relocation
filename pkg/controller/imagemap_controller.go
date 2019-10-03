/*

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

package controller

import (
	"context"

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ir "github.com/pivotal/image-relocation/pkg/api/v1alpha1"
	webhookv1alpha1 "github.com/pivotal/image-relocation/pkg/api/v1alpha1"
	"github.com/pivotal/image-relocation/pkg/multimap"
)

// ImageMapReconciler reconciles a ClusterImageMap object
type ImageMapReconciler struct {
	client.Client
	Log logr.Logger
	Map multimap.Composite
}

func (r *ImageMapReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("ClusterImageMap", req.NamespacedName)

	var imageMap ir.ClusterImageMap
	if err := r.Get(ctx, req.NamespacedName, &imageMap); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("deleting")
			_ = r.Map.Delete(req.NamespacedName.String()) // ignore error in case it has already been deleted
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to get ClusterImageMap")
		return ctrl.Result{}, err
	}

	log.Info("adding")
	if err := r.Map.Add(req.NamespacedName.String(), imageMap.Spec.Map); err != nil {
		log.Error(err, "unable to add ClusterImageMap")
		return ctrl.Result{
			Requeue: false, // or true for a soft failure FIXME
		}, err
	}

	return ctrl.Result{}, nil
}

func (r *ImageMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webhookv1alpha1.ClusterImageMap{}).
		Complete(r)
}
