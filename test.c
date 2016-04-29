#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "libaosp_su_daemon.h"

int main() {
	// Test exec
	char *cmdstr = "ls";
	GoString cmd;
       	cmd.p = cmdstr;
	cmd.n = strlen(cmdstr);

	char *argsstr = "-l -i -s -a";
	GoString args;
	args.p = argsstr;
	args.n = strlen(argsstr);

	Execv1(cmd, args, 1);
}
