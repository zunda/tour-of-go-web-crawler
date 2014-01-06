default: run

a.out: main.go
	gccgo -o $@ $<

run: a.out
	./a.out

