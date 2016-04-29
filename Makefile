PROG_NAME=aosp_su_daemon

vpath %.h $(INCLUDE)

LDFLAGS=-L.

sources=main su_daemon
sources_go=$(patsubst %,%.go,$(sources))

all: binary shared static

binary:
	go build -o $(PROG_NAME) $(sources_go)

shared:
	go build -buildmode=c-shared -o libaosp_su_daemon.so $(sources_go)

static:
	go build -buildmode=c-archive -o libaosp_su_daemon.a $(sources_go)

test: test.o
	gcc -static -o test $< $(LDFLAGS) -laosp_su_daemon -lpthread

%.o: %.c
	gcc -c $< -o $@
clean:
	rm -f $(PROG_NAME) lib$(PROG_NAME).* test

