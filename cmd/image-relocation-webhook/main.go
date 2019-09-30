/*
 * Copyright (c) 2019-Present Pivotal Software, Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	webhookv1alpha1 "github.com/pivotal/image-relocation/pkg/api/v1alpha1"
	"github.com/pivotal/image-relocation/pkg/multimap"

	"github.com/pivotal/image-relocation/pkg/controller"
	"github.com/pivotal/image-relocation/pkg/relocatingwebhook"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.Parse()

	log := logf.Log.WithName("image-relocation")
	logf.SetLogger(zap.New())
	entryLog := log.WithName("main")

	if debug {
		entryLog.Info("debug logging enabled")
	}

	entryLog.Info("setting up scheme")
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		entryLog.Error(err, "unable to add client-go types to scheme")
		os.Exit(1)
	}
	if err := webhookv1alpha1.AddToScheme(scheme); err != nil {
		entryLog.Error(err, "unable to add webhook types to scheme")
		os.Exit(1)
	}

	entryLog.Info("setting up controller manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{Scheme:scheme})
	if err != nil {
		entryLog.Error(err, "unable to set up controller manager")
		os.Exit(1)
	}

	stopCh := signals.SetupSignalHandler()
	comp := multimap.New(stopCh)

	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	entryLog.Info("registering webhook to the webhook server")
	hookServer.Register("/image-relocation", &webhook.Admission{
		Handler: relocatingwebhook.NewLoggingWebhookHandler(relocatingwebhook.NewImageReferenceRelocator(comp), log.WithName("handler"), debug),
	})

	if err = (&controller.ImageMapReconciler{
		Client: mgr.GetClient(),
		Log:    log.WithName("controller").WithName("ImageMap"),
		Map: comp,
	}).SetupWithManager(mgr); err != nil {
		entryLog.Error(err, "unable to create controller", "controller", "ImageMap")
		os.Exit(1)
	}

	entryLog.Info("starting controller manager")
	if err := mgr.Start(stopCh); err != nil {
		entryLog.Error(err, "unable to start controller manager")
		os.Exit(1)
	}
}
