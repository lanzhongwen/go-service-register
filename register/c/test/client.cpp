#include "register/c/register.h"
#include <pthread.h>

#include <chrono>
#include <iostream>
#include <memory>
#include <thread>

int main() {
	std::unique_ptr<std::thread> t(new std::thread([&] {
	   std::cout << "Thread Id: " << t->native_handle() << std::endl;
	   while (1) {
	   	std::this_thread::sleep_for(std::chrono::seconds(1));
	   }
	   std::cout << "Existing..." << std::endl;

	}));

	std::unique_ptr<std::thread> t1(new std::thread([&] {
	   std::cout << "Thread Id: " << t1->native_handle() << std::endl;
	   while (1) {
	   	std::this_thread::sleep_for(std::chrono::seconds(1));
	   }
	   std::cout << "Existing..." << std::endl;

	}));
	std::string etcd_addrs_str("127.0.0.1:2379");
	GoString etcd_addrs{etcd_addrs_str.c_str(), etcd_addrs_str.size()};
	std::string service("/lzw/jd");
	GoString serviceName{service.c_str(), service.size()};
	std::string value("xxxx");
	GoString serviceValue{value.c_str(), value.size()};
	GoInt timeout = 10;
	GoInt ttl = 10;
	int ret = RegisterWithAffinity(etcd_addrs, serviceName, serviceValue, timeout, ttl, 3);
	if (!ret) {
		std::cout << "Failed to Register" << std::endl;
	}
	std::this_thread::sleep_for(std::chrono::seconds(1000));
	return 0;
}
