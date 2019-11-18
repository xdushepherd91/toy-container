#define _GNU_SOURCE
#include <sys/types.h>
#include <sys/wait.h>
#include <stdio.h>
#include <sched.h>
#include <signal.h>
#include <unistd.h>
#include <stdlib.h>


#define STACK_SIZE (1024 * 1024)
static char child_stack[STACK_SIZE];

int flag = CLONE_NEWUTS;


int child_func(void *arg){
    int result =  sethostname("In Namespace", 12);
    printf("sethostname's result is %d",result);
    return 1;
}

void nsexec(void){

    /*
       xdushepherd 2019/11/18 9:15
       1. clone子进程，使得子进程可以进入新的命名空间
    */

    if(1){
    	return;
    }

    int child_pid = clone(child_func,child_stack+STACK_SIZE,CLONE_NEWUTS | SIGCHLD,NULL);

    printf("child pid is %d",child_pid);
    
    exit(0);
}
