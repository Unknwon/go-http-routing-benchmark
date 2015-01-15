// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package main

import (
	"net/http"
	"testing"
)

// Parse
// https://parse.com/docs/rest#summary
var parseAPI = []route{
	// Objects
	{"POST", "/1/classes/:className"},
	{"GET", "/1/classes/:className/:objectId"},
	{"PUT", "/1/classes/:className/:objectId"},
	{"GET", "/1/classes/:className"},
	{"DELETE", "/1/classes/:className/:objectId"},

	// Users
	{"POST", "/1/users"},
	{"GET", "/1/login"},
	{"GET", "/1/users/:objectId"},
	{"PUT", "/1/users/:objectId"},
	{"GET", "/1/users"},
	{"DELETE", "/1/users/:objectId"},
	{"POST", "/1/requestPasswordReset"},

	// Roles
	{"POST", "/1/roles"},
	{"GET", "/1/roles/:objectId"},
	{"PUT", "/1/roles/:objectId"},
	{"GET", "/1/roles"},
	{"DELETE", "/1/roles/:objectId"},

	// Files
	{"POST", "/1/files/:fileName"},

	// Analytics
	{"POST", "/1/events/:eventName"},

	// Push Notifications
	{"POST", "/1/push"},

	// Installations
	{"POST", "/1/installations"},
	{"GET", "/1/installations/:objectId"},
	{"PUT", "/1/installations/:objectId"},
	{"GET", "/1/installations"},
	{"DELETE", "/1/installations/:objectId"},

	// Cloud Functions
	{"POST", "/1/functions"},
}

var (
	parseBeego      http.Handler
	parseGin        http.Handler
	parseGoji       http.Handler
	parseGorillaMux http.Handler
	parseHttpRouter http.Handler
	parseMartini    http.Handler
	parseMacaron    http.Handler
	parseRevel      http.Handler
)

func init() {
	println("#ParseAPI Routes:", len(parseAPI))

	calcMem("Beego", func() {
		parseBeego = loadBeego(parseAPI)
	})
	calcMem("Gin", func() {
		parseGin = loadGin(parseAPI)
	})
	calcMem("Goji", func() {
		parseGoji = loadGoji(parseAPI)
	})
	calcMem("GorillaMux", func() {
		parseGorillaMux = loadGorillaMux(parseAPI)
	})
	calcMem("HttpRouter", func() {
		parseHttpRouter = loadHttpRouter(parseAPI)
	})
	calcMem("Martini", func() {
		parseMartini = loadMartini(parseAPI)
	})
	calcMem("Macaron", func() {
		parseMacaron = loadMacaron(parseAPI)
	})
	calcMem("Revel", func() {
		parseRevel = loadRevel(parseAPI)
	})

	println()
}

// Static
func BenchmarkBeego_ParseStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/users", nil)
	benchRequest(b, parseBeego, req)
}
func BenchmarkGin_ParseStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/users", nil)
	benchRequest(b, parseGin, req)
}
func BenchmarkGoji_ParseStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/users", nil)
	benchRequest(b, parseGoji, req)
}
func BenchmarkGorillaMux_ParseStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/users", nil)
	benchRequest(b, parseGorillaMux, req)
}
func BenchmarkHttpRouter_ParseStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/users", nil)
	benchRequest(b, parseHttpRouter, req)
}
func BenchmarkMartini_ParseStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/users", nil)
	benchRequest(b, parseMartini, req)
}
func BenchmarkMacaron_ParseStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/users", nil)
	benchRequest(b, parseMacaron, req)
}
func BenchmarkRevel_ParseStatic(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/users", nil)
	benchRequest(b, parseRevel, req)
}

// One Param
func BenchmarkBeego_ParseParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go", nil)
	benchRequest(b, parseBeego, req)
}
func BenchmarkGin_ParseParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go", nil)
	benchRequest(b, parseGin, req)
}
func BenchmarkGoji_ParseParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go", nil)
	benchRequest(b, parseGoji, req)
}
func BenchmarkGorillaMux_ParseParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go", nil)
	benchRequest(b, parseGorillaMux, req)
}
func BenchmarkHttpRouter_ParseParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go", nil)
	benchRequest(b, parseHttpRouter, req)
}
func BenchmarkMartini_ParseParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go", nil)
	benchRequest(b, parseMartini, req)
}
func BenchmarkMacaron_ParseParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go", nil)
	benchRequest(b, parseMacaron, req)
}

func BenchmarkRevel_ParseParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go", nil)
	benchRequest(b, parseRevel, req)
}

// Two Params
func BenchmarkBeego_Parse2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go/123456789", nil)
	benchRequest(b, parseBeego, req)
}
func BenchmarkGin_Parse2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go/123456789", nil)
	benchRequest(b, parseGin, req)
}
func BenchmarkGoji_Parse2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go/123456789", nil)
	benchRequest(b, parseGoji, req)
}
func BenchmarkGorillaMux_Parse2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go/123456789", nil)
	benchRequest(b, parseGorillaMux, req)
}
func BenchmarkHttpRouter_Parse2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go/123456789", nil)
	benchRequest(b, parseHttpRouter, req)
}
func BenchmarkMartini_Parse2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go/123456789", nil)
	benchRequest(b, parseMartini, req)
}
func BenchmarkMacaron_Parse2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go/123456789", nil)
	benchRequest(b, parseMacaron, req)
}
func BenchmarkRevel_Parse2Params(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1/classes/go/123456789", nil)
	benchRequest(b, parseRevel, req)
}

// All Routes
func BenchmarkBeego_ParseAll(b *testing.B) {
	benchRoutes(b, parseBeego, parseAPI)
}
func BenchmarkGin_ParseAll(b *testing.B) {
	benchRoutes(b, parseGin, parseAPI)
}
func BenchmarkGoji_ParseAll(b *testing.B) {
	benchRoutes(b, parseGoji, parseAPI)
}
func BenchmarkGorillaMux_ParseAll(b *testing.B) {
	benchRoutes(b, parseGorillaMux, parseAPI)
}
func BenchmarkHttpRouter_ParseAll(b *testing.B) {
	benchRoutes(b, parseHttpRouter, parseAPI)
}
func BenchmarkMartini_ParseAll(b *testing.B) {
	benchRoutes(b, parseMartini, parseAPI)
}
func BenchmarkMacaron_ParseAll(b *testing.B) {
	benchRoutes(b, parseMacaron, parseAPI)
}
func BenchmarkRevel_ParseAll(b *testing.B) {
	benchRoutes(b, parseRevel, parseAPI)
}
