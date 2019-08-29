package rest_api

type HealthApiRestController struct {
	BaseRestController
	HealthApiService HealthApiService
}

func (this *HealthApiRestController) InitController() {}
