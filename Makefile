include Makefile.standardgo

default: run

check:
	make -s run | tee actual-output.txt && \
	sort actual-output.txt > sorted-output.txt && \
	sort target-output.txt | diff -u - sorted-output.txt && \
	echo OK

clean:
	rm -f actual-output.txt sorted-output.txt
