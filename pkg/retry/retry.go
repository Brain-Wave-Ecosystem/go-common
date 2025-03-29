package retry

import (
	"context"
	"math/rand/v2"
	"time"
)

var StandardConfig = &Config{
	Attempts:  5,
	Delay:     1 * time.Second,
	MaxDelay:  10 * time.Second,
	MaxJitter: 1.1,
	Context:   context.Background(),
}

type Config struct {
	Attempts    uint
	Delay       time.Duration
	MaxDelay    time.Duration
	MaxJitter   float64
	OnRetry     func(uint, error) error
	RetryIf     func(error) bool
	ExposeDelay func(n uint, err error) time.Duration
	Context     context.Context
}

type Retry struct {
	*Config
}

func NewRetry(config *Config) *Retry {
	var r Retry
	if config == nil {
		config = StandardConfig
	}

	r.Config = config
	return &r
}

func (r *Retry) Do(doFunc func() error) error {
	err := doFunc()
	if err == nil {
		return nil
	}

	if r.RetryIf != nil {
		if ok := r.RetryIf(err); !ok {
			return err
		}
	}

	return r.onRetry(doFunc, 0, err)
}

func (r *Retry) onRetry(doFunc func() error, attempt uint, err error) error {
	if attempt >= r.Config.Attempts || r.Delay >= r.Config.MaxDelay {
		return err
	}

	attempt++
	time.Sleep(r.Delay)
	if r.ExposeDelay != nil {
		r.Delay = r.ExposeDelay(attempt, err)
	} else {
		r.Delay = r.exposeDelay(attempt, err)
	}

	err = doFunc()
	if err == nil {
		return nil
	}

	if ok := r.RetryIf(err); !ok {
		return err
	}

	return r.onRetry(doFunc, attempt, err)
}

func (r *Retry) exposeDelay(n uint, err error) time.Duration {
	if err == nil {
		return r.Delay
	}

	expBackoff := r.Delay * (1 << (n - 1))

	if expBackoff > r.MaxDelay {
		expBackoff = r.MaxDelay
	}

	jitter := rand.Float64() * r.MaxJitter
	expBackoff = time.Duration(float64(expBackoff) * (1 + jitter))

	return expBackoff
}

func (r *Retry) WithAttempts(n uint) *Retry {
	r.Attempts = n
	return r
}

func (r *Retry) WithDelay(d time.Duration) *Retry {
	r.Delay = d
	return r
}

func (r *Retry) WithMaxDelay(d time.Duration) *Retry {
	r.MaxDelay = d
	return r
}

func (r *Retry) WithMaxJitter(d float64) *Retry {
	r.MaxJitter = d
	return r
}

func (r *Retry) WithOnRetry(onRetry func(uint, error) error) *Retry {
	r.OnRetry = onRetry
	return r
}

func (r *Retry) WithRetryIf(f func(error) bool) *Retry {
	r.RetryIf = f
	return r
}

func (r *Retry) WithDelayType(t func(n uint, err error) time.Duration) *Retry {
	r.ExposeDelay = t
	return r
}

func (r *Retry) WithContext(ctx context.Context) *Retry {
	r.Context = ctx
	return r
}
