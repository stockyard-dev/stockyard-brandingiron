package server

import "github.com/stockyard-dev/stockyard-brandingiron/internal/license"

type Limits struct {
	MaxTemplates   int
	CustomFonts    bool
	RemoveWatermark bool
}

var freeLimits = Limits{MaxTemplates: 3, CustomFonts: false, RemoveWatermark: false}
var proLimits = Limits{MaxTemplates: 0, CustomFonts: true, RemoveWatermark: true}

func LimitsFor(info *license.Info) Limits {
	if info != nil && info.IsPro() {
		return proLimits
	}
	return freeLimits
}
