/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
	"gonum.org/v1/gonum/lapack/lapack64"
)

// LQ is a type for creating and using the LQ factorization of a matrix.
type LQ struct {
	lq   *Dense
	tau  []float64
	cond float64
}

func (lq *LQ) updateCond(norm lapack.MatrixNorm) {
	// Since A = L*Q, and Q is orthogonal, we get for the condition number κ
	//  κ(A) := |A| |A^-1| = |L*Q| |(L*Q)^-1| = |L| |Q^T * L^-1|
	//        = |L| |L^-1| = κ(L),
	// where we used that fact that Q^-1 = Q^T. However, this assumes that
	// the matrix norm is invariant under orthogonal transformations which
	// is not the case for CondNorm. Hopefully the error is negligible: κ
	// is only a qualitative measure anyway.
	m := lq.lq.mat.Rows
	work := getFloats(3*m, false)
	iwork := getInts(m, false)
	l := lq.lq.asTriDense(m, blas.NonUnit, blas.Lower)
	v := lapack64.Trcon(norm, l.mat, work, iwork)
	lq.cond = 1 / v
	putFloats(work)
	putInts(iwork)
}

// Factorize computes the LQ factorization of an m×n matrix a where n <= m. The LQ
// factorization always exists even if A is singular.
//
// The LQ decomposition is a factorization of the matrix A such that A = L * Q.
// The matrix Q is an orthonormal n×n matrix, and L is an m×n upper triangular matrix.
// L and Q can be extracted from the LTo and QTo methods.
func (lq *LQ) Factorize(a Matrix) {
	lq.factorize(a, CondNorm)
}

func (lq *LQ) factorize(a Matrix, norm lapack.MatrixNorm) {
	m, n := a.Dims()
	if m > n {
		panic(ErrShape)
	}
	k := min(m, n)
	if lq.lq == nil {
		lq.lq = &Dense{}
	}
	lq.lq.Clone(a)
	work := []float64{0}
	lq.tau = make([]float64, k)
	lapack64.Gelqf(lq.lq.mat, lq.tau, work, -1)
	work = getFloats(int(work[0]), false)
	lapack64.Gelqf(lq.lq.mat, lq.tau, work, len(work))
	putFloats(work)
	lq.updateCond(norm)
}

// Cond returns the condition number for the factorized matrix.
// Cond will panic if the receiver does not contain a successful factorization.
func (lq *LQ) Cond() float64 {
	if lq.lq == nil || lq.lq.IsZero() {
		panic("lq: no decomposition computed")
	}
	return lq.cond
}

// TODO(btracey): Add in the "Reduced" forms for extracting the m×m orthogonal
// and upper triangular matrices.

// LTo extracts the m×n lower trapezoidal matrix from a LQ decomposition.
// If dst is nil, a new matrix is allocated. The resulting L matrix is returned.
func (lq *LQ) LTo(dst *Dense) *Dense {
	r, c := lq.lq.Dims()
	if dst == nil {
		dst = NewDense(r, c, nil)
	} else {
		dst.reuseAs(r, c)
	}

	// Disguise the LQ as a lower triangular.
	t := &TriDense{
		mat: blas64.Triangular{
			N:      r,
			Stride: lq.lq.mat.Stride,
			Data:   lq.lq.mat.Data,
			Uplo:   blas.Lower,
			Diag:   blas.NonUnit,
		},
		cap: lq.lq.capCols,
	}
	dst.Copy(t)

	if r == c {
		return dst
	}
	// Zero right of the triangular.
	for i := 0; i < r; i++ {
		zero(dst.mat.Data[i*dst.mat.Stride+r : i*dst.mat.Stride+c])
	}

	return dst
}

// QTo extracts the n×n orthonormal matrix Q from an LQ decomposition.
// If dst is nil, a new matrix is allocated. The resulting Q matrix is returned.
func (lq *LQ) QTo(dst *Dense) *Dense {
	_, c := lq.lq.Dims()
	if dst == nil {
		dst = NewDense(c, c, nil)
	} else {
		dst.reuseAsZeroed(c, c)
	}
	q := dst.mat

	// Set Q = I.
	ldq := q.Stride
	for i := 0; i < c; i++ {
		q.Data[i*ldq+i] = 1
	}

	// Construct Q from the elementary reflectors.
	work := []float64{0}
	lapack64.Ormlq(blas.Left, blas.NoTrans, lq.lq.mat, lq.tau, q, work, -1)
	work = getFloats(int(work[0]), false)
	lapack64.Ormlq(blas.Left, blas.NoTrans, lq.lq.mat, lq.tau, q, work, len(work))
	putFloats(work)

	return dst
}

// Solve finds a minimum-norm solution to a system of linear equations defined
// by the matrices A and b, where A is an m×n matrix represented in its LQ factorized
// form. If A is singular or near-singular a Condition error is returned.
// See the documentation for Condition for more information.
//
// The minimization problem solved depends on the input parameters.
//  If trans == false, find the minimum norm solution of A * X = B.
//  If trans == true, find X such that ||A*X - B||_2 is minimized.
// The solution matrix, X, is stored in place into x.
func (lq *LQ) Solve(x *Dense, trans bool, b Matrix) error {
	r, c := lq.lq.Dims()
	br, bc := b.Dims()

	// The LQ solve algorithm stores the result in-place into the right hand side.
	// The storage for the answer must be large enough to hold both b and x.
	// However, this method's receiver must be the size of x. Copy b, and then
	// copy the result into x at the end.
	if trans {
		if c != br {
			panic(ErrShape)
		}
		x.reuseAs(r, bc)
	} else {
		if r != br {
			panic(ErrShape)
		}
		x.reuseAs(c, bc)
	}
	// Do not need to worry about overlap between x and b because w has its own
	// independent storage.
	w := getWorkspace(max(r, c), bc, false)
	w.Copy(b)
	t := lq.lq.asTriDense(lq.lq.mat.Rows, blas.NonUnit, blas.Lower).mat
	if trans {
		work := []float64{0}
		lapack64.Ormlq(blas.Left, blas.NoTrans, lq.lq.mat, lq.tau, w.mat, work, -1)
		work = getFloats(int(work[0]), false)
		lapack64.Ormlq(blas.Left, blas.NoTrans, lq.lq.mat, lq.tau, w.mat, work, len(work))
		putFloats(work)

		ok := lapack64.Trtrs(blas.Trans, t, w.mat)
		if !ok {
			return Condition(math.Inf(1))
		}
	} else {
		ok := lapack64.Trtrs(blas.NoTrans, t, w.mat)
		if !ok {
			return Condition(math.Inf(1))
		}
		for i := r; i < c; i++ {
			zero(w.mat.Data[i*w.mat.Stride : i*w.mat.Stride+bc])
		}
		work := []float64{0}
		lapack64.Ormlq(blas.Left, blas.Trans, lq.lq.mat, lq.tau, w.mat, work, -1)
		work = getFloats(int(work[0]), false)
		lapack64.Ormlq(blas.Left, blas.Trans, lq.lq.mat, lq.tau, w.mat, work, len(work))
		putFloats(work)
	}
	// x was set above to be the correct size for the result.
	x.Copy(w)
	putWorkspace(w)
	if lq.cond > ConditionTolerance {
		return Condition(lq.cond)
	}
	return nil
}

// SolveVec finds a minimum-norm solution to a system of linear equations.
// See LQ.Solve for the full documentation.
func (lq *LQ) SolveVec(x *VecDense, trans bool, b Vector) error {
	r, c := lq.lq.Dims()
	if _, bc := b.Dims(); bc != 1 {
		panic(ErrShape)
	}

	// The Solve implementation is non-trivial, so rather than duplicate the code,
	// instead recast the VecDenses as Dense and call the matrix code.
	bm := Matrix(b)
	if rv, ok := b.(RawVectorer); ok {
		bmat := rv.RawVector()
		if x != b {
			x.checkOverlap(bmat)
		}
		b := VecDense{mat: bmat, n: b.Len()}
		bm = b.asDense()
	}
	if trans {
		x.reuseAs(r)
	} else {
		x.reuseAs(c)
	}
	return lq.Solve(x.asDense(), trans, bm)
}
