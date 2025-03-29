package abstractions

type Requestable[T, S any] interface {
	Request(S) (*T, error)
}

func MakeRequest[T Requestable[T, S], S any](request S) (*T, error) {
	req := new(T)

	result, err := (*req).Request(request)

	return result, err
}

type Responseable[S any] interface {
	Response() (*S, error)
}

func MakeResponse[S any, T Responseable[S]](response T) (*S, error) {
	return response.Response()
}
