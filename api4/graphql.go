// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/graph-gophers/dataloader/v6"
	graphql "github.com/graph-gophers/graphql-go"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
)

type graphQLInput struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

//go:embed schema.graphqls
var schemaRaw string

func (api *API) InitGraphQL() error {
	// Guard with a feature flag.
	if !api.srv.Config().FeatureFlags.GraphQL {
		return nil
	}

	var err error
	opts := []graphql.SchemaOpt{
		graphql.UseFieldResolvers(),
		graphql.Logger(mlog.NewGraphQLLogger(api.srv.Log)),
		graphql.MaxParallelism(5),
	}

	if isProd() {
		opts = append(opts,
			// MaxDepth cannot be moved as a general param
			// because otherwise introspection also doesn't work
			// with just a depth of 4.
			graphql.MaxDepth(4),
			graphql.DisableIntrospection(),
		)
	}

	api.schema, err = graphql.ParseSchema(schemaRaw, &resolver{}, opts...)
	if err != nil {
		return err
	}

	api.BaseRoutes.APIRoot5.Handle("/graphql", api.APIHandlerTrustRequester(graphiQL)).Methods("GET")
	api.BaseRoutes.APIRoot5.Handle("/graphql", api.APISessionRequired(api.graphQL)).Methods("POST")
	return nil
}

// Unique type to hold our context.
type ctxKey int

const (
	webCtx            ctxKey = 0
	rolesLoaderCtx    ctxKey = 1
	channelsLoaderCtx ctxKey = 2
)

const loaderBatchCapacity = 200

func (api *API) graphQL(c *Context, w http.ResponseWriter, r *http.Request) {
	var response *graphql.Response
	defer func() {
		if response != nil {
			if err := json.NewEncoder(w).Encode(response); err != nil {
				mlog.Warn("Error while writing response", mlog.Err(err))
			}
		}
	}()

	// Limit bodies to 100KiB.
	// We need to enforce a lower limit than the file upload size,
	// to prevent the library doing unnecessary parsing.
	r.Body = http.MaxBytesReader(w, r.Body, 102400)

	var params graphQLInput
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		err2 := gqlerrors.Errorf("invalid request body: %v", err)
		response = &graphql.Response{Errors: []*gqlerrors.QueryError{err2}}
		return
	}

	if isProd() && params.OperationName == "" {
		err2 := gqlerrors.Errorf("operation name not passed")
		response = &graphql.Response{Errors: []*gqlerrors.QueryError{err2}}
		return
	}

	// Populate the context with required info.
	reqCtx := r.Context()
	reqCtx = context.WithValue(reqCtx, webCtx, c)

	rolesLoader := dataloader.NewBatchedLoader(graphQLRolesLoader, dataloader.WithBatchCapacity(loaderBatchCapacity))
	reqCtx = context.WithValue(reqCtx, rolesLoaderCtx, rolesLoader)

	channelsLoader := dataloader.NewBatchedLoader(graphQLChannelsLoader, dataloader.WithBatchCapacity(loaderBatchCapacity))
	reqCtx = context.WithValue(reqCtx, channelsLoaderCtx, channelsLoader)

	response = api.schema.Exec(reqCtx,
		params.Query,
		params.OperationName,
		params.Variables)

	if len(response.Errors) > 0 {
		logFunc := mlog.Error
		for _, gqlErr := range response.Errors {
			if gqlErr.Err != nil {
				if appErr, ok := gqlErr.Err.(*model.AppError); ok && appErr.StatusCode < http.StatusInternalServerError {
					logFunc = mlog.Debug
					break
				}
			}
		}
		logFunc("Error executing request", mlog.String("operation", params.OperationName),
			mlog.Array("errors", response.Errors))
	}
}

func graphiQL(c *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(graphiqlPage)
}

var graphiqlPage = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<title>GraphiQL editor | Mattermost</title>
		<link href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css" rel="stylesheet" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/es6-promise/4.1.1/es6-promise.auto.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/api/v5/graphql", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
					headers: {
						'X-Requested-With': 'XMLHttpRequest'
					}
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)

// isProd is a helper function to apply prod-specific graphQL validations.
func isProd() bool {
	return model.BuildNumber != "dev"
}
