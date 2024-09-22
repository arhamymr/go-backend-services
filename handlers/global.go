package handlers

import "sync"

// you can using syncRWMutex for better performance
var lock sync.Mutex
