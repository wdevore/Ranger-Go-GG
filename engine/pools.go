package engine

// ----------------------------------------------------------------
// AffineTransform Pool
// ----------------------------------------------------------------

// AffineTransformPool is a pool of AffineTransforms
type AffineTransformPool []*AffineTransform

// NewAffineTransformPool creates and adds transforms to the pool
func NewAffineTransformPool(count int) AffineTransformPool {
	p := AffineTransformPool{}
	p.Add(count)
	return p
}

// Add populates the pool with a fixed size
func (q *AffineTransformPool) Add(count int) {
	for i := 0; i < count; i++ {
		at := NewAffineTransform()
		q.Push(at)
	}
}

// Push puts a transform into the pool
func (q *AffineTransformPool) Push(n *AffineTransform) {
	*q = append(*q, n)
}

// Pop pulls a transform from the pool.
func (q *AffineTransformPool) Pop() (n *AffineTransform) {
	n = (*q)[0]
	*q = (*q)[1:]
	return
}

// IsEmpty is true if pool is empty
func (q *AffineTransformPool) IsEmpty() bool {
	return q.plen() == 0
}

func (q *AffineTransformPool) plen() int {
	return len(*q)
}

// ----------------------------------------------------------------
// Vector Pool
// ----------------------------------------------------------------

// VectorPool is a pool of vectors
type VectorPool []*Vector3

// NewVectorPool creates and adds vectors to the pool
func NewVectorPool(count int) VectorPool {
	p := VectorPool{}
	p.Add(count)
	return p
}

// Add populates the pool with a fixed size
func (q *VectorPool) Add(count int) {
	for i := 0; i < count; i++ {
		v := NewVector3()
		q.Push(v)
	}
}

// Push puts a transform into the pool
func (q *VectorPool) Push(n *Vector3) {
	*q = append(*q, n)
}

// Pop pulls a transform from the pool.
func (q *VectorPool) Pop() (n *Vector3) {
	n = (*q)[0]
	*q = (*q)[1:]
	return
}

// IsEmpty is true if pool is empty
func (q *VectorPool) IsEmpty() bool {
	return q.plen() == 0
}

func (q *VectorPool) plen() int {
	return len(*q)
}
