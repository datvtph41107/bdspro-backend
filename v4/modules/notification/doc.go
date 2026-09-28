// Package notification owns the v4 modular notification surface.
// Controllers expose gateway HTTP inbox paths; persistence lives in the v4 DSN
// (see migrate/000032_notification_core). LocalNotificationService serves the
// inbox; LocalNotificationProvider replaces remote writes for in-process callers.
package notification
