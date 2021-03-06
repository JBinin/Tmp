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

package mat

// TriKind represents the triangularity of the matrix.
type TriKind bool

const (
	// Upper specifies an upper triangular matrix.
	Upper TriKind = true
	// Lower specifies a lower triangular matrix.
	Lower TriKind = false
)

// SVDKind specifies the treatment of singular vectors during an SVD
// factorization.
type SVDKind int

const (
	// SVDNone specifies that no singular vectors should be computed during
	// the decomposition.
	SVDNone SVDKind = iota + 1
	// SVDThin computes the thin singular vectors, that is, it computes
	//  A = U~ * Σ * V~^T
	// where U~ is of size m×min(m,n), Σ is a diagonal matrix of size min(m,n)×min(m,n)
	// and V~ is of size n×min(m,n).
	SVDThin
	// SVDFull computes the full singular value decomposition,
	//  A = U * Σ * V^T
	// where U is of size m×m, Σ is an m×n diagonal matrix, and V is an n×n matrix.
	SVDFull
)

// GSVDKind specifies the treatment of singular vectors during a GSVD
// factorization.
type GSVDKind int

const (
	// GSVDU specifies that the U singular vectors should be computed during
	// the decomposition.
	GSVDU GSVDKind = 1 << iota
	// GSVDV specifies that the V singular vectors should be computed during
	// the decomposition.
	GSVDV
	// GSVDQ specifies that the Q singular vectors should be computed during
	// the decomposition.
	GSVDQ

	// GSVDNone specifies that no singular vector should be computed during
	// the decomposition.
	GSVDNone
)
