/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
package convert

import (
	swarmtypes "github.com/docker/docker/api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
)

// SecretFromGRPC converts a grpc Secret to a Secret.
func SecretFromGRPC(s *swarmapi.Secret) swarmtypes.Secret {
	secret := swarmtypes.Secret{
		ID: s.ID,
		Spec: swarmtypes.SecretSpec{
			Annotations: annotationsFromGRPC(s.Spec.Annotations),
			Data:        s.Spec.Data,
		},
	}

	secret.Version.Index = s.Meta.Version.Index
	// Meta
	secret.CreatedAt, _ = gogotypes.TimestampFromProto(s.Meta.CreatedAt)
	secret.UpdatedAt, _ = gogotypes.TimestampFromProto(s.Meta.UpdatedAt)

	return secret
}

// SecretSpecToGRPC converts Secret to a grpc Secret.
func SecretSpecToGRPC(s swarmtypes.SecretSpec) swarmapi.SecretSpec {
	return swarmapi.SecretSpec{
		Annotations: swarmapi.Annotations{
			Name:   s.Name,
			Labels: s.Labels,
		},
		Data: s.Data,
	}
}

// SecretReferencesFromGRPC converts a slice of grpc SecretReference to SecretReference
func SecretReferencesFromGRPC(s []*swarmapi.SecretReference) []*swarmtypes.SecretReference {
	refs := []*swarmtypes.SecretReference{}

	for _, r := range s {
		ref := &swarmtypes.SecretReference{
			SecretID:   r.SecretID,
			SecretName: r.SecretName,
		}

		if t, ok := r.Target.(*swarmapi.SecretReference_File); ok {
			ref.File = &swarmtypes.SecretReferenceFileTarget{
				Name: t.File.Name,
				UID:  t.File.UID,
				GID:  t.File.GID,
				Mode: t.File.Mode,
			}
		}

		refs = append(refs, ref)
	}

	return refs
}
