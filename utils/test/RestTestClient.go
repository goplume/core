package test

import (
    "github.com/goplume/core/rest_client"
    "github.com/goplume/core/utils/stereotypes"
    "github.com/stretchr/testify/assert"
    "testing"
)

type RestTestClient struct {
	RestClient *rest_client.RestClient
}

func (this *RestTestClient) T_(
	t *testing.T,
	method string,
	expectedCode int,
	controller stereotypes.ApiRestInteface,
	entityId string,
	body interface{},
) {
	entityPath := controller.GetRoot() + "/" + controller.GetEntityPath()

	request := this.RestClient.HttpClient.R()
	if len(entityId) > 0 {
		entityPath = entityPath + "/{" + controller.GetEntityId() + "}"
		request.SetPathParams(map[string]string{
			controller.GetEntityId(): entityId,
		})
	}
	if body != nil {
		request.SetBody(body)
	}

	response, e := request.Execute(method, entityPath)

	if assert.Nil(t, e) {
		assert.Equal(t, expectedCode, response.StatusCode())
	}
}
