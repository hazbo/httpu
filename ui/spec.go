package ui

type ViewSpec struct {
	x0    int
	y0    int
	x1    int
	y1    int
	Wrap  bool
	Title string
}

func (vs ViewSpec) Dimensions() []int {
	return []int{vs.x0, vs.y0, vs.x1, vs.y1}
}

// UiSpec represents the default dimensions needed to create a user interface
// that matches the builtin termnial user interface. Other Go packages /
// programs can use these dimensions / settings to get a similar look to the
// default.
type UiSpec struct {
	maxX                int
	maxY                int
	RequestViewSpec     ViewSpec
	ResponseViewSpec    ViewSpec
	CmdBarViewSpec      ViewSpec
	StatusCodeViewSpec  ViewSpec
	RequestTimeViewSpec ViewSpec
}

func NewUiSpec(maxX, maxY int) UiSpec {
	return UiSpec{
		maxX: maxX,
		maxY: maxY,
		RequestViewSpec: ViewSpec{
			x0:    0,
			y0:    0,
			x1:    maxX/2 - 2,
			y1:    maxY - 4,
			Wrap:  true,
			Title: "",
		},
		ResponseViewSpec: ViewSpec{
			x0:    maxX/2 + 2,
			y0:    0,
			x1:    maxX - 1,
			y1:    maxY - 4,
			Wrap:  true,
			Title: "",
		},
		CmdBarViewSpec: ViewSpec{
			x0:   0,
			y0:   maxY - 3,
			x1:   maxX/2 - 2,
			y1:   maxY - 1,
			Wrap: false,
		},
		StatusCodeViewSpec: ViewSpec{
			x0:    maxX/2 + 2,
			y0:    maxY - 3,
			x1:    maxX/2 + 11,
			y1:    maxY - 1,
			Wrap:  false,
			Title: "Status",
		},
		RequestTimeViewSpec: ViewSpec{
			x0:    maxX/2 + 2 + 11,
			y0:    maxY - 3,
			x1:    maxX/2 + 22,
			y1:    maxY - 1,
			Wrap:  false,
			Title: "Time",
		}}
}
