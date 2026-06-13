package common

// func GenerateCookies(c fiber.Ctx, at string, rt string, isRemember bool) {
// 	atExp := time.Minute * 5
// 	rtExp := enums.Day
// 	if isRemember {
// 		rtExp = enums.Month
// 	}

// 	isProd := configs.Env.ServerENV == "prod"

// 	cookie := fiber.Cookie{
// 		Name:     "at",
// 		Value:    at,
// 		HTTPOnly: true,
// 		Secure:   isProd,
// 		SameSite: "Strict",
// 		MaxAge:   int(atExp),
// 		Path:     "/",
// 	}
// 	c.Cookie(&cookie)

// 	cookie = fiber.Cookie{
// 		Name:     "rt",
// 		Value:    rt,
// 		HTTPOnly: true,
// 		Secure:   isProd,
// 		SameSite: "Strict",
// 		MaxAge:   int(rtExp),
// 		Path:     "/",
// 	}
// 	c.Cookie(&cookie)
// }
