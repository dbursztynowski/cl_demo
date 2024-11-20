package tools

import (
	closedlooppoocv1 "closedloop/api/v1"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Cr_type(nameCR string, kindCR string) client.Object {
	var clcm client.Object
	switch strings.Title(nameCR) {
	case "Monitoring":
		if strings.Title(kindCR) == "Monitoringv2" {
			clcm = &closedlooppoocv1.Monitoringv2{}
		} else {
			clcm = &closedlooppoocv1.Monitoring{}
		}
	case "Decision":
		clcm = &closedlooppoocv1.Decision{}
	case "Execution":
		clcm = &closedlooppoocv1.Execution{}
	default:
		return nil
	}
	return clcm
}
