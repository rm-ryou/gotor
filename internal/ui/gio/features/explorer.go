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
	"github.com/rm-ryou/gotor/internal/ui/gio/design/system"
)

type ExplorerView struct {
	theme   *system.Theme
	uc      *usecase.Explorer
	list    widget.List
	widgets map[*domain.Node]*widget.Clickable
}

func NewExplorerView(th *system.Theme, uc *usecase.Explorer) *ExplorerView {
	return &ExplorerView{
		theme: th,
		uc:    uc,
		list: widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		widgets: make(map[*domain.Node]*widget.Clickable),
	}
}

func (ev *ExplorerView) Layout(gtx layout.Context) layout.Dimensions {
	nodes := ev.uc.Tree().VisibleNodes()

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
	for _, node := range ev.uc.Tree().VisibleNodes() {
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

func (ev *ExplorerView) clickableFor(node *domain.Node) *widget.Clickable {
	c, ok := ev.widgets[node]
	if !ok {
		c = new(widget.Clickable)
		ev.widgets[node] = c
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
					Left: unit.Dp(float32(node.Depth) * 12),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layoutIcon(gtx, node, ev.theme.Theme)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
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

func layoutIcon(gtx layout.Context, node *domain.Node, th *material.Theme) layout.Dimensions {
	var i string
	var c color.NRGBA

	if node.IsDir {
		if node.Expanded {
			i = icon.FolderOpenIcon.Glyph
			c = icon.FolderOpenIcon.Color
		}
		i = icon.FolderClosedIcon.Glyph
		c = icon.FolderClosedIcon.Color
	} else {
		i = icon.DefaultFileIcon.Glyph
		c = icon.DefaultFileIcon.Color
	}

	size := gtx.Dp(unit.Dp(system.DefaultTextSize + 2))
	gtx.Constraints.Min = image.Pt(size, size)
	gtx.Constraints.Max = image.Pt(size, size)

	lbl := material.Body2(th, i)
	lbl.Color = c

	return layout.Center.Layout(gtx, lbl.Layout)
}
