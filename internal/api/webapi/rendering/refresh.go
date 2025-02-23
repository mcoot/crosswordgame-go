package rendering

import "slices"

type RenderRefreshTarget string

const (
	RefreshTargetNone        RenderRefreshTarget = ""
	RefreshTargetMain        RenderRefreshTarget = "main"
	RefreshTargetPageContent RenderRefreshTarget = "page-content"
)

func RefreshTargetSelector(targetId RenderRefreshTarget) string {
	if targetId == RefreshTargetNone {
		return ""
	}
	return "#" + string(targetId)
}

var pageOrContentTargets = []RenderRefreshTarget{
	RefreshTargetNone,
	RefreshTargetMain,
	RefreshTargetPageContent,
}

type RenderTarget struct {
	RefreshTarget RenderRefreshTarget
}

func (r RenderTarget) IsPageOrContentRefresh() bool {
	return slices.Contains(pageOrContentTargets, r.RefreshTarget)
}
