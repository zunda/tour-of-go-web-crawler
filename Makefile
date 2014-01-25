default: run

run: main.go
	if { which go > /dev/null; } then \
		go run main.go; \
	else \
		gccgo -o $@ $< && ./a.out; \
	fi

check:
	make -s run | tee actual-output.txt && \
	sort actual-output.txt > sorted-output.txt && \
	sort target-output.txt | diff -u - sorted-output.txt && \
	echo OK

clean:
	rm -f a.out actual-output.txt sorted-output.txt
