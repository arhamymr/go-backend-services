package handlers

import (
	"go-backend-services/types"
	"sync"
)

// you can using syncRWMutex for better performance
var lock sync.Mutex
var response types.Response
