package nsenter


/*
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

__attribute__((constructor)) void enter_namespace(void) {
    char *dockerEngine_pid;
    dockerEngine_pid = getenv("dockerEngine_pid");
    if (dockerEngine_pid) {
        //fprintf(stdout, "got dockerEngine_pid=%s\n", dockerEngine_pid);
    } else {
        //fprintf(stdout, "missing dockerEngine_pid env skip nsenter");
        return;
    }

    char *dockerEngine_cmd;
    dockerEngine_cmd = getenv("dockerEngine_cmd");
    if (dockerEngine_cmd) {
        //fprintf(stdout, "got dockerEngine_cmd=%s\n", dockerEngine_cmd);
    } else {
        //fprintf(stdout, "missing dockerEngine_cmd env skip nsenter");
        return;
    }

    int i;
    char nspath[1024];
    char *namespaces[] = {"ipc", "uts", "net", "pid", "mnt"};

    for (i=0; i<5; i++) {
        sprintf(nspath, "/proc/%s/ns/%s", dockerEngine_pid, namespaces[i]);
        int fd = open(nspath, O_RDONLY);
        if (setns(fd, 0) == -1) {
            //fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], stderr(errno));
        } else {
            //fprintf(stdout, "setns on %s namespace succeeded\n", namespace[i]);
        }

        close(fd);
    }

    int res = system(dockerEngine_cmd);
    exit(0)
    return;
}
*/
import "C"
