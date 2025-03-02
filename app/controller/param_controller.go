package controller

import (
	"bpkp-svc-portal/app/client"
	"bpkp-svc-portal/app/model"
	"bpkp-svc-portal/app/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type InterfaceParamController interface {
	GetParameterByKey(ctx context.Context, key string) (*model.Param, error)
	GetAllParam(ctx context.Context) ([]*model.Param, error)
	InsertNewParam(ctx context.Context, param *model.Param) error
	UpdateParam(ctx context.Context, param *model.Param) error
	DeleteParam(ctx context.Context, key string) error
}

type ParamController struct {
	redis  *redis.Client
	client client.InterfaceParamClient
}

func NewParamController(redis *redis.Client, client client.InterfaceParamClient) *ParamController {
	return &ParamController{
		redis:  redis,
		client: client,
	}
}

func (c *ParamController) GetParameterByKey(ctx context.Context, key string) (*model.Param, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetParameterByKey")
	defer span.Finish()

	utils.LogEvent(span, "Request", key)

	res, err := c.client.GetParameterByKey(ctx, key)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", res)

	return res, nil
}

func (c *ParamController) GetAllParam(ctx context.Context) ([]*model.Param, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetAllParam")
	defer span.Finish()

	utils.LogEvent(span, "Request", "All")

	res, err := c.client.GetAllParam(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", res)

	return res, nil
}

func (c *ParamController) InsertNewParam(ctx context.Context, param *model.Param) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: InsertNewParam")
	defer span.Finish()

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	param.UpdatedAt = utils.LocalTime()
	param.UpdatedBy = session.Username

	utils.LogEvent(span, "Request", param)

	err = c.client.InsertNewParam(ctx, param)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Insert New Param")

	return nil
}

func (c *ParamController) UpdateParam(ctx context.Context, param *model.Param) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: UpdateParam")
	defer span.Finish()

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	param.UpdatedAt = utils.LocalTime()
	param.UpdatedBy = session.Username

	utils.LogEvent(span, "Request", param)

	err = c.client.UpdateParam(ctx, param)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	newValueJSON, err := json.Marshal(param)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	if err := c.redis.Set(ctx, param.Key, newValueJSON, 6*time.Hour).Err(); err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Update Param")

	return nil
}

func (c *ParamController) DeleteParam(ctx context.Context, key string) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: DeleteParam")
	defer span.Finish()

	utils.LogEvent(span, "Request", key)

	err := c.client.DeleteParam(ctx, key)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	err = c.redis.Del(ctx, key).Err()
	if err != nil {
		fmt.Println("Hereee")
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Param")

	return nil
}
