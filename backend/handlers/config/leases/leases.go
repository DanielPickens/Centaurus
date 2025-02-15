package leases

import (
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/coordination/v1"
	"net/http"

	"github.com/danielpickens/centaurus/backend/container"
	"github.com/danielpickens/centaurus/backend/handlers/base"
	"github.com/danielpickens/centaurus/backend/handlers/helpers"
	"github.com/labstack/echo/v4"
)

type LeasesHandler struct {
	BaseHandler base.BaseHandler
}

func NewLeaseRouteHandler(container container.Container, routeType base.RouteType) echo.HandlerFunc {
	return func(c echo.Context) error {
		handler := NewLeasesHandler(c, container)

		switch routeType {
		case base.GetList:
			return handler.BaseHandler.GetList(c)
		case base.GetDetails:
			return handler.BaseHandler.GetDetails(c)
		case base.GetEvents:
			return handler.BaseHandler.GetEvents(c)
		case base.GetYaml:
			return handler.BaseHandler.GetYaml(c)
		case base.Delete:
			return handler.BaseHandler.Delete(c)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Unknown route type")
		}
	}
}

func NewLeasesHandler(c echo.Context, container container.Container) *LeasesHandler {
	config := c.QueryParam("config")
	cluster := c.QueryParam("cluster")

	informer := container.SharedInformerFactory(config, cluster).Coordination().V1().Leases().Informer()
	informer.SetTransform(helpers.StripUnusedFields)

	handler := &LeasesHandler{
		BaseHandler: base.BaseHandler{
			Kind:             "Lease",
			Container:        container,
			Informer:         informer,
			RestClient:       container.ClientSet(config, cluster).CoordinationV1().RESTClient(),
			QueryConfig:      config,
			QueryCluster:     cluster,
			InformerCacheKey: fmt.Sprintf("%s-%s-leaseInformer", config, cluster),
			TransformFunc:    transformItems,
		},
	}
	cache := base.ResourceEventHandler[*v1.Lease](&handler.BaseHandler)
	handler.BaseHandler.StartInformer(c, cache)
	handler.BaseHandler.WaitForSync(c)
	return handler
}

func transformItems(items []interface{}, b *base.BaseHandler) ([]byte, error) {
	var list []v1.Lease

	for _, obj := range items {
		if item, ok := obj.(*v1.Lease); ok {
			list = append(list, *item)
		}
	}

	t := TransformLeaseList(list)

	return json.Marshal(t)
}


func TransformLeaseList(list []v1.Lease) []Lease {
	var leases []Lease
	for _, lease := range list {
		leases = append(leases, Lease{
			Name: lease.Name,
			Namespace: lease.Namespace,
			Owner: lease.OwnerReferences,
			AcquisitionTime: lease.Spec.AcquisitionTime,
			LeaseDurationSeconds: lease.Spec.LeaseDurationSeconds,
			RenewTime: lease.Spec.RenewTime,
			LeaseTransitions: lease.Status.LeaseTransitions,
		})
	}
	return leases
}

func 
