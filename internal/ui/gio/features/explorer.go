package features

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	domain "github.com/rm-ryou/gotor/internal/core/domain/explorer"
	"github.com/rm-ryou/gotor/internal/core/usecase"
	"github.com/rm-ryou/gotor/internal/ui/gio/design/system"
)

type ExplorerView struct {
	theme *system.Theme
	uc    *usecase.Explorer
	list  widget.List
}

func NewExplorerView(th *system.Theme, uc *usecase.Explorer) *ExplorerView {
	return &ExplorerView{
		theme: th,
		uc:    uc,
		list: widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
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

func (ev *ExplorerView) layoutNode(gtx layout.Context, node *domain.Node) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, ev.theme.Palette.Bg)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(ev.theme.Theme, node.Name)
					lbl.Color = color.NRGBA{R: 204, G: 204, B: 204, A: 255}
					return lbl.Layout(gtx)
				}),
			)
		}),
	)
}
