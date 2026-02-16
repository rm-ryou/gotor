package features

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	domain "github.com/rm-ryou/gotor/internal/core/domain/explorer"
	"github.com/rm-ryou/gotor/internal/core/usecase"
	"github.com/rm-ryou/gotor/internal/ui/assets/icon"
	designlayout "github.com/rm-ryou/gotor/internal/ui/gio/design/layout"
	"github.com/rm-ryou/gotor/internal/ui/gio/design/system"
)

const nodeGap = 4

type ExplorerView struct {
	theme    *system.Theme
	uc       *usecase.Explorer
	layout   *designlayout.Explorer
	list     widget.List
	clickers map[*domain.Node]*widget.Clickable
}

func NewExplorerView(th *system.Theme, uc *usecase.Explorer) *ExplorerView {
	return &ExplorerView{
		theme:  th,
		uc:     uc,
		layout: designlayout.NewExplorer(system.DefaultTextSize),
		list: widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		clickers: make(map[*domain.Node]*widget.Clickable),
	}
}

func (ev *ExplorerView) Layout(gtx layout.Context) layout.Dimensions {
	nodes := ev.uc.Tree().VisibleNodes()
	ev.syncClickables(nodes)

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return material.List(ev.theme.Theme, &ev.list).Layout(
				gtx, len(nodes),
				func(gtx layout.Context, i int) layout.Dimensions {
					return ev.layoutNode(gtx, nodes[i])
				},
			)
		}),
	)
}

func (ev *ExplorerView) HandleNodeClicks(gtx layout.Context) {
	nodes := ev.uc.Tree().VisibleNodes()
	ev.syncClickables(nodes)

	for _, node := range nodes {
		c := ev.clickableFor(node)
		if c.Clicked(gtx) {
			if node.IsDir {
				_ = ev.uc.ToggleNode(node)
			} else {
				ev.uc.SelectFile(node)
			}
		}
	}
}

func (ev *ExplorerView) syncClickables(nodes []*domain.Node) {
	visible := make(map[*domain.Node]struct{}, len(nodes))
	for _, node := range nodes {
		visible[node] = struct{}{}
		if _, ok := ev.clickers[node]; !ok {
			ev.clickers[node] = new(widget.Clickable)
		}
	}

	for node := range ev.clickers {
		if _, ok := visible[node]; !ok {
			delete(ev.clickers, node)
		}
	}
}

func (ev *ExplorerView) clickableFor(node *domain.Node) *widget.Clickable {
	c, ok := ev.clickers[node]
	if !ok {
		c = new(widget.Clickable)
		ev.clickers[node] = c
	}

	return c
}

func (ev *ExplorerView) layoutNode(gtx layout.Context, node *domain.Node) layout.Dimensions {
	c := ev.clickableFor(node)

	return c.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, ev.theme.Palette.Bg)
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Left: ev.layout.Indent(node.Depth),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return ev.arrowIcon(gtx, node)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return ev.layoutIcon(gtx, node)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(nodeGap)}.Layout),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body2(ev.theme.Theme, node.Name)
							lbl.Color = ev.theme.Palette.Fg
							return lbl.Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}

func (ev *ExplorerView) arrowIcon(gtx layout.Context, node *domain.Node) layout.Dimensions {
	glyph := nodeArrowGlyph(node)

	return ev.layoutGlyph(gtx, glyph, ev.theme.Palette.Fg)
}

func (ev *ExplorerView) layoutIcon(gtx layout.Context, node *domain.Node) layout.Dimensions {
	glyph, c := nodeIcon(node)
	return ev.layoutGlyph(gtx, glyph, c)
}

func nodeArrowGlyph(node *domain.Node) string {
	if !node.IsDir {
		return ""
	}
	if node.Expanded {
		return icon.ArrowExpanded
	}
	return icon.ArrowCollapsed
}

func nodeIcon(node *domain.Node) (string, color.NRGBA) {
	if node.IsDir {
		if node.Expanded {
			return icon.FolderOpenIcon.Glyph, icon.FolderOpenIcon.Color
		}
		return icon.FolderClosedIcon.Glyph, icon.FolderClosedIcon.Color
	}
	return icon.DefaultFileIcon.Glyph, icon.DefaultFileIcon.Color
}

func (ev *ExplorerView) layoutGlyph(gtx layout.Context, glyph string, c color.NRGBA) layout.Dimensions {
	size := gtx.Dp(ev.layout.RowHeight())
	gtx.Constraints.Min = image.Pt(size, size)
	gtx.Constraints.Max = image.Pt(size, size)

	lbl := material.Body2(ev.theme.Theme, glyph)
	lbl.Color = c

	return layout.Center.Layout(gtx, lbl.Layout)
}
