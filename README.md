Go HTTP Router Benchmark
========================

This benchmark suite aims to compare the performance of HTTP request routers for [Go](https://golang.org) by implementing the routing structure of some real world APIs.
Some of the APIs are slightly adapted, since they can not be implemented 1:1 in some of the routers.

Of course the tested routers can be used for any kind of HTTP request → handler function routing, not only (REST) APIs.


#### Tested routers & frameworks:

 * [Beego](http://beego.me/)
 * [Goji](https://github.com/zenazn/goji/)
 * [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux)
 * [http.ServeMux](http://golang.org/pkg/net/http/#ServeMux)
 * [Martini](https://github.com/go-martini/martini)
 * [Macaron](https://github.com/Unknwon/macaron)
 * [Revel](https://github.com/revel/revel)

## Motivation

Go is a great language for web applications. Since the [default *request multiplexer*](http://golang.org/pkg/net/http/#ServeMux) of Go's net/http package is very simple and limited, an accordingly high number of HTTP request routers exist.

Unfortunately, most of the (early) routers use pretty bad routing algorithms. Moreover, many of them are very wasteful with memory allocations, which can become a problem in a language with Garbage Collection like Go, since every (heap) allocation results in more work for the Garbage Collector.

Lately more and more bloated frameworks pop up, outdoing one another in the number of features. This benchmark tries to measure their overhead.

Beware that we are comparing apples to oranges here, we compare feature-rich frameworks to packages with simple routing functionality only. But since we are only interested in decent request routing, I think this is not entirely unfair. The frameworks are configured to do as little additional work as possible.

If you care about performance, this benchmark can maybe help you find the right router, which scales with your application.

Personally, I prefer slim and optimized software, which is why I implemented [HttpRouter](https://github.com/julienschmidt/httprouter), which is also tested here. In fact, this benchmark suite started as part of the packages tests, but was then extended to a generic benchmark suite.
So keep in mind, that I am not completely unbiased :relieved:


## Results

Benchmark System:

 * 3.7 GHz Quad-Core Intel Xeon E5
 * 32 GB 1866 MHz DDR3 ECC
 * go version go1.5.2 darwin/amd64
 * Mac OS X 10.11.2

### Memory Consumption

Besides the micro-benchmarks, there are 3 sets of benchmarks where we play around with clones of some real-world APIs, and one benchmark with static routes only, to allow a comparison with [http.ServeMux](http://golang.org/pkg/net/http/#ServeMux).
The following table shows the memory required only for loading the routing structure for the respective API.
The best 3 values for each test are bold. I'm pretty sure you can detect a pattern :wink:

| Router       | Static    | GitHub     | Google+   | Parse     |
|:-------------|----------:|-----------:|----------:|----------:|
| HttpServeMux |__18352 B__|         -  |        -  |        -  |
| Beego        | 105112 B  |  167560 B  |  11480 B  |  21216 B  |
| Goji         |__27200 B__| __88392 B__| __3056 B__| __5456 B__|
| Gorilla Mux  | 668608 B  | 1496704 B  |  71168 B  | 122280 B  |
| Martini      | 309104 B  |  557440 B  |  24000 B  |  45952 B  |
| Macaron      |__37216 B__|__132328 B__| __8640 B__|__13632 B__|
| Revel        |  93632 B  |__141296 B__|__10768 B__|__15488 B__|

### Static Routes

The `Static` benchmark is not really a clone of a real-world API. It is just a collection of random static paths inspired by the structure of the Go directory. It might not be a realistic URL-structure.

The only intention of this benchmark is to allow a comparison with the default router of Go's net/http package, [http.ServeMux](http://golang.org/pkg/net/http/#ServeMux), which is limited to static routes and does not support parameters in the route pattern.

In the `StaticAll` benchmark each of 157 URLs is called once per repetition (op, *operation*). If you are unfamiliar with the `go test -bench` tool, the first number is the number of repetitions the `go test` tool made, to get a test running long enough for measurements. The second column shows the time in nanoseconds that a single repetition takes. The third number is the amount of heap memory allocated in bytes, the last one the average number of allocations made per repetition.

The logs below show, that http.ServeMux has only medium performance, compared to more feature-rich routers. The fastest router only needs 1.8% of the time http.ServeMux needs.

```
BenchmarkHttpServeMux_StaticAll 	    2000	    650124 ns/op	     128 B/op	       8 allocs/op

BenchmarkBeego_StaticAll        	   10000	    199040 ns/op	   57744 B/op	    1099 allocs/op
BenchmarkGoji_StaticAll         	   30000	     58481 ns/op	       0 B/op	       0 allocs/op
BenchmarkGorillaMux_StaticAll   	    1000	   1587820 ns/op	   75488 B/op	    1264 allocs/op
BenchmarkMartini_StaticAll      	     500	   2665141 ns/op	  132819 B/op	    2178 allocs/op
BenchmarkMacaron_StaticAll      	    5000	    287739 ns/op	  118065 B/op	    1256 allocs/op
BenchmarkRevel_StaticAll        	    2000	    769734 ns/op	  203280 B/op	    3611 allocs/op
```

### Micro Benchmarks

The following benchmarks measure the cost of some very basic operations.

In the first benchmark, only a single route, containing a parameter, is loaded into the routers. Then a request for a URL matching this pattern is made and the router has to call the respective registered handler function.

```
BenchmarkBeego_Param            	 1000000	      2118 ns/op	     720 B/op	      10 allocs/op
BenchmarkGoji_Param             	 2000000	       770 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_Param       	  500000	      3229 ns/op	     832 B/op	      11 allocs/op
BenchmarkMartini_Param          	  300000	      5119 ns/op	    1104 B/op	      11 allocs/op
BenchmarkMacaron_Param          	  500000	      3131 ns/op	    1072 B/op	      10 allocs/op
BenchmarkRevel_Param            	  200000	      6514 ns/op	    1664 B/op	      26 allocs/op
```

Same as before, but now with multiple parameters, all in the same single route. The intention is to see how the routers scale with the number of parameters. The values of the parameters must be passed to the handler function somehow, which requires allocations. Let's see how clever the routers solve this task with a route containing 5 and 20 parameters:

```
BenchmarkBeego_Param5           	  500000	      3196 ns/op	     992 B/op	      13 allocs/op
BenchmarkGoji_Param5            	 1000000	      1085 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_Param5      	  300000	      5619 ns/op	    1088 B/op	      19 allocs/op
BenchmarkMartini_Param5         	  300000	      6275 ns/op	    1232 B/op	      11 allocs/op
BenchmarkMacaron_Param5         	  500000	      3390 ns/op	    1072 B/op	      10 allocs/op
BenchmarkRevel_Param5           	  200000	      7302 ns/op	    2016 B/op	      33 allocs/op

BenchmarkBeego_Param20          	  200000	      9628 ns/op	    3867 B/op	      17 allocs/op
BenchmarkGoji_Param20           	  500000	      3312 ns/op	    1246 B/op	       2 allocs/op
BenchmarkGorillaMux_Param20     	  100000	     13992 ns/op	    3931 B/op	      51 allocs/op
BenchmarkMartini_Param20        	  100000	     12223 ns/op	    3595 B/op	      13 allocs/op
BenchmarkMacaron_Param20        	  200000	      8216 ns/op	    2924 B/op	      12 allocs/op
BenchmarkRevel_Param20          	  100000	     15028 ns/op	    5543 B/op	      52 allocs/op
```

Now let's see how expensive it is to access a parameter. The handler function reads the value (by the name of the parameter, e.g. with a map lookup; depends on the router) and writes it to our [web scale storage](https://www.youtube.com/watch?v=b2F-DItXtZs) (`/dev/null`).

```
BenchmarkBeego_ParamWrite       	 1000000	      2224 ns/op	     736 B/op	      11 allocs/op
BenchmarkGoji_ParamWrite        	 2000000	       843 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_ParamWrite  	  500000	      3515 ns/op	     832 B/op	      11 allocs/op
BenchmarkMartini_ParamWrite     	  200000	      6264 ns/op	    1216 B/op	      15 allocs/op
BenchmarkMacaron_ParamWrite     	  500000	      3738 ns/op	    1152 B/op	      13 allocs/op
BenchmarkRevel_ParamWrite       	  200000	      7267 ns/op	    2112 B/op	      31 allocs/op
```

### [Parse.com](https://parse.com/docs/rest#summary)

Enough of the micro benchmark stuff. Let's play a bit with real APIs. In the first set of benchmarks, we use a clone of the structure of [Parse](https://parse.com)'s decent medium-sized REST API, consisting of 26 routes.

The tasks are 1.) routing a static URL (no parameters), 2.) routing a URL containing 1 parameter, 3.) same with 2 parameters, 4.) route all of the routes once (like the StaticAll benchmark, but the routes now contain parameters).

Worth noting is, that the requested route might be a good case for some routing algorithms, while it is a bad case for another algorithm. The values might vary slightly depending on the selected route.

```
BenchmarkBeego_ParseStatic      	 1000000	      1091 ns/op	     368 B/op	       7 allocs/op
BenchmarkGoji_ParseStatic       	 5000000	       243 ns/op	       0 B/op	       0 allocs/op
BenchmarkGorillaMux_ParseStatic 	 1000000	      2501 ns/op	     480 B/op	       8 allocs/op
BenchmarkMartini_ParseStatic    	  500000	      4228 ns/op	     784 B/op	      10 allocs/op
BenchmarkMacaron_ParseStatic    	 1000000	      1775 ns/op	     752 B/op	       8 allocs/op
BenchmarkRevel_ParseStatic      	  300000	      4667 ns/op	    1280 B/op	      23 allocs/op

BenchmarkBeego_ParseParam       	 1000000	      1997 ns/op	     736 B/op	      10 allocs/op
BenchmarkGoji_ParseParam        	 2000000	       818 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_ParseParam  	  500000	      3221 ns/op	     832 B/op	      11 allocs/op
BenchmarkMartini_ParseParam     	  300000	      5018 ns/op	    1104 B/op	      11 allocs/op
BenchmarkMacaron_ParseParam     	 1000000	      2080 ns/op	    1040 B/op	       9 allocs/op
BenchmarkRevel_ParseParam       	  300000	      5219 ns/op	    1696 B/op	      26 allocs/op

BenchmarkGoji_Parse2Params      	 2000000	       739 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_Parse2Params	  500000	      3589 ns/op	     896 B/op	      13 allocs/op
BenchmarkMartini_Parse2Params   	  300000	      4781 ns/op	    1136 B/op	      11 allocs/op
BenchmarkMacaron_Parse2Params   	 1000000	      2208 ns/op	    1040 B/op	       9 allocs/op
BenchmarkRevel_Parse2Params     	  300000	      5518 ns/op	    1760 B/op	      28 allocs/op

BenchmarkBeego_ParseAll         	   30000	     43469 ns/op	   15600 B/op	     233 allocs/op
BenchmarkGoji_ParseAll          	  100000	     15808 ns/op	    5376 B/op	      32 allocs/op
BenchmarkGorillaMux_ParseAll    	   10000	    112735 ns/op	   18304 B/op	     262 allocs/op
BenchmarkMartini_ParseAll       	   10000	    134033 ns/op	   25600 B/op	     276 allocs/op
BenchmarkMacaron_ParseAll       	   30000	     53109 ns/op	   24160 B/op	     224 allocs/op
BenchmarkRevel_ParseAll         	   10000	    134209 ns/op	   40256 B/op	     652 allocs/op
```


### [GitHub](http://developer.github.com/v3/)

The GitHub API is rather large, consisting of 203 routes. The tasks are basically the same as in the benchmarks before.

```
BenchmarkBeego_GithubStatic     	 1000000	      1415 ns/op	     368 B/op	       7 allocs/op
BenchmarkGoji_GithubStatic      	 5000000	       267 ns/op	       0 B/op	       0 allocs/op
BenchmarkGorillaMux_GithubStatic	  100000	     17204 ns/op	     480 B/op	       8 allocs/op
BenchmarkMartini_GithubStatic   	  100000	     16235 ns/op	     784 B/op	      10 allocs/op
BenchmarkMacaron_GithubStatic   	 1000000	      2059 ns/op	     752 B/op	       8 allocs/op
BenchmarkRevel_GithubStatic     	  300000	      4973 ns/op	    1280 B/op	      23 allocs/op

BenchmarkBeego_GithubParam      	 1000000	      2351 ns/op	     784 B/op	      11 allocs/op
BenchmarkGoji_GithubParam       	 1000000	      1122 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_GithubParam 	  200000	     10487 ns/op	     896 B/op	      13 allocs/op
BenchmarkMartini_GithubParam    	  100000	     13191 ns/op	    1136 B/op	      11 allocs/op
BenchmarkMacaron_GithubParam    	 1000000	      2506 ns/op	    1040 B/op	       9 allocs/op
BenchmarkRevel_GithubParam      	  200000	      6093 ns/op	    1776 B/op	      28 allocs/op

BenchmarkBeego_GithubAll        	    3000	    451196 ns/op	  146272 B/op	    2092 allocs/op
BenchmarkGoji_GithubAll         	    3000	    508744 ns/op	   56113 B/op	     334 allocs/op
BenchmarkGorillaMux_GithubAll   	     200	   6135767 ns/op	  167232 B/op	    2469 allocs/op
BenchmarkMartini_GithubAll      	     300	   5610895 ns/op	  228216 B/op	    2483 allocs/op
BenchmarkMacaron_GithubAll      	    3000	    476765 ns/op	  201138 B/op	    1803 allocs/op
BenchmarkRevel_GithubAll        	    2000	   1151285 ns/op	  343936 B/op	    5512 allocs/op
```

### [Google+](https://developers.google.com/+/api/latest/)

Last but not least the Google+ API, consisting of 13 routes. In reality this is just a subset of a much larger API.

```
BenchmarkBeego_GPlusStatic      	 1000000	      1052 ns/op	     352 B/op	       7 allocs/op
BenchmarkGoji_GPlusStatic       	10000000	       194 ns/op	       0 B/op	       0 allocs/op
BenchmarkGorillaMux_GPlusStatic 	 1000000	      1642 ns/op	     480 B/op	       8 allocs/op
BenchmarkMartini_GPlusStatic    	  500000	      3924 ns/op	     784 B/op	      10 allocs/op
BenchmarkMacaron_GPlusStatic    	 1000000	      1882 ns/op	     752 B/op	       8 allocs/op
BenchmarkRevel_GPlusStatic      	  300000	      4809 ns/op	    1264 B/op	      23 allocs/op

BenchmarkBeego_GPlusParam       	 1000000	      1875 ns/op	     720 B/op	      10 allocs/op
BenchmarkGoji_GPlusParam        	 2000000	       688 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_GPlusParam  	  500000	      3803 ns/op	     832 B/op	      11 allocs/op
BenchmarkMartini_GPlusParam     	  300000	      4997 ns/op	    1104 B/op	      11 allocs/op
BenchmarkMacaron_GPlusParam     	 1000000	      2154 ns/op	    1040 B/op	       9 allocs/op
BenchmarkRevel_GPlusParam       	  300000	      5246 ns/op	    1696 B/op	      26 allocs/op

BenchmarkBeego_GPlus2Params     	 1000000	      2161 ns/op	     784 B/op	      11 allocs/op
BenchmarkGoji_GPlus2Params      	 1000000	      1029 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_GPlus2Params	  200000	      8052 ns/op	     896 B/op	      13 allocs/op
BenchmarkMartini_GPlusParam2    	  200000	     12181 ns/op	    1232 B/op	      15 allocs/op
BenchmarkMacaron_GPlusParam2    	 1000000	      2248 ns/op	    1040 B/op	       9 allocs/op
BenchmarkRevel_GPlus2Params     	  300000	      5653 ns/op	    1792 B/op	      28 allocs/op

BenchmarkBeego_GPlusAll         	  100000	     23877 ns/op	    8976 B/op	     129 allocs/op
BenchmarkGoji_GPlusAll          	  200000	      9744 ns/op	    3696 B/op	      22 allocs/op
BenchmarkGorillaMux_GPlusAll    	   30000	     59910 ns/op	   10432 B/op	     147 allocs/op
BenchmarkMartini_GPlusAll       	   20000	     92756 ns/op	   14448 B/op	     165 allocs/op
BenchmarkMacaron_GPlusAll       	   50000	     29502 ns/op	   12944 B/op	     115 allocs/op
BenchmarkRevel_GPlusAll         	   20000	     79955 ns/op	   21552 B/op	     342 allocs/op
```
