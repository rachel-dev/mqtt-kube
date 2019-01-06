/*
Package handlers : handle MQTT message and deploy object to kubernetes.
	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package handlers

import (
	"fmt"

	"go.uber.org/zap"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type configmapHandler struct {
	kubeClient kubernetes.Interface
	logger     *zap.SugaredLogger
}

func newConfigmapHandler(clientset *kubernetes.Clientset, logger *zap.SugaredLogger) *configmapHandler {
	return &configmapHandler{
		kubeClient: clientset,
		logger:     logger,
	}
}

func (h *configmapHandler) Apply(rawData runtime.Object) string {
	configmap := rawData.(*apiv1.ConfigMap)
	configmapsClient := h.kubeClient.CoreV1().ConfigMaps(apiv1.NamespaceDefault)
	name := configmap.ObjectMeta.Name
	current, getErr := configmapsClient.Get(name, metav1.GetOptions{})

	if current != nil && getErr == nil {
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			current.ObjectMeta.Labels = configmap.ObjectMeta.Labels
			current.ObjectMeta.Annotations = configmap.ObjectMeta.Annotations
			current.Data = configmap.Data
			_, err := configmapsClient.Update(current)
			return err
		})
		if err != nil {
			msg := fmt.Sprintf("update configmap err -- %s", name)
			h.logger.Errorf("%s: %s", msg, err.Error())
			return msg
		}
		msg := fmt.Sprintf("update configmap -- %s", name)
		h.logger.Infof(msg)
		return msg
	} else if errors.IsNotFound(getErr) {
		result, err := configmapsClient.Create(configmap)
		if err != nil {
			msg := fmt.Sprintf("create configmap err -- %s", name)
			h.logger.Errorf("%s: %s", msg, err.Error())
			return msg
		}
		msg := fmt.Sprintf("create configmap -- %s", result.GetObjectMeta().GetName())
		h.logger.Infof(msg)
		return msg
	} else {
		msg := fmt.Sprintf("get configmap err -- %s", name)
		h.logger.Errorf("%s: %s", msg, getErr.Error())
		return msg
	}
}

func (h *configmapHandler) Delete(rawData runtime.Object) string {
	configmap := rawData.(*apiv1.ConfigMap)
	configmapsClient := h.kubeClient.CoreV1().ConfigMaps(apiv1.NamespaceDefault)
	name := configmap.ObjectMeta.Name
	current, getErr := configmapsClient.Get(name, metav1.GetOptions{})

	if current != nil && getErr == nil {
		deletePolicy := metav1.DeletePropagationForeground
		if err := configmapsClient.Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			msg := fmt.Sprintf("delete configmap err -- %s", name)
			h.logger.Errorf("%s: %s", msg, err.Error())
			return msg
		}
		msg := fmt.Sprintf("delete configmap -- %s", name)
		h.logger.Infof(msg)
		return msg
	} else if errors.IsNotFound(getErr) {
		msg := fmt.Sprintf("configmap does not exist -- %s", name)
		h.logger.Infof(msg)
		return msg
	} else {
		msg := fmt.Sprintf("get configmap err -- %s", name)
		h.logger.Errorf("%s: %s", msg, getErr.Error())
		return msg
	}
}
