package utils

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/goerr"
)

func HandleError(ctx context.Context, msg string, err error) {
	// Sending error to Sentry
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		if goErr := goerr.Unwrap(err); goErr != nil {
			for k, v := range goErr.Values() {
				scope.SetExtra(k, v)
			}
		}
	})
	evID := hub.CaptureException(err)

	CtxLogger(ctx).Error(msg, ErrLog(err), "sentry.EventID", evID)
}
