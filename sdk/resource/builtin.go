// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/label"
	opentelemetry "go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/semconv"
)

type (
	// TelemetrySDK is a Detector that provides information about
	// the OpenTelemetry SDK used.  This Detector is included as a
	// builtin. If these resource attributes are not wanted, use
	// the WithTelemetrySDK(nil) or WithoutBuiltin() options to
	// explicitly disable them.
	TelemetrySDK struct{}

	// Host is a Detector that provides information about the host
	// being run on. This Detector is included as a builtin. If
	// these resource attributes are not wanted, use the
	// WithHost(nil) or WithoutBuiltin() options to explicitly
	// disable them.
	Host struct{}

	stringDetector struct {
		K label.Key
		F func() (string, error)
	}
)

var (
	_ Detector = TelemetrySDK{}
	_ Detector = Host{}
	_ Detector = stringDetector{}
)

// Detect returns a *Resource that describes the OpenTelemetry SDK used.
func (TelemetrySDK) Detect(context.Context) (*Resource, error) {
	return NewWithAttributes(
		semconv.TelemetrySDKNameKey.String("opentelemetry-go"),
		semconv.TelemetrySDKLanguageKey.String("go"),
		semconv.TelemetrySDKVersionKey.String(opentelemetry.Version()),
	), nil
}

// Detect returns a *Resource that describes the host being run on.
func (Host) Detect(ctx context.Context) (*Resource, error) {
	return StringDetector(semconv.HostNameKey, os.Hostname).Detect(ctx)
}

// StringDetector returns a Detector that will produce a *Resource
// containing the string as a value corresponding to k.
func StringDetector(k label.Key, f func() (string, error)) Detector {
	return stringDetector{K: k, F: f}
}

// Detect implements Detector.
func (sd stringDetector) Detect(ctx context.Context) (*Resource, error) {
	value, err := sd.F()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", string(sd.K), err)
	}
	return NewWithAttributes(sd.K.String(value)), nil
}
