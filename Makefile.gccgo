default: run

a.out: main.go
	gccgo -o $@ $<

run: a.out
	./a.out

check: a.out
	./a.out | tee actual-output.txt && \
	sort actual-output.txt > sorted-output.txt && \
	diff -u target-output.txt sorted-output.txt && echo OK

clean:
	rm -f actual-output.txt sorted-output.txt
