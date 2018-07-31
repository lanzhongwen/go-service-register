# Author: Zhongwen Lan(lanzhongwen@jd.com)
# Created: 2018/06/01
TOPDIR=$(shell while [ ! -d src ] && [ `pwd` != / ] ; do cd ..; done ; echo `pwd`)

CC := g++
LIBS := -lpthread
INCLUDE := -I. -I/usr/local -I$(TOPDIR)/src/github.com/go-service-register
CFLAGS := -Wall $(DEFINES) $(INCLUDE)
CXXFLAGS := $(CFLAGS) -std=c++11

.PHONY: clean 

clean:
	rm -rf *.o register_test client_main

test: register_test
	./register_test

client: client_main
	./client_main

SRCSTEST := $(filter %_test.cpp, $(wildcard *.cpp))
register_test: $(TARGET)
	$(CC) $(CXXFLAGS) $(SRCSTEST) ../register.a -o $@ $(INCLUDE) $(LIBS) -L/usr/local/lib -lgtest -lgtest_main

CLIENT := client.cpp
client_main:
	$(CC) $(CXXFLAGS) $(CLIENT) ../register.a -o $@ $(INCLUDE) $(LIBS) -L/usr/local/lib
