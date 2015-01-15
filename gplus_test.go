// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package main

import (
	"net/http"
	"testing"
)

// Google+
// https://developers.google.com/+/api/latest/
// (in reality this is just a subset of a much larger API)
var gplusAPI = []route{
	// People
	{"GET", "/people/:userId"},
	{"GET", "/people"},
	{"GET", "/activities/:activityId/people/:collection"},
	{"GET", "/people/:userId/people/:collection"},
	{"GET", "/people/:userId/openIdConnect"},

	// Activities
	{"GET", "/people/:userId/activities/:collection"},
	{"GET", "/activities/:activityId"},
	{"GET", "/activities"},

	// Comments
	{"GET", "/activities/:activityId/comments"},
	{"GET", "/comments/:commentId"},

	// Moments
	{"POST", "/people/:userId/moments/:collection"},
	{"GET", "/people/:userId/moments/:collection"},
	{"DELETE", "/moments/:id"},
}

var (
	gplusBeego      http.Handler
	gplusGin        http.Handler
	gplusGoji       http.Handler
	gplusGorillaMux http.Handler
	gplusHttpRouter http.Handler
	gplusMartini    http.Handler
	gplusMacaron    http.Handler
	gplusRevel      http.Handler
)

func init() {
	println("#GPlusAPI Routes:", len(gplusAPI))

	calcMem("Beego", func() {
		gplusBeego = loadBeego(gplusAPI)
	})
	calcMem("Gin", func() {
		gplusGin = loadGin(gplusAPI)
	})
	calcMem("Goji", func() {
		gplusGoji = loadGoji(gplusAPI)
	})
	calcMem("GorillaMux", func() {
		gplusGorillaMux = loadGorillaMux(gplusAPI)
	})
	calcMem("HttpRouter", func() {
		gplusHttpRouter = loadHttpRouter(gplusAPI)
	})
	calcMem("Martini", func() {
		gplusMartini = loadMartini(gplusAPI)
	})
	calcMem("Macaron", func() {
		gplusMacaron = loadMacaron(gplusAPI)
	})
	calcMem("Revel", func() {
		gplusRevel = loadRevel(gplusAPI)
	})

	println()
}

// Static
func BenchmarkBeego_GPlusStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people", nil)
	benchRequest(b, gplusBeego, req)
}
func BenchmarkGin_GPlusStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people", nil)
	benchRequest(b, gplusGin, req)
}
func BenchmarkGoji_GPlusStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people", nil)
	benchRequest(b, gplusGoji, req)
}
func BenchmarkGorillaMux_GPlusStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people", nil)
	benchRequest(b, gplusGorillaMux, req)
}
func BenchmarkHttpRouter_GPlusStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people", nil)
	benchRequest(b, gplusHttpRouter, req)
}
func BenchmarkMartini_GPlusStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people", nil)
	benchRequest(b, gplusMartini, req)
}
func BenchmarkMacaron_GPlusStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people", nil)
	benchRequest(b, gplusMacaron, req)
}
func BenchmarkRevel_GPlusStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people", nil)
	benchRequest(b, gplusRevel, req)
}

// One Param
func BenchmarkBeego_GPlusParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327", nil)
	benchRequest(b, gplusBeego, req)
}
func BenchmarkGin_GPlusParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327", nil)
	benchRequest(b, gplusGin, req)
}
func BenchmarkGoji_GPlusParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327", nil)
	benchRequest(b, gplusGoji, req)
}
func BenchmarkGorillaMux_GPlusParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327", nil)
	benchRequest(b, gplusGorillaMux, req)
}
func BenchmarkHttpRouter_GPlusParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327", nil)
	benchRequest(b, gplusHttpRouter, req)
}
func BenchmarkMartini_GPlusParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327", nil)
	benchRequest(b, gplusMartini, req)
}
func BenchmarkMacaron_GPlusParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327", nil)
	benchRequest(b, gplusMacaron, req)
}
func BenchmarkRevel_GPlusParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327", nil)
	benchRequest(b, gplusRevel, req)
}

// Two Params
func BenchmarkBeego_GPlus2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327/activities/123456789", nil)
	benchRequest(b, gplusBeego, req)
}
func BenchmarkGin_GPlus2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327/activities/123456789", nil)
	benchRequest(b, gplusGin, req)
}
func BenchmarkGoji_GPlus2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327/activities/123456789", nil)
	benchRequest(b, gplusGoji, req)
}
func BenchmarkGorillaMux_GPlus2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327/activities/123456789", nil)
	benchRequest(b, gplusGorillaMux, req)
}
func BenchmarkHttpRouter_GPlus2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327/activities/123456789", nil)
	benchRequest(b, gplusHttpRouter, req)
}
func BenchmarkMartini_GPlusParam2(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327/activities/123456789", nil)
	benchRequest(b, gplusMartini, req)
}
func BenchmarkMacaron_GPlusParam2(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327/activities/123456789", nil)
	benchRequest(b, gplusMacaron, req)
}
func BenchmarkRevel_GPlus2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/people/118051310819094153327/activities/123456789", nil)
	benchRequest(b, gplusRevel, req)
}

// All Routes
func BenchmarkBeego_GPlusAll(b *testing.B) {
	benchRoutes(b, gplusBeego, gplusAPI)
}
func BenchmarkGin_GPlusAll(b *testing.B) {
	benchRoutes(b, gplusGin, gplusAPI)
}
func BenchmarkGoji_GPlusAll(b *testing.B) {
	benchRoutes(b, gplusGoji, gplusAPI)
}
func BenchmarkGorillaMux_GPlusAll(b *testing.B) {
	benchRoutes(b, gplusGorillaMux, gplusAPI)
}
func BenchmarkHttpRouter_GPlusAll(b *testing.B) {
	benchRoutes(b, gplusHttpRouter, gplusAPI)
}
func BenchmarkMartini_GPlusAll(b *testing.B) {
	benchRoutes(b, gplusMartini, gplusAPI)
}
func BenchmarkMacaron_GPlusAll(b *testing.B) {
	benchRoutes(b, gplusMacaron, gplusAPI)
}
func BenchmarkRevel_GPlusAll(b *testing.B) {
	benchRoutes(b, gplusRevel, gplusAPI)
}
