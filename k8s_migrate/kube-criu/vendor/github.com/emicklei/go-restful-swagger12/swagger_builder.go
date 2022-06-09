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
package swagger

type SwaggerBuilder struct {
	SwaggerService
}

func NewSwaggerBuilder(config Config) *SwaggerBuilder {
	return &SwaggerBuilder{*newSwaggerService(config)}
}

func (sb SwaggerBuilder) ProduceListing() ResourceListing {
	return sb.SwaggerService.produceListing()
}

func (sb SwaggerBuilder) ProduceAllDeclarations() map[string]ApiDeclaration {
	return sb.SwaggerService.produceAllDeclarations()
}

func (sb SwaggerBuilder) ProduceDeclarations(route string) (*ApiDeclaration, bool) {
	return sb.SwaggerService.produceDeclarations(route)
}
