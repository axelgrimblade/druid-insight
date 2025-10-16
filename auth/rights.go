package auth

import (
	"druid-insight/config"
)

func CheckRights(payload map[string]interface{}, druidCfg *config.DruidConfig, datasource string, isAdmin bool) []string {
	problems := []string{}
	ds, ok := druidCfg.Datasources[datasource]
	if !ok {
		return []string{"datasource_not_found"}
	}
	if dims, ok := payload["dimensions"].([]interface{}); ok {
		for _, dimRaw := range dims {
			dim, _ := dimRaw.(string)
			if dim == "time" {
				// La dimension "time" est TOUJOURS autorisée
				continue
			}
			f, ok := ds.Dimensions[dim]
			if !ok {
				problems = append(problems, "dimension:"+dim+":unknown")
			} else if f.Reserved && !isAdmin {
				problems = append(problems, "dimension:"+dim+":forbidden")
			}
		}
	}
	if mets, ok := payload["metrics"].([]interface{}); ok {
		for _, mRaw := range mets {
			metric, _ := mRaw.(string)
			f, ok := ds.Metrics[metric]
			if !ok {
				problems = append(problems, "metric:"+metric+":unknown")
			} else if f.Reserved && !isAdmin {
				problems = append(problems, "metric:"+metric+":forbidden")
			}
		}
	}
	if filters, ok := payload["filters"].([]interface{}); ok {
		for _, filterRaw := range filters {
			filter, _ := filterRaw.(map[string]interface{})
			if dimensionRaw, exists := filter["dimension"]; exists {
				dimension, _ := dimensionRaw.(string)
				if dimension == "time" {
					// La dimension "time" n'est pas autorisée dans les filtres
					problems = append(problems, "filter_dimension:time:forbidden")
				}
				f, ok := ds.Dimensions[dimension]
				if !ok {
					problems = append(problems, "filter_dimension:"+dimension+":unknown")
				} else if f.Reserved && !isAdmin {
					problems = append(problems, "filter_dimension:"+dimension+":forbidden")
				}
			} else {
				problems = append(problems, "filter_missing_dimension")
			}
			if valueRaw, exists := filter["values"]; exists {
				_, ok := valueRaw.(string)
				if !ok {
					_, ok := valueRaw.([]interface{})
					if !ok {
						problems = append(problems, "filter_value_invalid")
					}
				}
			} else {
				problems = append(problems, "filter_missing_value")
			}
		}
	}
	return problems
}
