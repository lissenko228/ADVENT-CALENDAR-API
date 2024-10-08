package handler

import (
	"advent-calendar/internal/config"
	"advent-calendar/internal/repository"
	"advent-calendar/pkg/utils"
	"advent-calendar/pkg/validators"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Exp          int64  `json:"exp"`
}

// @Tags Users
// @Param request formData repository.LoginDTO true "-"
// @Failure 401 {object} validators.GlobalHandlerResp
// @Success 200 {object} Tokens
// @Router /api/users/login [post]
func Login(c *fiber.Ctx) error {
	data := new(repository.LoginDTO)

	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(400, err.Error())
	}

	if errs := config.Validator.Validate(data); len(errs) > 0 && errs[0].Error {
		return validators.ValidateError(errs)
	}

	user, _ := repository.UserService.Get(repository.User{Email: data.Email})

	if user.Role == "admin" {
		if !utils.CheckPasswordHash(data.Password, user.Password) {
			return fiber.NewError(401, "Неправильный пароль")
		}

		jwt, exp, err := utils.NewJWT(user.ID, user.Role)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		tokens := Tokens{AccessToken: jwt, RefreshToken: user.RefreshToken, Exp: exp}

		return c.JSON(tokens)
	}
	code := utils.GenerateCode()

	if len(user.Email) > 0 {
		if err := user.Update(repository.User{Code: code}); err != nil {
			return fiber.NewError(500, err.Error())
		}
	} else {
		if _, err := repository.UserService.Create(repository.User{Email: data.Email, Code: code}); err != nil {
			return fiber.NewError(500, err.Error())
		}
	}

	type ToSend struct {
		Code string
		Year int
	}

	body, err := utils.LoadTemplate("registration.email", ToSend{Code: code, Year: time.Now().Year()})
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	if err := utils.SendMail(data.Email, "Завершение регистрации - Кибербезопасный Новый год", body.String()); err != nil {
		return fiber.NewError(500, "Ошибка при отправке письма")
	}

	return c.JSON(validators.GlobalHandlerResp{Success: true, Message: "Проверьте Ваш email, для продолжения регистрации"})

}

// @Tags Users
// @Param Authorization header string true "Authorization"
// @Failure 401 {object} validators.GlobalHandlerResp
// @Success 200 {object} repository.User
// @Router /api/users/check [get]
func Check(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(utils.Claims)

	user, err := repository.UserService.Get(repository.User{ID: userClaims.ID})
	if err != nil {
		return fiber.NewError(401, "Ошибка при проверке авторизации")
	}

	return c.JSON(user)
}

// @Tags Users
// @Param RefreshToken header string true "RefreshToken"
// @Failure 500 {object} validators.GlobalHandlerResp
// @Failure 401 {object} validators.GlobalHandlerResp
// @Success 200 {object} Tokens
// @Router /api/users/refresh [patch]
func Refresh(c *fiber.Ctx) error {
	refreshToken := c.Get("RefreshToken")

	user, err := repository.UserService.Get(repository.User{RefreshToken: refreshToken})
	if err != nil {
		return fiber.NewError(401, "Неправильный токен")
	}

	jwt, exp, err := utils.NewJWT(user.ID, user.Role)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	newToken, err := utils.NewRefreshToken()
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	err = user.Update(repository.User{RefreshToken: newToken})
	if err != nil {
		return fiber.NewError(500, "Ошибка при обновлении токена")
	}

	tokens := Tokens{AccessToken: jwt, RefreshToken: user.RefreshToken, Exp: exp}

	return c.JSON(tokens)
}

// @Tags Users
// @Param request formData repository.ConfirmUser true "-"
// @Failure 500 {object} validators.GlobalHandlerResp
// @Failure 401 {object} validators.GlobalHandlerResp
// @Success 200 {object} Tokens
// @Router /api/users/confirm [patch]
func ConfirmRegister(c *fiber.Ctx) error {
	data := new(repository.ConfirmUser)

	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(400, err.Error())
	}

	if errs := config.Validator.Validate(data); len(errs) > 0 && errs[0].Error {
		return validators.ValidateError(errs)
	}

	user, err := repository.UserService.Get(repository.User{Code: data.Code, Email: data.Email})
	if err != nil {
		return fiber.NewError(401, "Неправильный код")
	}

	refreshToken, err := utils.NewRefreshToken()
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	if err := user.Update(repository.User{Code: " ", RefreshToken: refreshToken}); err != nil {
		return fiber.NewError(500, err.Error())
	}

	jwt, exp, err := utils.NewJWT(user.ID, user.Role)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(Tokens{AccessToken: jwt, RefreshToken: refreshToken, Exp: exp})
}

// @Tags Users
// @Param request formData repository.SubscribeDTO true "-"
// @Failure 500 {object} validators.GlobalHandlerResp
// @Failure 400 {object} validators.GlobalHandlerResp
// @Success 200 {object} validators.GlobalHandlerResp
// @Router /api/users/subscribe [post]
func Subscribe(c *fiber.Ctx) error {
	data := new(repository.SubscribeDTO)

	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(400, err.Error())
	}

	if errs := config.Validator.Validate(data); len(errs) > 0 && errs[0].Error {
		return validators.ValidateError(errs)
	}

	if err := repository.UserService.Subscribe(data); err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(validators.GlobalHandlerResp{Success: true, Message: "Подписка оформлена"})
}

// @Tags Users
// @Param request query repository.UnSubscribeDTO true "-"
// @Failure 500 {object} validators.GlobalHandlerResp
// @Failure 400 {object} validators.GlobalHandlerResp
// @Success 200 {object} validators.GlobalHandlerResp
// @Router /api/users/subscribe [delete]
func UnSubscribe(c *fiber.Ctx) error {
	data := new(repository.UnSubscribeDTO)

	err := c.QueryParser(data)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	if errs := config.Validator.Validate(data); len(errs) > 0 && errs[0].Error {
		return validators.ValidateError(errs)
	}

	if err := repository.UserService.UnSubscribe(data.Email); err != nil {
		return fiber.NewError(500, err.Error())
	}

	if c.Method() == "GET" {
		return c.Redirect(config.Config.CLIENT_URI)
	}

	return c.JSON(validators.GlobalHandlerResp{Success: true, Message: "Подписка отменена"})
}
