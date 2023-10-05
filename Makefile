

clean:
	go clean ./...
	go clean -modcache
	rm -rf ./*.log \
	    ./*.svg \
	    ./go.* \
	    pprof \
	    bench \
	    vendor

init:
	GO111MODULE=on go mod init github.com/kpango/glsmt
	GO111MODULE=on go mod tidy

bench: clean init
	sleep 3
	go test -count=1 -timeout=30m -run=NONE -bench . -benchmem

profile: clean init
	rm -rf bench
	mkdir bench
	mkdir pprof
	\
	# go test -count=3 -timeout=30m -run=NONE -bench=BenchmarkChangeOutAllInt_glsmt -benchmem -o pprof/glsmt-test.bin -cpuprofile pprof/cpu-glsmt.out -memprofile pprof/mem-glsmt.out
	go test -count=3 -timeout=30m -run=NONE -bench . -benchmem -o pprof/glsmt-test.bin -cpuprofile pprof/cpu-glsmt.out -memprofile pprof/mem-glsmt.out
	go tool pprof --svg pprof/glsmt-test.bin pprof/cpu-glsmt.out > cpu-glsmt.svg
	go tool pprof --svg pprof/glsmt-test.bin pprof/mem-glsmt.out > mem-glsmt.svg
	\
	mv ./*.svg bench/

profile-web-cpu:
	go tool pprof -http=":6061" \
		pprof/glsmt-test.bin \
		pprof/cpu-glsmt.out

profile-web-mem:
	go tool pprof -http=":6062" \
		pprof/glsmt-test.bin \
		pprof/mem-glsmt.out
