package terraform_utils

import (
	"log"
	"regexp"
)

func ConnectServices(importResources map[string][]Resource, resourceConnections map[string]map[string][]string) map[string][]Resource {
	for resource, connection := range resourceConnections {
		if _, exist := importResources[resource]; exist {
			for k, v := range connection {
				if cc, ok := importResources[k]; ok {
					for _, ccc := range cc {
						for i := range importResources[resource] {
							key := v[1]
							if v[1] == "self_link" || v[1] == "id" {
								key = ccc.GetIDKey()
							}
							keyValue := ccc.InstanceInfo.Type + "_" + ccc.ResourceName + "_" + key
							linkValue := "${data.terraform_remote_state." + k + "." + keyValue + "}"

							tfResource := importResources[resource][i]
							if ccc.InstanceState.Attributes[key] == tfResource.InstanceState.Attributes[v[0]] {
								importResources[resource][i].InstanceState.Attributes[v[0]] = linkValue
								importResources[resource][i].Item[v[0]] = linkValue
							} else {
								for keyAttributes, j := range tfResource.InstanceState.Attributes {
									match, err := regexp.MatchString(v[0]+".\\d+$", keyAttributes)
									if match && err == nil {
										if j == ccc.InstanceState.Attributes[key] {
											importResources[resource][i].InstanceState.Attributes[keyAttributes] = linkValue
											switch ar := tfResource.Item[v[0]].(type) {
											case []interface{}:
												for j, l := range ar {
													if l == ccc.InstanceState.Attributes[key] {
														importResources[resource][i].Item[v[0]].([]interface{})[j] = linkValue
													}
												}
											default:
												log.Println("type not supported", ar)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return importResources
}
