package pods

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/danielpickens/centaurus/backend/handlers/workloads/replicaset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"sync"

	"github.com/danielpickens/centaurus/backend/handlers/base"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/gorilla/websocket"
	"github.com/danielpickens/centaurus/backend/event"
	"github.com/danielpickens/centaurus/backend/handlers/helpers"

	"net/http"
	"strings"
	"time"

	"github.com/danielpickens/centaurus/backend/container"
	"github.com/labstack/echo/v4"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

type PodsHandler struct {
	BaseHandler       base.BaseHandler
	clientSet         *kubernetes.Clientset
	restConfig        *rest.Config
	replicasetHandler *replicaset.ReplicaSetHandler
}

func NewPodsRouteHandler(container container.Container, routeType base.RouteType) echo.HandlerFunc {
	return func(c echo.Context) error {
		handler := NewPodsHandler(c, container)

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
		case base.GetLogsWS:
			return handler.GetLogsWS(c)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Unknown route type")
		}
	}
}

func NewPodsHandler(c echo.Context, container container.Container) *PodsHandler {
	config := c.QueryParam("config")
	cluster := c.QueryParam("cluster")

	informer := container.SharedInformerFactory(config, cluster).Core().V1().Pods().Informer()
	informer.SetTransform(helpers.StripUnusedFields)
	clientSet := container.ClientSet(config, cluster)

	handler := &PodsHandler{
		BaseHandler: base.BaseHandler{
			Kind:             "Pod",
			Container:        container,
			RestClient:       clientSet.CoreV1().RESTClient(),
			Informer:         informer,
			QueryConfig:      config,
			QueryCluster:     cluster,
			InformerCacheKey: fmt.Sprintf("%s-%s-podInformer", config, cluster),
			TransformFunc:    transformItems,
		},
		restConfig:        container.RestConfig(config, cluster),
		clientSet:         clientSet,
		replicasetHandler: replicaset.NewReplicaSetHandler(c, container),
	}

	additionalEvents := []map[string]func(){
		{
			"pods-deployments": func() {
				handler.DeploymentsPods(c)
			},
		},
	}

	cache := base.ResourceEventHandler[*v1.Pod](&handler.BaseHandler, additionalEvents...)
	handler.BaseHandler.StartInformer(c, cache)
	handler.BaseHandler.WaitForSync(c)

	return handler
}

func transformItems(items []interface{}, b *base.BaseHandler) ([]byte, error) {
	var list []v1.Pod
	for _, obj := range items {
		if item, ok := obj.(*v1.Pod); ok {
			list = append(list, *item)
		}
	}
	podMetricsList := GetPodsMetricsList(b)
	t := TransformPodList(list, podMetricsList)

	return json.Marshal(t)
}

func GetPodsMetricsList(b *base.BaseHandler) *v1beta1.PodMetricsList {
	cacheKey := fmt.Sprintf(helpers.IsMetricServerAvailableCacheKeyFormat, b.QueryConfig, b.QueryCluster)
	value, exists := b.Container.Cache().Get(cacheKey)
	if value == false || exists == false {
		return nil
	}
	podMetrics, err := b.Container.
		MetricClient(b.QueryConfig, b.QueryCluster).
		MetricsV1beta1().
		PodMetricses("").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Info("failed to get pod metrics", "err", err)
	}
	return podMetrics
}

func (h *PodsHandler) GetLogsWS(c echo.Context) error {
	ws, err := h.BaseHandler.Container.SocketUpgrader().Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	name := c.Param("name")
	namespace := c.QueryParam("namespace")
	container := c.QueryParam("container")
	isAllContainers := strings.EqualFold(c.QueryParam("all-containers"), "true")

	var containerNames []string

	if isAllContainers {
		podObj, _, err := h.BaseHandler.Informer.GetStore().GetByKey(fmt.Sprintf("%s/%s", c.QueryParam("namespace"), c.Param("name")))
		if err != nil {
			return err
		}
		pod := podObj.(*v1.Pod)
		for _, logContainer := range pod.Spec.Containers {
			containerNames = append(containerNames, logContainer.Name)
		}
	} else {
		containerNames = []string{container}
	}

	logsChannel := make(chan LogMessage)
	defer close(logsChannel)

	for _, containerName := range containerNames {
		go h.fetchLogs(c.Request().Context(), namespace, name, containerName, logsChannel)
	}
	event := event.NewEventCounter(250 * time.Millisecond)
	var logMessages []LogMessage
	var mu sync.RWMutex

	go event.Run()

	for logMsg := range logsChannel {
		select {
		case <-c.Request().Context().Done():
			c.Logger().Info("request context cancelled, closing logs channel")
			return nil
		default:
			mu.Lock()
			event.AddEvent(fmt.Sprintf("pod-logs-%s", container), func() {
				mu.Lock()
				defer mu.Unlock()
				if len(logMessages) > 0 {
					j, err := json.Marshal(logMessages)
					if err != nil {
						c.Logger().Errorf("failed to marshal log message: %v", err)
					}

					err = ws.WriteMessage(websocket.TextMessage, j)
					if err != nil {
						c.Logger().Error(err)
					}
				}
				logMessages = []LogMessage{}
			})
			logMessages = append(logMessages, logMsg)
			mu.Unlock()
		}
	}

	return nil
}
