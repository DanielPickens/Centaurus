package horizontalpodautoscalers

import (
	"encoding/json"
	"fmt"
	autoScalingV2 "k8s.io/api/autoscaling/v2"
	"net/http"

	"github.com/danielpickens/centaurus/backend/container"
	"github.com/danielpickens/centaurus/backend/handlers/base"
	"github.com/danielpickens/centaurus/backend/handlers/helpers"
	"github.com/labstack/echo/v4"
)

type HorizontalPodAutoScalerHandler struct {
	BaseHandler base.BaseHandler
}

func NewHorizontalPodAutoscalersRouteHandler(container container.Container, routeType base.RouteType) echo.HandlerFunc {
	return func(c echo.Context) error {
		handler := NewHorizontalPodAutoScalerHandler(c, container)

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

func NewHorizontalPodAutoScalerHandler(c echo.Context, container container.Container) *HorizontalPodAutoScalerHandler {
	config := c.QueryParam("config")
	cluster := c.QueryParam("cluster")

	informer := container.SharedInformerFactory(config, cluster).Autoscaling().V2().HorizontalPodAutoscalers().Informer()
	informer.SetTransform(helpers.StripUnusedFields)

	handler := &HorizontalPodAutoScalerHandler{
		BaseHandler: base.BaseHandler{
			Kind:             "HorizontalPodAutoscaler",
			Container:        container,
			Informer:         informer,
			RestClient:       container.ClientSet(config, cluster).AutoscalingV2().RESTClient(),
			QueryConfig:      config,
			QueryCluster:     cluster,
			InformerCacheKey: fmt.Sprintf("%s-%s-horizontalPodAutoscalerInformer", config, cluster),
			TransformFunc:    transformItems,
		},
	}
	cache := base.ResourceEventHandler[*autoScalingV2.HorizontalPodAutoscaler](&handler.BaseHandler)
	handler.BaseHandler.StartInformer(c, cache)
	handler.BaseHandler.WaitForSync(c)
	return handler
}

func transformItems(items []interface{}, _ *base.BaseHandler) ([]byte, error) {
	var list []autoScalingV2.HorizontalPodAutoscaler

	for _, obj := range items {
		if item, ok := obj.(*autoScalingV2.HorizontalPodAutoscaler); ok {
			list = append(list, *item)
		}
	}

	t := TransformHorizontalPodAutoscaler(list)

	return json.Marshal(t)
}

func TransformHorizontalPodAutoscaler(list []autoScalingV2.HorizontalPodAutoscaler) []HorizontalPodAutoscaler {
	var hpaList []HorizontalPodAutoscaler

	for _, hpa := range list {
		hpaList = append(hpaList, HorizontalPodAutoscaler{
			Name: hpa.Name,
			UID:  string(hpa.UID),
		})
	}

	return hpaList
}

type HorizontalPodAutoscaler struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

func (h *HorizontalPodAutoScalerHandler) GetList(c echo.Context) error {
	return h.BaseHandler.GetList(c)
}

func (h *HorizontalPodAutoScalerHandler) GetDetails(c echo.Context) error {
	return h.BaseHandler.GetDetails(c)
}

func (h *HorizontalPodAutoScalerHandler) GetEvents(c echo.Context) error {
	return h.BaseHandler.GetEvents(c)
}


