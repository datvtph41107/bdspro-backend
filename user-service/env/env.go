package env

import "os"

// ENV_RUNTIME is the process runtime label used by shared TCP bootstrap helpers.
var ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
