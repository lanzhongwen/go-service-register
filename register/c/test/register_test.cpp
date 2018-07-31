#include "gtest/gtest.h"
#include "register/c/register.h"
#include <pthread.h>

#include <chrono>
#include <iostream>
#include <memory>
#include <thread>

namespace {
class RegisterTest : public testing::Test {
	protected:
		RegisterTest() {
		}
		void SetUp() {
			std::cout << "SetUp()" << std::endl;
		}
		void TearDown() {
			std::cout << "TearDown()" << std::endl;
		}
};

TEST_F(RegisterTest, Affinity) {
	std::unique_ptr<std::thread> thread(new std::thread([=] {
	   std::string etcd_addrs_str("127.0.0.1:2379");
	   GoString etcd_addrs{etcd_addrs_str.c_str(), etcd_addrs_str.size()};
	   std::string service("/lzw/jd");
	   GoString serviceName{service.c_str(), service.size()};
	   std::string value("xxxx");
	   GoString serviceValue{value.c_str(), value.size()};
	   GoInt timeout = 10;
	   GoInt ttl = 10;
	   int ret = Register(etcd_addrs, serviceName, serviceValue, timeout, ttl);
	   if (!ret) {
	      std::cout << "Failed to Register" << std::endl;
	   }
	   std::this_thread::sleep_for(std::chrono::seconds(1000));
	}));
	cpu_set_t cpuset;
	CPU_ZERO(&cpuset);
	CPU_SET(1, &cpuset);
	int rc = pthread_setaffinity_np(thread->native_handle(), sizeof(cpu_set_t), &cpuset);
	if (rc != 0) {
	   std::cerr << "Error calling pthread_setaffinity_np: " << rc << "\n";
	}
	thread->join();
}
}
