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

package relocatingwebhook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type loggingWebhookHandler struct {
	Handler ExtendedHandler
	Log     logr.Logger
	Debug   bool
}

func NewLoggingWebhookHandler(handler ExtendedHandler, log logr.Logger, debug bool) *loggingWebhookHandler {
	return &loggingWebhookHandler{Handler: handler, Log: log, Debug: debug}
}

var _ ExtendedHandler = &loggingWebhookHandler{}

func (l *loggingWebhookHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	resp := l.Handler.Handle(ctx, req)
	if l.Debug {
		reqBs, _ := json.Marshal(req)
		respBs, _ := json.Marshal(resp)
		l.Log.Info(fmt.Sprintf("req: %s resp: %s\n", reqBs, respBs))
	}
	return resp
}

func (l *loggingWebhookHandler) InjectClient(c client.Client) error {
	return l.Handler.InjectClient(c)
}

func (l *loggingWebhookHandler) InjectDecoder(d *admission.Decoder) error {
	return l.Handler.InjectDecoder(d)
}