package features

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	domain "github.com/rm-ryou/gotor/internal/core/domain/explorer"
	"github.com/rm-ryou/gotor/internal/core/usecase"
	"github.com/rm-ryou/gotor/internal/ui/assets/icon"
	"github.com/rm-ryou/gotor/internal/ui/gio/config"
	designlayout "github.com/rm-ryou/gotor/internal/ui/gio/design/layout"
	"github.com/rm-ryou/gotor/internal/ui/gio/design/system"
)

type ExplorerView struct {
	theme    *system.Theme
	uc       *usecase.Explorer
	cfg      config.Explorer
	layout   *designlayout.Explorer
	list     widget.List
	hList    widget.List
	clickers map[*domain.Node]*widget.Clickable

	addFile      widget.Clickable
	cancelInput  widget.Clickable
	nameEditor   widget.Editor
	inputVisible bool
	pendingFocus bool

	OnError func(error)
}

func NewExplorerView(th *system.Theme, uc *usecase.Explorer, cfg config.Explorer) *ExplorerView {
	return &ExplorerView{
		theme:  th,
		uc:     uc,
		cfg:    cfg,
		layout: designlayout.NewExplorer(int(th.TextSize), cfg.IndentPerDepth, cfg.RowHeightDelta),
		list: widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		hList: widget.List{
			List: layout.List{Axis: layout.Horizontal},
		},
		clickers: make(map[*domain.Node]*widget.Clickable),
	}
}

func (ev *ExplorerView) Layout(gtx layout.Context) layout.Dimensions {
	nodes := ev.uc.Tree().VisibleNodes()
	ev.syncClickables(nodes)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ev.layoutToolbar(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !ev.inputVisible {
				return layout.Dimensions{}
			}
			return ev.layoutFileInput(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			contentWidth := ev.measureContentWidth(gtx, nodes)

			return ev.hList.List.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
				if contentWidth < gtx.Constraints.Max.X {
					contentWidth = gtx.Constraints.Max.X
				}
				gtx.Constraints.Min.X = contentWidth
				gtx.Constraints.Max.X = contentWidth

				return ev.list.List.Layout(
					gtx, len(nodes),
					func(gtx layout.Context, i int) layout.Dimensions {
						return ev.layoutNode(gtx, nodes[i])
					},
				)
			})
		}),
	)
}

func (ev *ExplorerView) HandleNodeClicks(gtx layout.Context) {
	nodes := ev.uc.Tree().VisibleNodes()
	ev.syncClickables(nodes)

	for ev.addFile.Clicked(gtx) {
		ev.inputVisible = true
		ev.pendingFocus = true
		ev.nameEditor.SetText("")
	}

	if ev.inputVisible {
		for {
			e, ok := ev.nameEditor.Update(gtx)
			if !ok {
				break
			}
			if _, ok := e.(widget.SubmitEvent); ok {
				ev.submitNewFile()
			}
		}
		for ev.cancelInput.Clicked(gtx) {
			ev.inputVisible = false
			ev.nameEditor.SetText("")
		}
	}

	for _, node := range nodes {
		c := ev.clickableFor(node)
		if c.Clicked(gtx) {
			if node.IsDir {
				if err := ev.uc.ToggleNode(node); err != nil {
					ev.reportError(err)
				}
			} else {
				if err := ev.uc.SelectFile(node); err != nil {
					ev.reportError(err)
				}
			}
		}
	}
}

func (ev *ExplorerView) InputActive() bool {
	return ev.inputVisible
}

func (ev *ExplorerView) submitNewFile() {
	name := strings.TrimSpace(ev.nameEditor.Text())
	ev.inputVisible = false
	ev.nameEditor.SetText("")
	if name == "" {
		return
	}
	dirPath := ev.uc.Tree().Root().Path
	if err := ev.uc.CreateFile(dirPath, name); err != nil {
		ev.reportError(err)
	}
}

func (ev *ExplorerView) layoutToolbar(gtx layout.Context) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(6),
		Right:  unit.Dp(8),
		Bottom: unit.Dp(6),
		Left:   unit.Dp(8),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceBetween,
		}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Caption(ev.theme.Theme, "Explorer")
				lbl.Color = color.NRGBA{R: 140, G: 140, B: 140, A: 255}
				lbl.MaxLines = 1
				return lbl.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ev.addFile.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					size := gtx.Dp(unit.Dp(22))
					gtx.Constraints.Min = image.Pt(size, size)
					gtx.Constraints.Max = image.Pt(size, size)

					return layout.Stack{}.Layout(gtx,
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							bg := color.NRGBA{R: 52, G: 52, B: 54, A: 255}
							if ev.addFile.Hovered() {
								bg = color.NRGBA{R: 68, G: 68, B: 71, A: 255}
							}
							defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(unit.Dp(5))).Push(gtx.Ops).Pop()
							paint.Fill(gtx.Ops, bg)
							return layout.Dimensions{Size: gtx.Constraints.Min}
						}),
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body1(ev.theme.Theme, "+")
							lbl.Color = ev.theme.Palette.Fg
							lbl.Alignment = text.Middle
							return layout.Center.Layout(gtx, lbl.Layout)
						}),
					)
				})
			}),
		)
	})
}

func (ev *ExplorerView) layoutFileInput(gtx layout.Context) layout.Dimensions {
	ev.nameEditor.Submit = true
	ev.nameEditor.SingleLine = true
	if ev.pendingFocus {
		ev.pendingFocus = false
		gtx.Execute(key.FocusCmd{Tag: &ev.nameEditor})
	}

	return layout.Inset{
		Top:    unit.Dp(2),
		Right:  unit.Dp(6),
		Bottom: unit.Dp(4),
		Left:   unit.Dp(6),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				ed := material.Editor(ev.theme.Theme, &ev.nameEditor, "filename...")
				ed.Color = ev.theme.Palette.Fg
				ed.HintColor = color.NRGBA{R: 100, G: 100, B: 100, A: 255}

				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(unit.Dp(4))).Push(gtx.Ops).Pop()
						paint.Fill(gtx.Ops, color.NRGBA{R: 40, G: 40, B: 42, A: 255})
						return layout.Dimensions{Size: gtx.Constraints.Min}
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    unit.Dp(3),
							Right:  unit.Dp(4),
							Bottom: unit.Dp(3),
							Left:   unit.Dp(4),
						}.Layout(gtx, ed.Layout)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ev.cancelInput.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					size := gtx.Dp(unit.Dp(18))
					gtx.Constraints.Min = image.Pt(size, size)
					gtx.Constraints.Max = image.Pt(size, size)

					lbl := material.Body2(ev.theme.Theme, "×")
					lbl.Color = color.NRGBA{R: 160, G: 160, B: 160, A: 255}
					return layout.Center.Layout(gtx, lbl.Layout)
				})
			}),
		)
	})
}

func (ev *ExplorerView) reportError(err error) {
	if err == nil || ev.OnError == nil {
		return
	}
	ev.OnError(err)
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
						layout.Rigid(layout.Spacer{Width: ev.cfg.NodeGap}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body2(ev.theme.Theme, node.Name)
							lbl.Color = ev.theme.Palette.Fg
							lbl.MaxLines = 1
							lbl.WrapPolicy = text.WrapWords
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

func (ev *ExplorerView) measureContentWidth(gtx layout.Context, nodes []*domain.Node) int {
	maxWidth := 0

	for _, node := range nodes {
		width := ev.measureNodeWidth(gtx, node)
		if width > maxWidth {
			maxWidth = width
		}
	}

	return maxWidth
}

func (ev *ExplorerView) measureNodeWidth(gtx layout.Context, node *domain.Node) int {
	width := gtx.Dp(ev.layout.Indent(node.Depth))
	width += gtx.Dp(ev.layout.RowHeight()) * 2
	width += gtx.Dp(ev.cfg.NodeGap)
	width += ev.measureTextWidth(gtx, node.Name)
	return width
}

func (ev *ExplorerView) measureTextWidth(gtx layout.Context, value string) int {
	var ops op.Ops
	measureGtx := layout.Context{
		Constraints: layout.Constraints{
			Min: image.Point{},
			Max: image.Pt(1_000_000, 1_000_000),
		},
		Metric: gtx.Metric,
		Now:    gtx.Now,
		Locale: gtx.Locale,
		Values: gtx.Values,
		Ops:    &ops,
	}

	lbl := material.Body2(ev.theme.Theme, value)
	lbl.MaxLines = 1
	lbl.WrapPolicy = text.WrapWords

	return lbl.Layout(measureGtx).Size.X
}
