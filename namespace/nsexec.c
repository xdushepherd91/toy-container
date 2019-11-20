#define _GNU_SOURCE
#include <sys/types.h>
#include <sys/wait.h>
#include <stdio.h>
#include <sched.h>
#include <signal.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <setjmp.h>
#include <sys/mount.h>

#define STACK_SIZE (1024 * 1024)
static char child_stack[STACK_SIZE];
char* const child_args[] = {
   "/bin/sh",
   "/home/test.sh",
   NULL 
};

int flag = CLONE_NEWUTS;

struct clone_t {
	jmp_buf *env;
	int jmpval;
};

int child_func(void *arg) __attribute__ ((noinline));
int clone_parent(jmp_buf *env, int jmpval)
{
	struct clone_t ca = {
		.env = env,
		.jmpval = jmpval,
	};
        printf("clone parent %d\n ",jmpval);
	return clone(child_func, child_stack,  CLONE_NEWUTS | CLONE_NEWIPC | CLONE_NEWPID | CLONE_NEWNS | SIGCHLD, &ca);
}

int child_func(void *arg){

       struct clone_t *ca = (struct clone_t *)arg;
       longjmp(*ca->env, ca->jmpval); 
}

void nsexec(void){

    jmp_buf env;

    FILE *fp = NULL;
    fp = fopen("/opt/toy-container/test.txt", "w+");
    fprintf(fp, "This is testing for fprintf...%s  \n",getenv("_LIBCONTAINER_INITPIPE"));


    char* str_ptr = getenv("_LIBCONTAINER_INITPIPE");


    if(str_ptr == NULL){
    	return;
    }

    /*
       xdushepherd 2019/11/18 9:15
       1. clone子进程，使得子进程可以进入新的命名空间
    */


    switch (setjmp(env)) {

	   case 0:{
                   
	           int child_pid = clone_parent(&env,1);

                   printf("child pid is %d\n",child_pid);

                   fp = fopen("/opt/toy-container/child_pid.txt", "w+");
                   fprintf(fp, "%d",child_pid);
       
       		   //printf("pid data persistence");		   
		   fclose(fp);
		   waitpid(child_pid, NULL, 0);
                   exit(0);

		  }
           case 1: {
                  int result = sethostname("container",12);
		  printf("Container [%d] - inside the container!\n", getpid());
		  printf("\ninside container now %d \n",result);
		      //remount "/proc" to make sure the "top" and "ps" show container's information
                  if (mount("proc", "/opt/toy-container/default-id/target/proc", "proc", 0, NULL) !=0 ) {
                         perror("proc");
                  }
                  if (mount("sysfs", "/opt/toy-container/default-id/target/sys", "sysfs", 0, NULL)!=0) {
                          perror("sys");
                  }
                  if (mount("none", "/opt/toy-container/default-id/target/tmp", "tmpfs", 0, NULL)!=0) {
                        perror("tmp");
                  }
                  if (mount("udev", "/opt/toy-container/rootfs/dev", "devtmpfs", 0, NULL)!=0) {
                       perror("dev");
                  }
                 if (mount("devpts", "/opt/toy-container/default-id/target/dev/pts", "devpts", 0, NULL)!=0) {
                        perror("dev/pts");
                 }
                 if (mount("shm", "/opt/toy-container/default-id/target/dev/shm", "tmpfs", 0, NULL)!=0) {
                        perror("dev/shm");
                 }
                 if (mount("tmpfs", "/opt/toy-container/default-id/target/run", "tmpfs", 0, NULL)!=0) {
                        perror("run");
                 }
		
		 printf("before chroot\n");
		 /* chroot 隔离目录 */
                 if ( chdir("/opt/toy-container/default-id/target") != 0 || chroot("./") != 0 ){
                       perror("chdir/chroot");
                 }		  

                 printf("\nok\n");
                 setenv("PS1", "[\\u@\\H \\W] #", 0);
		 execv(child_args[0], child_args);	
		 
//		 return;
		 exit(0);
	 }

    }
}
