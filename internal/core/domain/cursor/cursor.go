package cursor

type Cursor struct {
	Row int
	Col int
}

func New(row, col int) *Cursor {
	return &Cursor{
		Row: row,
		Col: col,
	}
}

func (c *Cursor) MoveTo(row, col int) {
	c.Row = row
	c.Col = col
}

func (c *Cursor) MoveUp() {
	c.Row--
}

func (c *Cursor) MoveDown() {
	c.Row++
}

func (c *Cursor) MoveLeft() {
	c.Col--
}

func (c *Cursor) MoveRight() {
	c.Col++
}

func (c *Cursor) MoveToStartLine() {
	c.Col = 0
}

func (c *Cursor) MoveToEndLine(lineLen int) {
	c.Col = lineLen
}
