package xproc

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <libproc.h>

int getFilenameByPid(int id) {
	pid_t pid = id;
	int ret;
	char pathbuf[PROC_PIDPATHINFO_MAXSIZE];
        ret = proc_pidpath (pid, pathbuf, sizeof(pathbuf));
        if ( ret <= 0 ) {
            fprintf(stderr, "PID %d: proc_pidpath ();\n", pid);
            fprintf(stderr, "    %s\n", strerror(errno));
        } else {
            printf("proc %d: %s\n", pid, pathbuf);
        }
        return ret;
}
*/
import "C"

func getFilenameByPid(pid ProcId) int {
	return 0
}
