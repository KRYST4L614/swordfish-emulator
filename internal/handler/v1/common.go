package v1

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

// contextResourceKey - type for inserting service defined values in context.
// Should be used to avoid collision with mux (or other) package context values.
type contextResourceKey string

const (
	// idContext - is key for Id value in context
	idContext contextResourceKey = "Id"
)

const (
	// idPathRegex - regex that accepts all characters except URI unsafe (https://datatracker.ietf.org/doc/html/rfc1738)
	idPathRegex string = "/{id:[^/<>\"#%{}|\\\\^~[\\]]+}"
	// idPathTag just tag in 'idPathRegex'
	idPathTag string = "id"
)

// injectId - gets resource Id from request route
// and inserts value in request context, returning new
// request value
//
// 'tag' parameter identifies route tag for 'Id' mapping
func injectId(request *http.Request, tag string) *http.Request {
	if mux.Vars(request)[tag] != "" {
		request = request.WithContext(
			context.WithValue(request.Context(), idContext, mux.Vars(request)[tag]),
		)
	}

	return request
}

// resourceGetter - returns closure for common resource GET handler
//
// GET can be common for all resources, as no validation required
func resourceGetter(service service.ResourceService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		resource, err := service.Get(request.Context(), request.RequestURI)
		if err != nil {
			util.WriteJSONError(writer, err)
			return
		}

		util.WriteJSON(writer, resource)
	}
}

// resourceDeleter - returns closure for common resource DELETE handler
//
// DELETE can be common for all resources, as no validation required
func resourceDeleter(service service.ResourceService) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		resourceId := request.RequestURI
		resource, err := service.DeleteResourceFromCollection(request.Context(), util.GetParent(resourceId), resourceId)
		if err != nil {
			util.WriteJSONError(writer, err)
			return
		}

		util.WriteJSON(writer, resource)
	}
}

// resourceCreatorFromNotCollectionEndpoint - returns closure that routes request from
// /<Collection>/<id> or /<Collection>/Members to /<Collection> endpoint
//
// Not a good way, should be removed when Ansible modules will handle
// resources creation in right way (without /<Collection>/<id> request)
func resourceCreatorFromNotCollectionEndpoint(
	createResourceFromCollectionEndpoint func(writer http.ResponseWriter, request *http.Request),
) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		request = injectId(request, idPathTag)
		request.RequestURI = util.GetParent(request.RequestURI)
		createResourceFromCollectionEndpoint(writer, request)
	}
}
