#define _GNU_SOURCE
#include <sys/types.h>
#include <sys/wait.h>
#include <stdio.h>
#include <sched.h>
#include <signal.h>
#include <unistd.h>

#define STACK_SIZE (1024 * 1024)
static char child_stack[STACK_SIZE];

int flag = CLONE_NEWIPC | CLONE_NEWNET | CLONE_NEWNS | CLONE_NEWPID | CLONE_NEWUTS |CLONE_NEWCGROUP;


void child_func(void){
    sethostname("In Namespace", 12);
    return;
}

void nsexec(void){

    /*
       xdushepherd 2019/11/18 9:15
       1. clone子进程，使得子进程可以进入新的命名空间
    */


    int child_pid = clone(child_func,child_stack+STACK_SIZE,flag);



    exit(0);
}
