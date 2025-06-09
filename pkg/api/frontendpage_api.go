package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	frontendv1alpha1 "github.com/yourusername/k8s-controller-tutorial/pkg/apis/frontend/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FrontendPageAPI provides handlers for FrontendPage resources.
type FrontendPageAPI struct {
	K8sClient client.Client
	Namespace string // default namespace for simplicity
}

// ListFrontendPages godoc
// @Summary List all FrontendPages
// @Description Get all FrontendPage resources
// @Tags frontendpages
// @Produce json
// @Success 200 {array} frontendv1alpha1.FrontendPage
// @Router /api/frontendpages [get]
func (api *FrontendPageAPI) ListFrontendPages(ctx *fasthttp.RequestCtx) {
	list := &frontendv1alpha1.FrontendPageList{}
	err := api.K8sClient.List(context.Background(), list, client.InNamespace(api.Namespace))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(fmt.Sprintf(`{"error":"%v"}`, err))
		return
	}
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(list.Items)
}

// GetFrontendPage godoc
// @Summary Get a FrontendPage
// @Description Get a FrontendPage by name
// @Tags frontendpages
// @Produce json
// @Param name path string true "FrontendPage name"
// @Success 200 {object} frontendv1alpha1.FrontendPage
// @Failure 404 {object} map[string]string
// @Router /api/frontendpages/{name} [get]
func (api *FrontendPageAPI) GetFrontendPage(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	obj := &frontendv1alpha1.FrontendPage{}
	err := api.K8sClient.Get(context.Background(), client.ObjectKey{Namespace: api.Namespace, Name: name}, obj)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBodyString(fmt.Sprintf(`{"error":"%v"}`, err))
		return
	}
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(obj)
}

// CreateFrontendPage godoc
// @Summary Create a FrontendPage
// @Description Create a new FrontendPage
// @Tags frontendpages
// @Accept json
// @Produce json
// @Param body body frontendv1alpha1.FrontendPage true "FrontendPage object"
// @Success 201 {object} frontendv1alpha1.FrontendPage
// @Failure 400 {object} map[string]string
// @Router /api/frontendpages [post]
func (api *FrontendPageAPI) CreateFrontendPage(ctx *fasthttp.RequestCtx) {
	obj := &frontendv1alpha1.FrontendPage{}
	if err := json.Unmarshal(ctx.PostBody(), obj); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(fmt.Sprintf(`{"error":"%v"}`, err))
		return
	}
	if err := api.K8sClient.Create(context.Background(), obj); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(fmt.Sprintf(`{"error":"%v"}`, err))
		return
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(obj)
}

// UpdateFrontendPage godoc
// @Summary Update a FrontendPage
// @Description Update an existing FrontendPage
// @Tags frontendpages
// @Accept json
// @Produce json
// @Param name path string true "FrontendPage name"
// @Param body body frontendv1alpha1.FrontendPage true "FrontendPage object"
// @Success 200 {object} frontendv1alpha1.FrontendPage
// @Failure 400 {object} map[string]string
// @Router /api/frontendpages/{name} [put]
func (api *FrontendPageAPI) UpdateFrontendPage(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	obj := &frontendv1alpha1.FrontendPage{}
	if err := json.Unmarshal(ctx.PostBody(), obj); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(fmt.Sprintf(`{"error":"%v"}`, err))
		return
	}
	obj.Name = name
	obj.Namespace = api.Namespace
	if err := api.K8sClient.Update(context.Background(), obj); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(fmt.Sprintf(`{"error":"%v"}`, err))
		return
	}
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(obj)
}

// DeleteFrontendPage godoc
// @Summary Delete a FrontendPage
// @Description Delete a FrontendPage by name
// @Tags frontendpages
// @Param name path string true "FrontendPage name"
// @Success 204 {object} nil
// @Failure 404 {object} map[string]string
// @Router /api/frontendpages/{name} [delete]
func (api *FrontendPageAPI) DeleteFrontendPage(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	obj := &frontendv1alpha1.FrontendPage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: api.Namespace,
		},
	}
	if err := api.K8sClient.Delete(context.Background(), obj); err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBodyString(fmt.Sprintf(`{"error":"%v"}`, err))
		return
	}
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

// Note: Wire these handlers into your FastHTTP server in cmd/server.go using the appropriate routing logic.
//       You can use a router like fasthttprouter or manually parse the path and method.
