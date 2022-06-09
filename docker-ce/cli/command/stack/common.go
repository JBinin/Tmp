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
package stack

import (
	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/compose/convert"
	"github.com/docker/docker/client"
	"github.com/docker/docker/opts"
)

func getStackFilter(namespace string) filters.Args {
	filter := filters.NewArgs()
	filter.Add("label", convert.LabelNamespace+"="+namespace)
	return filter
}

func getStackFilterFromOpt(namespace string, opt opts.FilterOpt) filters.Args {
	filter := opt.Value()
	filter.Add("label", convert.LabelNamespace+"="+namespace)
	return filter
}

func getAllStacksFilter() filters.Args {
	filter := filters.NewArgs()
	filter.Add("label", convert.LabelNamespace)
	return filter
}

func getServices(
	ctx context.Context,
	apiclient client.APIClient,
	namespace string,
) ([]swarm.Service, error) {
	return apiclient.ServiceList(
		ctx,
		types.ServiceListOptions{Filters: getStackFilter(namespace)})
}

func getStackNetworks(
	ctx context.Context,
	apiclient client.APIClient,
	namespace string,
) ([]types.NetworkResource, error) {
	return apiclient.NetworkList(
		ctx,
		types.NetworkListOptions{Filters: getStackFilter(namespace)})
}

func getStackSecrets(
	ctx context.Context,
	apiclient client.APIClient,
	namespace string,
) ([]swarm.Secret, error) {
	return apiclient.SecretList(
		ctx,
		types.SecretListOptions{Filters: getStackFilter(namespace)})
}
