package service

import (
	"bpkp-svc-portal/app/controller"
	"bpkp-svc-portal/app/model"
	"bpkp-svc-portal/app/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type InterfaceUserService interface {
	CreateNewUser(e echo.Context) error
	GetUserDetail(e echo.Context) error
	UpdateUser(e echo.Context) error
	DeleteUser(e echo.Context) error
	Login(e echo.Context) error
	GetAllUser(e echo.Context) error
	GetInstitutionList(e echo.Context) error
	EmbedMetabase(e echo.Context) error
}

type UserService struct {
	uc controller.InterfaceUserController
}

func NewUserService(uc controller.InterfaceUserController) InterfaceUserService {
	return &UserService{
		uc: uc,
	}
}

func (s *UserService) CreateNewUser(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "CreateNewUser")
	defer span.Finish()

	var request *model.User

	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	err := s.uc.CreateNewUser(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success Create New User")

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create New User",
		Data:    nil,
	})
}

func (s *UserService) GetUserDetail(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetUserDetail")
	defer span.Finish()

	username := e.Param("id")

	utils.LogEvent(span, "Request", username)

	user, err := s.uc.GetUserDetail(ctx, username)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get User Detail",
		Data:    user,
	})
}

func (s *UserService) UpdateUser(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "UpdateUser")
	defer span.Finish()

	var request *model.User

	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	err := s.uc.UpdateUser(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Update User",
		Data:    nil,
	})
}

func (s *UserService) DeleteUser(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "DeleteUser")
	defer span.Finish()

	username := e.Param("id")

	utils.LogEvent(span, "Request", username)

	err := s.uc.DeleteUser(ctx, username)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Delete User",
		Data:    nil,
	})
}

func (s *UserService) Login(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "Login")
	defer span.Finish()

	var request *model.RequestLogin

	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	response, err := s.uc.Login(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Login",
		Data:    response,
	})
}

func (s *UserService) GetAllUser(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetAlluser")
	defer span.Finish()

	users, err := s.uc.GetAllUser(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", users)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get All User",
		Data:    users,
	})
}

func (s *UserService) GetInstitutionList(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetInstitutionList")
	defer span.Finish()

	institutionList, err := s.uc.GetInstitutionList(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", institutionList)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Institution List",
		Data:    institutionList,
	})
}

func (s *UserService) EmbedMetabase(e echo.Context) error {
	_, span := utils.StartSpan(e, "EmbedMetabase")
	defer span.Finish()

	claims := jwt.MapClaims{
		"resource": map[string]int{"dashboard": 3},
		"params":   map[string]interface{}{},
		"exp":      time.Now().Add(10 * time.Minute).Unix(), // 10 minutes expiration
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte("2e96319d77bc2bcff04d923e4ebc9ef0b3636dac8fad1584ade05516f92e8796"))
	if err != nil {
		fmt.Println("Error signing token:", err)
		return utils.LogError(e, err, nil)
	}

	// Generate the iframe URL
	iframeURL := fmt.Sprintf("%s/embed/dashboard/%s#bordered=true&titled=true", "https://metabase.eventarry.com", signedToken)

	utils.LogEvent(span, "Response", iframeURL)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Embed Metabase",
		Data:    iframeURL,
	})
}
