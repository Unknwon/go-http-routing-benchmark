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

## Results

Benchmark System:

 * 2 GHz Intel Core i7
 * 8 GB 1600 MHz DDR3
 * go version go1.7.5 darwin/amd64
 * Mac OS X 10.12.3

### Memory Consumption

Besides the micro-benchmarks, there are 3 sets of benchmarks where we play around with clones of some real-world APIs, and one benchmark with static routes only, to allow a comparison with [http.ServeMux](http://golang.org/pkg/net/http/#ServeMux).
The following table shows the memory required only for loading the routing structure for the respective API.
The best 3 values for each test are bold. I'm pretty sure you can detect a pattern :wink:

| Router       | Static    | GitHub     | Google+   | Parse     |
|:-------------|----------:|-----------:|----------:|----------:|
| HttpServeMux |__17344 B__|         -  |        -  |        -  |
| Beego        |  93896 B  |  145080 B  |   9840 B  |  18544 B  |
| Goji         |__27200 B__| __86088 B__| __2912 B__| __5232 B__|
| Gorilla Mux  | 668496 B  | 1494864 B  |  71072 B  | 122184 B  |
| Martini      | 309040 B  |  476960 B  |  23904 B  |  45952 B  |
| Macaron      |__37856 B__|__132536 B__| __8656 B__|__13648 B__|

### Static Routes

The `Static` benchmark is not really a clone of a real-world API. It is just a collection of random static paths inspired by the structure of the Go directory. It might not be a realistic URL-structure.

The only intention of this benchmark is to allow a comparison with the default router of Go's net/http package, [http.ServeMux](http://golang.org/pkg/net/http/#ServeMux), which is limited to static routes and does not support parameters in the route pattern.

In the `StaticAll` benchmark each of 157 URLs is called once per repetition (op, *operation*). If you are unfamiliar with the `go test -bench` tool, the first number is the number of repetitions the `go test` tool made, to get a test running long enough for measurements. The second column shows the time in nanoseconds that a single repetition takes. The third number is the amount of heap memory allocated in bytes, the last one the average number of allocations made per repetition.

The logs below show, that http.ServeMux has only medium performance, compared to more feature-rich routers. The fastest router only needs 1.8% of the time http.ServeMux needs.

```
BenchmarkHttpServeMux_StaticAll  	    2000	    856194 ns/op	      96 B/op	       8 allocs/op
BenchmarkBeego_StaticAll         	   10000	    210652 ns/op	   57776 B/op	     628 allocs/op
BenchmarkGoji_StaticAll          	   20000	     62771 ns/op	       0 B/op	       0 allocs/op
BenchmarkGorillaMux_StaticAll    	    1000	   1770567 ns/op	  108112 B/op	    1421 allocs/op
BenchmarkMartini_StaticAll       	     500	   2750423 ns/op	  132819 B/op	    2178 allocs/op
BenchmarkMacaron_StaticAll       	    5000	    299464 ns/op	  120577 B/op	    1413 allocs/op
```

### Micro Benchmarks

The following benchmarks measure the cost of some very basic operations.

In the first benchmark, only a single route, containing a parameter, is loaded into the routers. Then a request for a URL matching this pattern is made and the router has to call the respective registered handler function.

```
BenchmarkBeego_Param             	 1000000	      1466 ns/op	     368 B/op	       4 allocs/op
BenchmarkGoji_Param              	 2000000	       790 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_Param        	  500000	      2935 ns/op	     992 B/op	      10 allocs/op
BenchmarkMartini_Param           	  300000	      6086 ns/op	    1104 B/op	      11 allocs/op
BenchmarkMacaron_Param           	 1000000	      2731 ns/op	    1072 B/op	      11 allocs/op
```

Same as before, but now with multiple parameters, all in the same single route. The intention is to see how the routers scale with the number of parameters. The values of the parameters must be passed to the handler function somehow, which requires allocations. Let's see how clever the routers solve this task with a route containing 5 and 20 parameters:

```
BenchmarkBeego_Param5            	 1000000	      1539 ns/op	     368 B/op	       4 allocs/op
BenchmarkGoji_Param5             	 1000000	      1070 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_Param5       	  300000	      4484 ns/op	    1056 B/op	      10 allocs/op
BenchmarkMartini_Param5          	  200000	      6809 ns/op	    1232 B/op	      11 allocs/op
BenchmarkMacaron_Param5          	  500000	      2721 ns/op	    1072 B/op	      11 allocs/op

BenchmarkBeego_Param20           	  500000	      2959 ns/op	     368 B/op	       4 allocs/op
BenchmarkGoji_Param20            	  300000	      3355 ns/op	    1247 B/op	       2 allocs/op
BenchmarkGorillaMux_Param20      	  200000	     11170 ns/op	    3163 B/op	      12 allocs/op
BenchmarkMartini_Param20         	  100000	     13837 ns/op	    3597 B/op	      13 allocs/op
BenchmarkMacaron_Param20         	  200000	      8527 ns/op	    2923 B/op	      13 allocs/op
```

Now let's see how expensive it is to access a parameter. The handler function reads the value (by the name of the parameter, e.g. with a map lookup; depends on the router) and writes it to our [web scale storage](https://www.youtube.com/watch?v=b2F-DItXtZs) (`/dev/null`).

```
BenchmarkBeego_ParamWrite        	 1000000	      1479 ns/op	     376 B/op	       5 allocs/op
BenchmarkGoji_ParamWrite         	 2000000	       927 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_ParamWrite   	  500000	      2972 ns/op	    1000 B/op	      11 allocs/op
BenchmarkMartini_ParamWrite      	  200000	      6831 ns/op	    1208 B/op	      15 allocs/op
BenchmarkMacaron_ParamWrite      	  500000	      3360 ns/op	    1160 B/op	      14 allocs/op
```

### [Parse.com](https://parse.com/docs/rest#summary)

Enough of the micro benchmark stuff. Let's play a bit with real APIs. In the first set of benchmarks, we use a clone of the structure of [Parse](https://parse.com)'s decent medium-sized REST API, consisting of 26 routes.

The tasks are 1.) routing a static URL (no parameters), 2.) routing a URL containing 1 parameter, 3.) same with 2 parameters, 4.) route all of the routes once (like the StaticAll benchmark, but the routes now contain parameters).

Worth noting is, that the requested route might be a good case for some routing algorithms, while it is a bad case for another algorithm. The values might vary slightly depending on the selected route.

```
BenchmarkBeego_ParseStatic       	 1000000	      1202 ns/op	     368 B/op	       4 allocs/op
BenchmarkGoji_ParseStatic        	 5000000	       275 ns/op	       0 B/op	       0 allocs/op
BenchmarkGorillaMux_ParseStatic  	 1000000	      2604 ns/op	     688 B/op	       9 allocs/op
BenchmarkMartini_ParseStatic     	  300000	      5311 ns/op	     784 B/op	      10 allocs/op
BenchmarkMacaron_ParseStatic     	 1000000	      1963 ns/op	     768 B/op	       9 allocs/op

BenchmarkBeego_ParseParam        	 1000000	      1130 ns/op	     368 B/op	       4 allocs/op
BenchmarkGoji_ParseParam         	 2000000	       761 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_ParseParam   	  500000	      2797 ns/op	     992 B/op	      10 allocs/op
BenchmarkMartini_ParseParam      	  300000	      5776 ns/op	    1104 B/op	      11 allocs/op
BenchmarkMacaron_ParseParam      	 1000000	      2339 ns/op	    1056 B/op	      10 allocs/op

BenchmarkBeego_Parse2Params      	 1000000	      1268 ns/op	     368 B/op	       4 allocs/op
BenchmarkGoji_Parse2Params       	 2000000	       819 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_Parse2Params 	  500000	      3120 ns/op	    1008 B/op	      10 allocs/op
BenchmarkMartini_Parse2Params    	  300000	      5705 ns/op	    1136 B/op	      11 allocs/op
BenchmarkMacaron_Parse2Params    	 1000000	      2395 ns/op	    1056 B/op	      10 allocs/op

BenchmarkBeego_ParseAll          	   50000	     32461 ns/op	    9568 B/op	     104 allocs/op
BenchmarkGoji_ParseAll           	  100000	     16428 ns/op	    5376 B/op	      32 allocs/op
BenchmarkGorillaMux_ParseAll     	   10000	    112740 ns/op	   22800 B/op	     250 allocs/op
BenchmarkMartini_ParseAll        	   10000	    148264 ns/op	   25600 B/op	     276 allocs/op
BenchmarkMacaron_ParseAll        	   30000	     57537 ns/op	   24576 B/op	     250 allocs/o
```


### [GitHub](http://developer.github.com/v3/)

The GitHub API is rather large, consisting of 203 routes. The tasks are basically the same as in the benchmarks before.

```
BenchmarkBeego_GithubStatic      	 1000000	      1394 ns/op	     368 B/op	       4 allocs/op
BenchmarkGoji_GithubStatic       	 5000000	       269 ns/op	       0 B/op	       0 allocs/op
BenchmarkGorillaMux_GithubStatic 	  100000	     19705 ns/op	     688 B/op	       9 allocs/op
BenchmarkMartini_GithubStatic    	  100000	     16992 ns/op	     784 B/op	      10 allocs/op
BenchmarkMacaron_GithubStatic    	 1000000	      2124 ns/op	     768 B/op	       9 allocs/op

BenchmarkBeego_GithubParam       	 1000000	      1522 ns/op	     368 B/op	       4 allocs/op
BenchmarkGoji_GithubParam        	 1000000	      1191 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_GithubParam  	  200000	     11385 ns/op	    1008 B/op	      10 allocs/op
BenchmarkMartini_GithubParam     	  100000	     13737 ns/op	    1136 B/op	      11 allocs/op
BenchmarkMacaron_GithubParam     	  500000	      2622 ns/op	    1056 B/op	      10 allocs/op

BenchmarkBeego_GithubAll         	    5000	    287480 ns/op	   74706 B/op	     812 allocs/op
BenchmarkGoji_GithubAll          	    3000	    506663 ns/op	   56113 B/op	     334 allocs/op
BenchmarkGorillaMux_GithubAll    	     200	   6869192 ns/op	  193184 B/op	    1994 allocs/op
BenchmarkMartini_GithubAll       	     300	   5414121 ns/op	  228216 B/op	    2483 allocs/op
BenchmarkMacaron_GithubAll       	    3000	    512117 ns/op	  204387 B/op	    2006 allocs/op
```

### [Google+](https://developers.google.com/+/api/latest/)

Last but not least the Google+ API, consisting of 13 routes. In reality this is just a subset of a much larger API.

```
BenchmarkBeego_GPlusStatic       	 1000000	      1189 ns/op	     368 B/op	       4 allocs/op
BenchmarkGoji_GPlusStatic        	10000000	       208 ns/op	       0 B/op	       0 allocs/op
BenchmarkGorillaMux_GPlusStatic  	 1000000	      1518 ns/op	     688 B/op	       9 allocs/op
BenchmarkMartini_GPlusStatic     	  300000	      4946 ns/op	     784 B/op	      10 allocs/op
BenchmarkMacaron_GPlusStatic     	 1000000	      1871 ns/op	     768 B/op	       9 allocs/op

BenchmarkBeego_GPlusParam        	 1000000	      1284 ns/op	     368 B/op	       4 allocs/op
BenchmarkGoji_GPlusParam         	 2000000	       690 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_GPlusParam   	  500000	      3744 ns/op	     992 B/op	      10 allocs/op
BenchmarkMartini_GPlusParam      	  300000	      5860 ns/op	    1104 B/op	      11 allocs/op
BenchmarkMacaron_GPlusParam      	 1000000	      2279 ns/op	    1056 B/op	      10 allocs/op

BenchmarkBeego_GPlus2Params      	 1000000	      1420 ns/op	     368 B/op	       4 allocs/op
BenchmarkGoji_GPlus2Params       	 1000000	      1007 ns/op	     336 B/op	       2 allocs/op
BenchmarkGorillaMux_GPlus2Params 	  200000	      7933 ns/op	    1008 B/op	      10 allocs/op
BenchmarkMartini_GPlusParam2     	  100000	     13707 ns/op	    1232 B/op	      15 allocs/op
BenchmarkMacaron_GPlusParam2     	 1000000	      2498 ns/op	    1056 B/op	      10 allocs/op

BenchmarkBeego_GPlusAll          	  100000	     16770 ns/op	    4784 B/op	      52 allocs/op
BenchmarkGoji_GPlusAll           	  200000	     10003 ns/op	    3696 B/op	      22 allocs/op
BenchmarkGorillaMux_GPlusAll     	   30000	     55007 ns/op	   12368 B/op	     128 allocs/op
BenchmarkMartini_GPlusAll        	   20000	     97129 ns/op	   14448 B/op	     165 allocs/op
BenchmarkMacaron_GPlusAll        	   50000	     28788 ns/op	   13152 B/op	     128 allocs/op
```
