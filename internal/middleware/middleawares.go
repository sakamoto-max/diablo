package middleware

type ctxKey string

type Middlewares struct {
	Auth *Auth
}

func NewMiddlewares() *Middlewares {
	return &Middlewares{
		Auth: &Auth{},
	}
}
