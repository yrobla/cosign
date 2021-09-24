//
// Copyright 2021 The Sigstore Authors.
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

package static

import (
	"github.com/go-openapi/swag"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/sigstore/cosign/pkg/oci"
	ctypes "github.com/sigstore/cosign/pkg/types"
)

// Option is a functional option for customizing static signatures.
type Option func(*options)

type options struct {
	MediaType   types.MediaType
	Bundle      *oci.Bundle
	Cert        []byte
	Chain       []byte
	Annotations map[string]string
}

func makeOptions(opts ...Option) (*options, error) {
	o := &options{
		MediaType:   ctypes.SimpleSigningMediaType,
		Annotations: make(map[string]string),
	}

	for _, opt := range opts {
		opt(o)
	}

	if o.Cert != nil {
		o.Annotations[CertificateAnnotationKey] = string(o.Cert)
		o.Annotations[ChainAnnotationKey] = string(o.Chain)
	}

	if o.Bundle != nil {
		b, err := swag.WriteJSON(o.Bundle)
		if err != nil {
			return nil, err
		}
		o.Annotations[BundleAnnotationKey] = string(b)
	}

	return o, nil
}

// WithMediaType sets the media type of the signature.
func WithMediaType(mt types.MediaType) Option {
	return func(o *options) {
		o.MediaType = mt
	}
}

// WithAnnotations sets the annotations that will be associated.
func WithAnnotations(ann map[string]string) Option {
	return func(o *options) {
		o.Annotations = ann
	}
}

// WithBundle sets the bundle to attach to the signature
func WithBundle(b *oci.Bundle) Option {
	return func(o *options) {
		o.Bundle = b
	}
}

// WithCertChain sets the certificate chain for this signature.
func WithCertChain(cert, chain []byte) Option {
	return func(o *options) {
		o.Cert = cert
		o.Chain = chain
	}
}