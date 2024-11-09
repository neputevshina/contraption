// Code generated by "stringer -type=tagkind -trimprefix=tag"; DO NOT EDIT.

package contraption

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[tagCompound-0]
	_ = x[tagCircle-1]
	_ = x[tagRect-2]
	_ = x[tagRoundrect-3]
	_ = x[tagVoid-4]
	_ = x[tagEquation-5]
	_ = x[tagText-6]
	_ = x[tagCanvas-7]
	_ = x[tagVectorText-8]
	_ = x[tagTopDownText-9]
	_ = x[tagBottomUpText-10]
	_ = x[tagSequence-11]
	_ = x[tagIllustration-12]
	_ = x[tagHalign - -1]
	_ = x[tagValign - -2]
	_ = x[tagFill - -3]
	_ = x[tagStroke - -4]
	_ = x[tagStrokewidth - -5]
	_ = x[tagIdentity - -6]
	_ = x[tagCond - -7]
	_ = x[tagCondfill - -8]
	_ = x[tagCondstroke - -9]
	_ = x[tagBetween - -10]
	_ = x[tagScroll - -11]
	_ = x[tagSource - -12]
	_ = x[tagSink - -13]
	_ = x[tagPosttransform - -101]
	_ = x[tagTransform - -102]
	_ = x[tagScissor - -103]
	_ = x[tagHshrink - -104]
	_ = x[tagVshrink - -105]
	_ = x[tagLimit - -106]
	_ = x[tagVfollow - -107]
	_ = x[tagHfollow - -108]
	_ = x[tagDontDecimate - -109]
	_ = x[tagDecimate - -110]
}

const (
	_tagkind_name_0 = "DecimateDontDecimateHfollowVfollowLimitVshrinkHshrinkScissorTransformPosttransform"
	_tagkind_name_1 = "SinkSourceScrollBetweenCondstrokeCondfillCondIdentityStrokewidthStrokeFillValignHalignCompoundCircleRectRoundrectVoidEquationTextCanvasVectorTextTopDownTextBottomUpTextSequenceIllustration"
)

var (
	_tagkind_index_0 = [...]uint8{0, 8, 20, 27, 34, 39, 46, 53, 60, 69, 82}
	_tagkind_index_1 = [...]uint8{0, 4, 10, 16, 23, 33, 41, 45, 53, 64, 70, 74, 80, 86, 94, 100, 104, 113, 117, 125, 129, 135, 145, 156, 168, 176, 188}
)

func (i tagkind) String() string {
	switch {
	case -110 <= i && i <= -101:
		i -= -110
		return _tagkind_name_0[_tagkind_index_0[i]:_tagkind_index_0[i+1]]
	case -13 <= i && i <= 12:
		i -= -13
		return _tagkind_name_1[_tagkind_index_1[i]:_tagkind_index_1[i+1]]
	default:
		return "tagkind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
