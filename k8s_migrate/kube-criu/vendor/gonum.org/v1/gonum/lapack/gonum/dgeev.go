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
// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

// Dgeev computes the eigenvalues and, optionally, the left and/or right
// eigenvectors for an n×n real nonsymmetric matrix A.
//
// The right eigenvector v_j of A corresponding to an eigenvalue λ_j
// is defined by
//  A v_j = λ_j v_j,
// and the left eigenvector u_j corresponding to an eigenvalue λ_j is defined by
//  u_j^H A = λ_j u_j^H,
// where u_j^H is the conjugate transpose of u_j.
//
// On return, A will be overwritten and the left and right eigenvectors will be
// stored, respectively, in the columns of the n×n matrices VL and VR in the
// same order as their eigenvalues. If the j-th eigenvalue is real, then
//  u_j = VL[:,j],
//  v_j = VR[:,j],
// and if it is not real, then j and j+1 form a complex conjugate pair and the
// eigenvectors can be recovered as
//  u_j     = VL[:,j] + i*VL[:,j+1],
//  u_{j+1} = VL[:,j] - i*VL[:,j+1],
//  v_j     = VR[:,j] + i*VR[:,j+1],
//  v_{j+1} = VR[:,j] - i*VR[:,j+1],
// where i is the imaginary unit. The computed eigenvectors are normalized to
// have Euclidean norm equal to 1 and largest component real.
//
// Left eigenvectors will be computed only if jobvl == lapack.ComputeLeftEV,
// otherwise jobvl must be lapack.None. Right eigenvectors will be computed
// only if jobvr == lapack.ComputeRightEV, otherwise jobvr must be lapack.None.
// For other values of jobvl and jobvr Dgeev will panic.
//
// wr and wi contain the real and imaginary parts, respectively, of the computed
// eigenvalues. Complex conjugate pairs of eigenvalues appear consecutively with
// the eigenvalue having the positive imaginary part first.
// wr and wi must have length n, and Dgeev will panic otherwise.
//
// work must have length at least lwork and lwork must be at least max(1,4*n) if
// the left or right eigenvectors are computed, and at least max(1,3*n) if no
// eigenvectors are computed. For good performance, lwork must generally be
// larger.  On return, optimal value of lwork will be stored in work[0].
//
// If lwork == -1, instead of performing Dgeev, the function only calculates the
// optimal vaule of lwork and stores it into work[0].
//
// On return, first is the index of the first valid eigenvalue. If first == 0,
// all eigenvalues and eigenvectors have been computed. If first is positive,
// Dgeev failed to compute all the eigenvalues, no eigenvectors have been
// computed and wr[first:] and wi[first:] contain those eigenvalues which have
// converged.
func (impl Implementation) Dgeev(jobvl lapack.LeftEVJob, jobvr lapack.RightEVJob, n int, a []float64, lda int, wr, wi []float64, vl []float64, ldvl int, vr []float64, ldvr int, work []float64, lwork int) (first int) {
	var wantvl bool
	switch jobvl {
	default:
		panic("lapack: invalid LeftEVJob")
	case lapack.ComputeLeftEV:
		wantvl = true
	case lapack.None:
	}
	var wantvr bool
	switch jobvr {
	default:
		panic("lapack: invalid RightEVJob")
	case lapack.ComputeRightEV:
		wantvr = true
	case lapack.None:
	}
	switch {
	case n < 0:
		panic(nLT0)
	case len(work) < lwork:
		panic(shortWork)
	}
	var minwrk int
	if wantvl || wantvr {
		minwrk = max(1, 4*n)
	} else {
		minwrk = max(1, 3*n)
	}
	if lwork != -1 {
		checkMatrix(n, n, a, lda)
		if wantvl {
			checkMatrix(n, n, vl, ldvl)
		}
		if wantvr {
			checkMatrix(n, n, vr, ldvr)
		}
		switch {
		case len(wr) != n:
			panic("lapack: bad length of wr")
		case len(wi) != n:
			panic("lapack: bad length of wi")
		case lwork < minwrk:
			panic(badWork)
		}
	}

	// Quick return if possible.
	if n == 0 {
		work[0] = 1
		return 0
	}

	maxwrk := 2*n + n*impl.Ilaenv(1, "DGEHRD", " ", n, 1, n, 0)
	if wantvl || wantvr {
		maxwrk = max(maxwrk, 2*n+(n-1)*impl.Ilaenv(1, "DORGHR", " ", n, 1, n, -1))
		impl.Dhseqr(lapack.EigenvaluesAndSchur, lapack.OriginalEV, n, 0, n-1,
			nil, 1, nil, nil, nil, 1, work, -1)
		maxwrk = max(maxwrk, max(n+1, n+int(work[0])))
		side := lapack.LeftEV
		if wantvr {
			side = lapack.RightEV
		}
		impl.Dtrevc3(side, lapack.AllEVMulQ, nil, n, nil, 1, nil, 1, nil, 1,
			n, work, -1)
		maxwrk = max(maxwrk, n+int(work[0]))
		maxwrk = max(maxwrk, 4*n)
	} else {
		impl.Dhseqr(lapack.EigenvaluesOnly, lapack.None, n, 0, n-1,
			nil, 1, nil, nil, nil, 1, work, -1)
		maxwrk = max(maxwrk, max(n+1, n+int(work[0])))
	}
	maxwrk = max(maxwrk, minwrk)

	if lwork == -1 {
		work[0] = float64(maxwrk)
		return 0
	}

	// Get machine constants.
	smlnum := math.Sqrt(dlamchS) / dlamchP
	bignum := 1 / smlnum

	// Scale A if max element outside range [smlnum,bignum].
	anrm := impl.Dlange(lapack.MaxAbs, n, n, a, lda, nil)
	var scalea bool
	var cscale float64
	if 0 < anrm && anrm < smlnum {
		scalea = true
		cscale = smlnum
	} else if anrm > bignum {
		scalea = true
		cscale = bignum
	}
	if scalea {
		impl.Dlascl(lapack.General, 0, 0, anrm, cscale, n, n, a, lda)
	}

	// Balance the matrix.
	workbal := work[:n]
	ilo, ihi := impl.Dgebal(lapack.PermuteScale, n, a, lda, workbal)

	// Reduce to upper Hessenberg form.
	iwrk := 2 * n
	tau := work[n : iwrk-1]
	impl.Dgehrd(n, ilo, ihi, a, lda, tau, work[iwrk:], lwork-iwrk)

	var side lapack.EVSide
	if wantvl {
		side = lapack.LeftEV
		// Copy Householder vectors to VL.
		impl.Dlacpy(blas.Lower, n, n, a, lda, vl, ldvl)
		// Generate orthogonal matrix in VL.
		impl.Dorghr(n, ilo, ihi, vl, ldvl, tau, work[iwrk:], lwork-iwrk)
		// Perform QR iteration, accumulating Schur vectors in VL.
		iwrk = n
		first = impl.Dhseqr(lapack.EigenvaluesAndSchur, lapack.OriginalEV, n, ilo, ihi,
			a, lda, wr, wi, vl, ldvl, work[iwrk:], lwork-iwrk)
		if wantvr {
			// Want left and right eigenvectors.
			// Copy Schur vectors to VR.
			side = lapack.RightLeftEV
			impl.Dlacpy(blas.All, n, n, vl, ldvl, vr, ldvr)
		}
	} else if wantvr {
		side = lapack.RightEV
		// Copy Householder vectors to VR.
		impl.Dlacpy(blas.Lower, n, n, a, lda, vr, ldvr)
		// Generate orthogonal matrix in VR.
		impl.Dorghr(n, ilo, ihi, vr, ldvr, tau, work[iwrk:], lwork-iwrk)
		// Perform QR iteration, accumulating Schur vectors in VR.
		iwrk = n
		first = impl.Dhseqr(lapack.EigenvaluesAndSchur, lapack.OriginalEV, n, ilo, ihi,
			a, lda, wr, wi, vr, ldvr, work[iwrk:], lwork-iwrk)
	} else {
		// Compute eigenvalues only.
		iwrk = n
		first = impl.Dhseqr(lapack.EigenvaluesOnly, lapack.None, n, ilo, ihi,
			a, lda, wr, wi, nil, 1, work[iwrk:], lwork-iwrk)
	}

	if first > 0 {
		if scalea {
			// Undo scaling.
			impl.Dlascl(lapack.General, 0, 0, cscale, anrm, n-first, 1, wr[first:], 1)
			impl.Dlascl(lapack.General, 0, 0, cscale, anrm, n-first, 1, wi[first:], 1)
			impl.Dlascl(lapack.General, 0, 0, cscale, anrm, ilo, 1, wr, 1)
			impl.Dlascl(lapack.General, 0, 0, cscale, anrm, ilo, 1, wi, 1)
		}
		work[0] = float64(maxwrk)
		return first
	}

	if wantvl || wantvr {
		// Compute left and/or right eigenvectors.
		impl.Dtrevc3(side, lapack.AllEVMulQ, nil, n,
			a, lda, vl, ldvl, vr, ldvr, n, work[iwrk:], lwork-iwrk)
	}
	bi := blas64.Implementation()
	if wantvl {
		// Undo balancing of left eigenvectors.
		impl.Dgebak(lapack.PermuteScale, lapack.LeftEV, n, ilo, ihi, workbal, n, vl, ldvl)
		// Normalize left eigenvectors and make largest component real.
		for i, wii := range wi {
			if wii < 0 {
				continue
			}
			if wii == 0 {
				scl := 1 / bi.Dnrm2(n, vl[i:], ldvl)
				bi.Dscal(n, scl, vl[i:], ldvl)
				continue
			}
			scl := 1 / impl.Dlapy2(bi.Dnrm2(n, vl[i:], ldvl), bi.Dnrm2(n, vl[i+1:], ldvl))
			bi.Dscal(n, scl, vl[i:], ldvl)
			bi.Dscal(n, scl, vl[i+1:], ldvl)
			for k := 0; k < n; k++ {
				vi := vl[k*ldvl+i]
				vi1 := vl[k*ldvl+i+1]
				work[iwrk+k] = vi*vi + vi1*vi1
			}
			k := bi.Idamax(n, work[iwrk:iwrk+n], 1)
			cs, sn, _ := impl.Dlartg(vl[k*ldvl+i], vl[k*ldvl+i+1])
			bi.Drot(n, vl[i:], ldvl, vl[i+1:], ldvl, cs, sn)
			vl[k*ldvl+i+1] = 0
		}
	}
	if wantvr {
		// Undo balancing of right eigenvectors.
		impl.Dgebak(lapack.PermuteScale, lapack.RightEV, n, ilo, ihi, workbal, n, vr, ldvr)
		// Normalize right eigenvectors and make largest component real.
		for i, wii := range wi {
			if wii < 0 {
				continue
			}
			if wii == 0 {
				scl := 1 / bi.Dnrm2(n, vr[i:], ldvr)
				bi.Dscal(n, scl, vr[i:], ldvr)
				continue
			}
			scl := 1 / impl.Dlapy2(bi.Dnrm2(n, vr[i:], ldvr), bi.Dnrm2(n, vr[i+1:], ldvr))
			bi.Dscal(n, scl, vr[i:], ldvr)
			bi.Dscal(n, scl, vr[i+1:], ldvr)
			for k := 0; k < n; k++ {
				vi := vr[k*ldvr+i]
				vi1 := vr[k*ldvr+i+1]
				work[iwrk+k] = vi*vi + vi1*vi1
			}
			k := bi.Idamax(n, work[iwrk:iwrk+n], 1)
			cs, sn, _ := impl.Dlartg(vr[k*ldvr+i], vr[k*ldvr+i+1])
			bi.Drot(n, vr[i:], ldvr, vr[i+1:], ldvr, cs, sn)
			vr[k*ldvr+i+1] = 0
		}
	}

	if scalea {
		// Undo scaling.
		impl.Dlascl(lapack.General, 0, 0, cscale, anrm, n-first, 1, wr[first:], 1)
		impl.Dlascl(lapack.General, 0, 0, cscale, anrm, n-first, 1, wi[first:], 1)
	}

	work[0] = float64(maxwrk)
	return first
}
