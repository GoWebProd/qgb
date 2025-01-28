# Quick Go Builder

Простой конструктор SQL запросов для PostgreSQL c возможностью подставлять модели как параметры

## Бенчмарки

```
goos: darwin
goarch: arm64
cpu: Apple M3 Pro
BenchmarkFullBuild-12     	 1471616	       810.1 ns/op	    1152 B/op	      18 allocs/op
BenchmarkPrepare-12       	 7400949	       161.6 ns/op	     344 B/op	       3 allocs/op
BenchmarkSquirrel-12      	  292369	      4069 ns/op	    3722 B/op	      79 allocs/op
BenchmarkSQLBuilder-12    	 1311027	       919.0 ns/op	    1136 B/op	      31 allocs/op
BenchmarkBob-12           	  543602	      2184 ns/op	    2489 B/op	      67 allocs/op
PASS
```
