package constants

type ContextKey string

var CONTEXT_USER_KEY ContextKey = "CONTEXT_USER_KEY"

// Redis keys
const USER_ONLINE_STATUS_KEY_PREFIX = "user:online:"
const USER_APP_STATE_KEY_PREFIX = "user:app_state:"
