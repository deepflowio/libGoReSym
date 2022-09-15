// gcc main.c --static -L../ -lGoReSym  -lpthread
#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
#include <assert.h>
#include <stddef.h>
#include <string.h>
#include <errno.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/resource.h>
#include "../libGoReSym.h"

int main()
{
	char *file_name = "GoELF";
	char *func_name = "runtime.casgstatus";
	char *itab_name = "go.itab.*net.TCPConn,net.Conn";

	struct function_address_return func = {};
	GoUintptr itab = 0;

	func = function_address(file_name, func_name);
	itab = itab_address(file_name, itab_name);

	printf("func: addr=[%p], size=[%d]\n", func.r0, func.r1);
	printf("itab: addr=[%p]\n", itab);
	return 0;
}
