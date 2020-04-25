package engine

import (
	"fmt"
	"math"
)

/**
 * A minified affine transform.
 *  column major (form used by this class)
 *     x'   |a c tx| |x|
 *     y' = |b d ty| |y|
 *     1    |0 0  1| |1|
 *  or
 *  Row major
 *                           |a  b   0|
 *     |x' y' 1| = |x y 1| x |c  d   0|
 *                           |tx ty  1|
 *
 */
type AffineTransform struct {
	a, b, c, d float64
	tx, ty     float64
}

func NewAffineTransform() *AffineTransform {
	at := new(AffineTransform)
	at.ToIdentity()
	return at
}

func (at *AffineTransform) AsTranslate(tx, ty float64) *AffineTransform {
	//  t := new (1.0, 0.0, 0.0, 1.0, tx, ty)
	t := NewAffineTransform()
	t.SetToTranslate(tx, ty)
	return t
}

func (at *AffineTransform) AsScale(sx, sy float64) *AffineTransform {
	t := NewAffineTransform()
	t.SetToScale(sx, sy)
	return t
}

// ----------------------------------------------------------
// Operators
// ----------------------------------------------------------
//   bool operator ==(AffineTransform t) {
//     return (a == t.a && b == t.b && c == t.c && d == t.d && tx == t.tx && ty == t.ty);
//   }

// ----------------------------------------------------------
// Methods
// ----------------------------------------------------------
func (at *AffineTransform) ToIdentity() {
	at.a = 1.0
	at.d = 1.0
	at.b = 0.0
	at.c = 0.0
	at.tx = 0.0
	at.ty = 0.0
}

func (at *AffineTransform) ApplyToVector(point *Vector3, out *Vector3) {
	out.Set2Components(
		(at.a*point.X)+(at.c*point.Y),
		(at.b*point.X)+(at.d*point.Y),
	)
}

func (at *AffineTransform) ApplyTo(point *Vector3, out *Vector3) {
	out.Set2Components(
		(at.a*point.X)+(at.c*point.Y)+at.tx,
		(at.b*point.X)+(at.d*point.Y)+at.ty,
	)
}

//   func (at *AffineTransform) ApplyToSize(Size size) {
//     size.width = (a * size.width + c * size.height).toInt();
//     size.height = (b * size.width + d * size.height).toInt();
//   }

//   func (at *AffineTransform) ApplyToRect(MutableRectangle<double> rect) {
//     double top    = rect.bottom;
//     double left   = rect.left;
//     double right  = rect.width;
//     double bottom = rect.height;

//     Point topLeft = new Point(left, top);
//     Point topRight = new Point(right, top);
//     Point bottomLeft = new Point(left, bottom);
//     Point bottomRight = new Point(right, bottom);
//     ApplyToPoint(topLeft);
//     ApplyToPoint(topRight);
//     ApplyToPoint(bottomLeft);
//     ApplyToPoint(bottomRight);

//     double minX = math.min(math.min(topLeft.x, topRight.x), math.min(bottomLeft.x, bottomRight.x));
//     double maxX = math.max(math.max(topLeft.x, topRight.x), math.max(bottomLeft.x, bottomRight.x));
//     double minY = math.min(math.min(topLeft.y, topRight.y), math.min(bottomLeft.y, bottomRight.y));
//     double maxY = math.max(math.max(topLeft.y, topRight.y), math.max(bottomLeft.y, bottomRight.y));
//   }

func (at *AffineTransform) Set(a, b, c, d, tx, ty float64) {
	at.a = a
	at.b = b
	at.c = c
	at.d = d
	at.tx = tx
	at.ty = ty
}

func (at *AffineTransform) SetWithAT(t *AffineTransform) {
	at.a = t.a
	at.b = t.b
	at.c = t.c
	at.d = t.d
	at.tx = t.tx
	at.ty = t.ty
}

/// Concatenate translation
func (at *AffineTransform) Translate(x, y float64) {
	at.tx += (at.a * x) + (at.c * y)
	at.ty += (at.b * x) + (at.d * y)
}

func (at *AffineTransform) SetToTranslate(tx, ty float64) {
	at.Set(1.0, 0.0, 0.0, 1.0, tx, ty)
}

func (at *AffineTransform) SetToScale(sx, sy float64) {
	at.Set(sx, 0.0, 0.0, sy, 0.0, 0.0)
}

/// Concatenate scale
func (at *AffineTransform) Scale(sx, sy float64) {
	at.a *= sx
	at.b *= sx
	at.c *= sy
	at.d *= sy
}

// If Y axis is downward (default for SDL and Image) then:
// +angle yields a CW rotation
// -angle yeilds a CCW rotation.
//
// else
// -angle yields a CW rotation
// +angle yeilds a CCW rotation.
func (at *AffineTransform) SetToRotate(angle float64) {
	at.a = math.Cos(angle)
	at.b = -math.Sin(angle)
	at.c = math.Sin(angle)
	at.d = math.Cos(angle)
	at.tx = 0
	at.ty = 0
}

/**
 * Concatinate a rotation (radians) onto this transform.
 *
 * Rotation is just a matter of perspective. A CW rotation can be seen as
 * CCW depending on what you are talking about rotating. For example,
 * if the coordinate system is thought as rotating CCW then objects are
 * seen as rotating CW, and that is what the 2x2 matrix below represents
 * with Canvas2D. It is also the frame of reference we use.
 *     |cos  -sin|   object appears to rotate CW.
 *     |sin   cos|
 *
 * In the matrix below the object appears to rotate CCW.
 *     |cos  sin|
 *     |-sin cos|
 *
 *     |a  c|    |cos  -sin|
 *     |b  d|  x |sin   cos|
 *
 */
func (at *AffineTransform) Rotate(angle float64) {
	sin := math.Sin(angle)
	cos := math.Cos(angle)
	_a := at.a
	_b := at.b
	_c := at.c
	_d := at.d

	/*
	 * |a1 c1|   |a2 c2|   |a1a2 + a1b2, a1c2 + c1d2|
	 * |b1 d1| x |b2 d2| = |b1a2 + d1b2, b1c2 + d1d2|
	 *
	 * |_a, _c|   |cos, -sin|   |_acos + _csin, _a(-sin) + _ccos|
	 * |_b, _d| x |sin,  cos| = |_bcos + _dsin, _b(-sin) + _dcos|
	 */
	at.a = _a*cos + _c*sin
	at.b = _b*cos + _d*sin
	at.c = _c*cos - _a*sin
	at.d = _d*cos - _b*sin

	/*
	 * |_a, _c|   |cos,  sin|   |_acos + _c(-sin), _a(sin) + _ccos|
	 * |_b, _d| x |-sin, cos| = |_bcos + _d(-sin), _b(sin) + _dcos|
	 */
	//    a = _a * cos - _c * sin;
	//    b = _b * cos - _d * sin;
	//    c = _c * cos + _a * sin;
	//    d = _d * cos + _b * sin;
}

/**
 * A minified affine transform.
 *     |a c tx|
 *     |b d ty|
 *     |0 0  1|
 *
 *     |- y -|
 *     |x - -|
 *     |0 0 1|
 */
/// Concatenate skew/shear
/// [x] and [y] are in radians
func (at *AffineTransform) Skew(x, y float64) {
	at.c += math.Tan(y)
	at.b += math.Tan(x)
}

/**
 * A pre multiply order.
 */
func (at *AffineTransform) PreMultiply(t *AffineTransform) {
	_a := at.a
	_b := at.b
	_c := at.c
	_d := at.d
	_tx := at.tx
	_ty := at.ty

	at.a = _a*t.a + _b*t.c
	at.b = _a*t.b + _b*t.d
	at.c = _c*t.a + _d*t.c
	at.d = _c*t.b + _d*t.d
	at.tx = (_tx * t.a) + (_ty * t.c) + t.tx
	at.ty = (_tx * t.b) + (_ty * t.d) + t.ty
}

func (at *AffineTransform) Multiply(t1, t2 *AffineTransform) {
	at.a = t1.a*t2.a + t1.b*t2.c
	at.b = t1.a*t2.b + t1.b*t2.d
	at.c = t1.c*t2.a + t1.d*t2.c
	at.d = t1.c*t2.b + t1.d*t2.d
	at.tx = t1.tx*t2.a + t1.ty*t2.c + t2.tx
	at.ty = t1.tx*t2.b + t1.ty*t2.d + t2.ty
}

func (at *AffineTransform) Invert() {
	determinant := 1.0 / (at.a*at.d - at.b*at.c)
	_a := at.a
	_b := at.b
	_c := at.c
	_d := at.d
	_tx := at.tx
	_ty := at.ty

	at.a = determinant * _d
	at.b = -determinant * _b
	at.c = -determinant * _c
	at.d = determinant * _a
	at.tx = determinant * (_c*_ty - _d*_tx)
	at.ty = determinant * (_b*_tx - _a*_ty)
}

/**
 * Converts either from or to pre or post multiplication.
 *     a c
 *     b d
 * to
 *     a b
 *     c d
 */
func (at *AffineTransform) Transpose() {
	_c := at.c

	at.c = at.b
	at.b = _c
	// tx and ty are implied for partial 2x3 matrices
}

//   func (at *AffineTransform) ExtractUniformScale() {
//     Vector2P p = new Vector2P.withCoords(0.0, 0.0);
//     double length = 0.0;

//     CompApplyAffineTransformTo(1.0, 0.0, p.v, this);
//     length = p.v.length;
//     p.moveToPool();

//     return length;
//   }

/**
 * Reliable as long as the transform does not have a Skew present.
 * http://stackoverflow.com/tags/affinetransform/info
 */
//   double get extractRotation => math.atan2(c, d);

/**
 * Reliable as long as the transform does not have a Skew present.
 * http://stackoverflow.com/tags/affinetransform/info
 */
//   double get extractScaleX => math.sqrt(a*a + b*b);
/**
 * Reliable as long as the transform does not have a Skew present.
 * http://stackoverflow.com/tags/affinetransform/info
 */
//   double get extractScaleY => math.sqrt(d*d + c*c);

// func (at *AffineTransform) String() {
//     StringBuffer s = new StringBuffer();
//     s.writeln("|${a.toStringAsFixed(2)}, ${b.toStringAsFixed(2)}, ${tx.toStringAsFixed(2)}|");
//     s.writeln("|${c.toStringAsFixed(2)}, ${d.toStringAsFixed(2)}, ${ty.toStringAsFixed(2)}|");
//     return s.toString();
//   }

func PointApplyAffineTransform(point *Vector3, t *AffineTransform) *Vector3 {
	return CompApplyAffineTransform(point.X, point.Y, t)
}

func CompApplyAffineTransform(x, y float64, t *AffineTransform) *Vector3 {
	v := NewVector3With2Components(
		(t.a*x)+(t.c*y)+t.tx,
		(t.b*x)+(t.d*y)+t.ty,
	)
	return v
}

func CompApplyAffineTransformTo(x, y float64, out *Vector3, t *AffineTransform) {
	out.Set2Components(
		(t.a*x)+(t.c*y)+t.tx,
		(t.b*x)+(t.d*y)+t.ty,
	)
}

// func (at *AffineTransform) SizeApplyAffineTransform(Size size, AffineTransform t) {
//   size.width = (t.a * size.width + t.c * size.height).toInt();
//   size.height = (t.b * size.width + t.d * size.height).toInt();
// }

/// Returns a poolable object.
// MutableRectangle<double> RectApplyAffineTransform(MutableRectangle<double> rect, AffineTransform at) {
//   double top    = rect.top;
//   double right  = rect.right;
//   double left   = rect.left;
//   double bottom = rect.bottom;

//   Vector2 topLeft = Vectors.v[0];
//   Vector2 topRight = Vectors.v[1];
//   Vector2 bottomLeft = Vectors.v[2];
//   Vector2 bottomRight = Vectors.v[3];

//   topLeft.setValues(
//       (at.a * left) + (at.c * top) + at.tx,
//       (at.b * left) + (at.d * top) + at.ty);
//   topRight.setValues(
//       (at.a * right) + (at.c * top) + at.tx,
//       (at.b * right) + (at.d * top) + at.ty);
//   bottomLeft.setValues(
//       (at.a * left) + (at.c * bottom) + at.tx,
//       (at.b * left) + (at.d * bottom) + at.ty);
//   bottomRight.setValues(
//       (at.a * right) + (at.c * bottom) + at.tx,
//       (at.b * right) + (at.d * bottom) + at.ty);

//   double mm1 = topLeft.x < topRight.x ? topLeft.x : topRight.x;
//   double mm2 = bottomLeft.x < bottomRight.x ? bottomLeft.x : bottomRight.x;
//   double minX = mm1 < mm2 ? mm1 : mm2;

//   mm1 = topLeft.x > topRight.x ? topLeft.x : topRight.x;
//   mm2 = bottomLeft.x > bottomRight.x ? bottomLeft.x : bottomRight.x;
//   double maxX = mm1 > mm2 ? mm1 : mm2;

//   mm1 = topLeft.y < topRight.y ? topLeft.y : topRight.y;
//   mm2 = bottomLeft.y < bottomRight.y ? bottomLeft.y : bottomRight.y;
//   double minY = mm1 < mm2 ? mm1 : mm2;

//   mm1 = topLeft.y > topRight.y ? topLeft.y : topRight.y;
//   mm2 = bottomLeft.y > bottomRight.y ? bottomLeft.y : bottomRight.y;
//   double maxY = mm1 > mm2 ? mm1 : mm2;

//   return new MutableRectangle<double>.withP(minX, minY, (maxX - minX), (maxY - minY));
// }
/**
 * [rect] rectangle to transform.
 * [rectOut] the results.
 * [at] The transform to use.
 */
// void RectApplyAffineTransformTo(MutableRectangle<double> rect, MutableRectangle<double> rectOut, AffineTransform at) {
//   double top    = rect.top;
//   double right  = rect.right;
//   double left   = rect.left;
//   double bottom = rect.bottom;

//   Vector2 topLeft = Vectors.v[0];
//   Vector2 topRight = Vectors.v[1];
//   Vector2 bottomLeft = Vectors.v[2];
//   Vector2 bottomRight = Vectors.v[3];

//   topLeft.setValues(
//       (at.a * left) + (at.c * top) + at.tx,
//       (at.b * left) + (at.d * top) + at.ty);
//   topRight.setValues(
//       (at.a * right) + (at.c * top) + at.tx,
//       (at.b * right) + (at.d * top) + at.ty);
//   bottomLeft.setValues(
//       (at.a * left) + (at.c * bottom) + at.tx,
//       (at.b * left) + (at.d * bottom) + at.ty);
//   bottomRight.setValues(
//       (at.a * right) + (at.c * bottom) + at.tx,
//       (at.b * right) + (at.d * bottom) + at.ty);

//   //CompApplyAffineTransformTo(left, top, topLeft, at);
//   //CompApplyAffineTransformTo(right, top, topRight, at);
//   //CompApplyAffineTransformTo(left, bottom, bottomLeft, at);
//   //CompApplyAffineTransformTo(right, bottom, bottomRight, at);

//   double mm1 = topLeft.x < topRight.x ? topLeft.x : topRight.x;
//   double mm2 = bottomLeft.x < bottomRight.x ? bottomLeft.x : bottomRight.x;
//   double minX = mm1 < mm2 ? mm1 : mm2;
//   //double minX = math.min(math.min(topLeft.x, topRight.x), math.min(bottomLeft.x, bottomRight.x));

//   mm1 = topLeft.x > topRight.x ? topLeft.x : topRight.x;
//   mm2 = bottomLeft.x > bottomRight.x ? bottomLeft.x : bottomRight.x;
//   double maxX = mm1 > mm2 ? mm1 : mm2;
//   //double maxX = math.max(math.max(topLeft.x, topRight.x), math.max(bottomLeft.x, bottomRight.x));

//   mm1 = topLeft.y < topRight.y ? topLeft.y : topRight.y;
//   mm2 = bottomLeft.y < bottomRight.y ? bottomLeft.y : bottomRight.y;
//   double minY = mm1 < mm2 ? mm1 : mm2;
//   //double minY = math.min(math.min(topLeft.y, topRight.y), math.min(bottomLeft.y, bottomRight.y));

//   mm1 = topLeft.y > topRight.y ? topLeft.y : topRight.y;
//   mm2 = bottomLeft.y > bottomRight.y ? bottomLeft.y : bottomRight.y;
//   double maxY = mm1 > mm2 ? mm1 : mm2;
//   //double maxY = math.max(math.max(topLeft.y, topRight.y), math.max(bottomLeft.y, bottomRight.y));

//   rectOut.left = minX;
//   rectOut.bottom = minY;
//   rectOut.width = maxX - minX;
//   rectOut.height = maxY - minY;
// }

/// [rect] is overlayed
// void RectangleApplyAffineTransform(MutableRectangle<double> rect, AffineTransform at) {
//   double top    = rect.top;
//   double right  = rect.right;
//   double left   = rect.left;
//   double bottom = rect.bottom;

//   Vector2 topLeft = Vectors.v[0];
//   Vector2 topRight = Vectors.v[1];
//   Vector2 bottomLeft = Vectors.v[2];
//   Vector2 bottomRight = Vectors.v[3];

//   topLeft.setValues(
//       (at.a * left) + (at.c * top) + at.tx,
//       (at.b * left) + (at.d * top) + at.ty);
//   topRight.setValues(
//       (at.a * right) + (at.c * top) + at.tx,
//       (at.b * right) + (at.d * top) + at.ty);
//   bottomLeft.setValues(
//       (at.a * left) + (at.c * bottom) + at.tx,
//       (at.b * left) + (at.d * bottom) + at.ty);
//   bottomRight.setValues(
//       (at.a * right) + (at.c * bottom) + at.tx,
//       (at.b * right) + (at.d * bottom) + at.ty);

//   double mm1 = topLeft.x < topRight.x ? topLeft.x : topRight.x;
//   double mm2 = bottomLeft.x < bottomRight.x ? bottomLeft.x : bottomRight.x;
//   double minX = mm1 < mm2 ? mm1 : mm2;

//   mm1 = topLeft.x > topRight.x ? topLeft.x : topRight.x;
//   mm2 = bottomLeft.x > bottomRight.x ? bottomLeft.x : bottomRight.x;
//   double maxX = mm1 > mm2 ? mm1 : mm2;

//   mm1 = topLeft.y < topRight.y ? topLeft.y : topRight.y;
//   mm2 = bottomLeft.y < bottomRight.y ? bottomLeft.y : bottomRight.y;
//   double minY = mm1 < mm2 ? mm1 : mm2;

//   mm1 = topLeft.y > topRight.y ? topLeft.y : topRight.y;
//   mm2 = bottomLeft.y > bottomRight.y ? bottomLeft.y : bottomRight.y;
//   double maxY = mm1 > mm2 ? mm1 : mm2;

//   rect.left = minX;
//   rect.bottom = minY;
//   rect.width = maxX - minX;
//   rect.height = maxY - minY;
// }

func AffineTransformTranslate(t *AffineTransform, out *AffineTransform, tx, ty float64) {
	out.Set(
		t.a,
		t.b,
		t.c,
		t.d,
		t.tx+t.a*tx+t.c*ty,
		t.ty+t.b*tx+t.d*ty)
}

func AffineTransformScale(t *AffineTransform, out *AffineTransform, sx, sy float64) {
	out.Set(
		t.a*sx,
		t.b*sx,
		t.c*sy,
		t.d*sy,
		t.tx,
		t.ty)
}

/**
 * Rotation is just a matter of perspective. A CW rotation can be seen as
 * CCW depending on what you are talking about rotating. For example,
 * if the coordinate system is thought as rotating CCW then objects are
 * seen as rotating CW. And that is what the 2x2 matrix below represents
 * with Canvas2D.
 *     |cos  -sin|   object appears to rotate CW.
 *     |sin   cos|
 *
 * In the matrix below the object appears to rotate CCW.
 *     |cos  sin|
 *     |-sin cos|
 *
 *     |a  c|    |cos  -sin|
 *     |b  d|  x |sin   cos|
 *
 */
func AffineTransformRotate(t *AffineTransform, out *AffineTransform, anAngle float64) {
	sin := math.Sin(anAngle)
	cos := math.Cos(anAngle)

	out.Set(
		t.a*cos+t.c*sin,
		t.b*cos+t.d*sin,
		t.c*cos-t.a*sin,
		t.d*cos-t.b*sin,
		t.tx,
		t.ty)

	//  AffineTransform at = new AffineTransform._poolable(
	//      t.a * cos - t.c * sin,
	//      t.b * cos - t.d * sin,
	//      t.c * cos + t.a * sin,
	//      t.d * cos + t.b * sin,
	//      t.tx,
	//      t.ty);

	// return at
}

/**
 * Concatenate `t2' to `t1' and return the result: t' = t1 * t2
 * returns a [Poolable]ed [AffineTransform].
 */
func AffineTransformMultiply(t1 *AffineTransform, t2 *AffineTransform, out *AffineTransform) {
	out.Set(
		t1.a*t2.a+t1.b*t2.c,
		t1.a*t2.b+t1.b*t2.d,
		t1.c*t2.a+t1.d*t2.c,
		t1.c*t2.b+t1.d*t2.d,
		t1.tx*t2.a+t1.ty*t2.c+t2.tx,
		t1.tx*t2.b+t1.ty*t2.d+t2.ty)
	// return t
}

/**
 * Multiply [tA] x [tB] and place in [tB]
 */
func AffineTransformMultiplyTo(tA, tB *AffineTransform) {
	a := tA.a*tB.a + tA.b*tB.c
	b := tA.a*tB.b + tA.b*tB.d
	c := tA.c*tB.a + tA.d*tB.c
	d := tA.c*tB.b + tA.d*tB.d
	tx := tA.tx*tB.a + tA.ty*tB.c + tB.tx
	ty := tA.tx*tB.b + tA.ty*tB.d + tB.ty
	tB.a = a
	tB.b = b
	tB.c = c
	tB.d = d
	tB.tx = tx
	tB.ty = ty
}

/**
 * Multiply [tA] x [tB] and place in [tA]
 * A pre-multiply
 *
 * http://www-evasion.imag.fr/~Francois.Faure/doc/inventorMentor/sgi_html/ch09.html
 */
func AffineTransformMultiplyFrom(tA, tB *AffineTransform) {
	a := tA.a*tB.a + tA.b*tB.c
	b := tA.a*tB.b + tA.b*tB.d
	c := tA.c*tB.a + tA.d*tB.c
	d := tA.c*tB.b + tA.d*tB.d
	tx := tA.tx*tB.a + tA.ty*tB.c + tB.tx
	ty := tA.tx*tB.b + tA.ty*tB.d + tB.ty
	tA.a = a
	tA.b = b
	tA.c = c
	tA.d = d
	tA.tx = tx
	tA.ty = ty
}

/* Return true if `t1' and `t2' are equal, false otherwise. */
func AffineTransformEqualToTransform(t1, t2 *AffineTransform) bool {
	return (t1.a == t2.a && t1.b == t2.b && t1.c == t2.c && t1.d == t2.d && t1.tx == t2.tx && t1.ty == t2.ty)
}

func AffineTransformInvert(t *AffineTransform, out *AffineTransform) {
	determinant := 1.0 / (t.a*t.d - t.b*t.c)

	out.Set(
		determinant*t.d,
		-determinant*t.b,
		-determinant*t.c,
		determinant*t.a,
		determinant*(t.c*t.ty-t.d*t.tx),
		determinant*(t.b*t.tx-t.a*t.ty))

	// return at
}

/**
 * Invert [t] to [to].
 */
func AffineTransformInvertTo(t, to *AffineTransform) {
	determinant := 1.0 / (t.a*t.d - t.b*t.c)

	to.a = determinant * t.d
	to.b = -determinant * t.b
	to.c = -determinant * t.c
	to.d = determinant * t.a
	to.tx = determinant * (t.c*t.ty - t.d*t.tx)
	to.ty = determinant * (t.b*t.tx - t.a*t.ty)
}

/**
 *     |a  c tx|   |cos  -sin|
 *     |b  d ty|   |sin   cos|
 */
func (v AffineTransform) String() string {
	return fmt.Sprintf("|%f, %f, %f|\n|%f, %f, %f|", v.a, v.c, v.tx, v.b, v.d, v.ty)
}
