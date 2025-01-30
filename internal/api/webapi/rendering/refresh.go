package rendering

type RenderRefreshTarget string

const (
	RefreshTargetNone        RenderRefreshTarget = ""
	RefreshTargetMain        RenderRefreshTarget = "#main"
	RefreshTargetPageContent RenderRefreshTarget = "#page-content"
)

type RenderRefreshLevel int

const (
	BrowserLevelRefresh RenderRefreshLevel = iota
	PageChangeRefresh
	ContentRefresh
	TargetedRefresh
)

type RenderTarget struct {
	RefreshLevel  RenderRefreshLevel
	RefreshTarget RenderRefreshTarget
}
