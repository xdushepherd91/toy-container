##include

int flag = CLONE_NEWIPC | CLONE_NEWNET | CLONE_NEWNS | CLONE_NEWPID | CLONE_NEWUTS |CLONE_NEWCGROUP;

/*
    xdushepherd 2019/11/18 10:04
    复制runc中的代码，用来作为clone函数的第二参数
 */
struct clone_t {

	char stack[4096] __attribute__ ((aligned(16)));
	char stack_ptr[0];

};

void child_func(void){
    sethostname("In Namespace", 12);
    return;
}

void nsexec(void){

    /*
       xdushepherd 2019/11/18 9:15
       1. clone子进程，使得子进程可以进入新的命名空间
    */

    struct clone_t cstack;

    int child_pid = clone(child_func,cstack.stack_ptr,flag);



    exit(0);
}
