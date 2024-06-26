// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdaobuilder

package util

import "strings"

func ToUpperFirstLetter(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}
