package mapper

var noisyKeys = map[string]struct{}{
	// Filter out video frame counters that change constantly
	"videoElapseFrame":    {},
	"videoElapseInterval": {},
	// Raw video state flags
	"video":  {},
	"video1": {},
}
