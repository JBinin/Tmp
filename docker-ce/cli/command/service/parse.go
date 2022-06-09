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
package service

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	swarmtypes "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// ParseSecrets retrieves the secrets from the requested names and converts
// them to secret references to use with the spec
func ParseSecrets(client client.SecretAPIClient, requestedSecrets []*types.SecretRequestOption) ([]*swarmtypes.SecretReference, error) {
	secretRefs := make(map[string]*swarmtypes.SecretReference)
	ctx := context.Background()

	for _, secret := range requestedSecrets {
		if _, exists := secretRefs[secret.Target]; exists {
			return nil, fmt.Errorf("duplicate secret target for %s not allowed", secret.Source)
		}
		secretRef := &swarmtypes.SecretReference{
			File: &swarmtypes.SecretReferenceFileTarget{
				Name: secret.Target,
				UID:  secret.UID,
				GID:  secret.GID,
				Mode: secret.Mode,
			},
			SecretName: secret.Source,
		}

		secretRefs[secret.Target] = secretRef
	}

	args := filters.NewArgs()
	for _, s := range secretRefs {
		args.Add("names", s.SecretName)
	}

	secrets, err := client.SecretList(ctx, types.SecretListOptions{
		Filters: args,
	})
	if err != nil {
		return nil, err
	}

	foundSecrets := make(map[string]string)
	for _, secret := range secrets {
		foundSecrets[secret.Spec.Annotations.Name] = secret.ID
	}

	addedSecrets := []*swarmtypes.SecretReference{}

	for _, ref := range secretRefs {
		id, ok := foundSecrets[ref.SecretName]
		if !ok {
			return nil, fmt.Errorf("secret not found: %s", ref.SecretName)
		}

		// set the id for the ref to properly assign in swarm
		// since swarm needs the ID instead of the name
		ref.SecretID = id
		addedSecrets = append(addedSecrets, ref)
	}

	return addedSecrets, nil
}
