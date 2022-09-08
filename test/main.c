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
#include "../libGoReSym.h"

int main()
{
	char *file_name = "GoReSym";
	char *func_name = "runtime.casgstatus";
	char *itab_name = "go.itab.syscall.Errno,error";

	GoString fileName = {};
	GoString funcName = {};
	GoString itabName = {};
	struct FunctionAddress_return funcRet = {};
	GoUintptr itabRet = 0;

	fileName.p = file_name;
	fileName.n = strlen(file_name);

	funcName.p = func_name;
	funcName.n = strlen(func_name);

	itabName.p = itab_name;
	itabName.n = strlen(itab_name);

	funcRet = FunctionAddress(fileName, funcName);
	itabRet = ITabAddress(fileName, itabName);

	printf("func: addr=[%p], size=[%d]\n", funcRet.r0, funcRet.r1);
	printf("itab: addr=[%p]\n", itabRet);
	return 0;
}
