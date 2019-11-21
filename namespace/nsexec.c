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


#include <stdint.h>
#include <linux/netlink.h>

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

struct config_t {
	char *data;

	uint32_t cloneflags;
	char * pidpath;
	char * container_root;
	char * container_id;
	char *namespaces;
};


/*
 * List of netlink message types sent to us as part of bootstrapping the init.
 * These constants are defined in libcontainer/message_linux.go.
 */
#define INIT_MSG			62000
#define CLONE_FLAGS_ATTR	27281
#define NS_PATHS_ATTR		27282
#define UIDMAP_ATTR			27283
#define GIDMAP_ATTR			27284
#define SETGROUP_ATTR		27285
#define OOM_SCORE_ADJ_ATTR	27286
#define ROOTLESS_EUID_ATTR	27287
#define UIDMAPPATH_ATTR	    27288
#define GIDMAPPATH_ATTR	    27289
#define GIDMAPPATH_ATTR	    27289
#define PID_PATH            27290
#define CONTAINER_PATH      27291
#define CONTAINER_ID        27292


#define PATH_MAX            260





void nl_parse(int fd ,struct config_t * config){
    size_t len,size;
    struct nlmsghdr hdr;
    char *data,*current;
    len = read(fd,&hdr,NLMSG_HDRLEN);

    if(len!=NLMSG_HDRLEN){
        perror("netlink header length does match error");
    }

	size = NLMSG_PAYLOAD(&hdr, 0);

	current = data = malloc(size);
	if (!data)
		perror("failed to allocate  bytes of memory for nl_payload");

	/* Parse the netlink payload. */
	config->data = data;
	while (current < data + size) {
		struct nlattr *nlattr = (struct nlattr *)current;
		size_t payload_len = nlattr->nla_len - NLA_HDRLEN;

		/* Advance to payload. */
		current += NLA_HDRLEN;

		/* Handle payload. */
		switch (nlattr->nla_type) {
		case CLONE_FLAGS_ATTR:
			config->cloneflags = *(uint32_t *) current;
			break;
		case PID_PATH:
			config->pidpath = current;
			break;
		case CONTAINER_PATH:
			config->container_root = current;
			break;
		case CONTAINER_ID:
			config->container_id = current;
			break;
		case NS_PATHS_ATTR:
			config->namespaces = current;
			break;
		default:
			perror("unknown netlink message type ");
		}

		current += NLA_ALIGN(payload_len);
	}



}

int child_func(void *arg) __attribute__ ((noinline));
int clone_parent(jmp_buf *env, int flags , int jmpval){
	struct clone_t ca = {
		.env = env,
		.jmpval = jmpval,
	};
        printf("clone parent %d\n ",jmpval);
	return clone(child_func, child_stack,  flags , &ca);
}

int child_func(void *arg){

       struct clone_t *ca = (struct clone_t *)arg;
       longjmp(*ca->env, ca->jmpval); 
}

// 从环境变量解析pipe的fd
int initpipe(void){
	int pipenum;
	char *initpipe, *endptr;

	initpipe = getenv("_LIBCONTAINER_INITPIPE");
	if (initpipe == NULL || *initpipe == '\0')
		return -1;

	pipenum = strtol(initpipe, &endptr, 10);
	if (*endptr != '\0')
		perror("unable to parse _LIBCONTAINER_INITPIPE");

	return pipenum;
}

int parseExec(){
    int exec;
	char *exec_str, *endptr;
    exec_str = getenv("_TOYCONTAINER_EXEC");
	if (exec_str == NULL || *exec_str == '\0')
		return -1;
	exec = strtol(exec_str,&endptr,10);
    return exec;
}


/* Returns the clone(2) flag for a namespace, given the name of a namespace. */
int nsflag(char *name)
{
	if (!strcmp(name, "cgroup"))
		return CLONE_NEWCGROUP;
	else if (!strcmp(name, "ipc"))
		return CLONE_NEWIPC;
	else if (!strcmp(name, "mnt"))
		return CLONE_NEWNS;
	else if (!strcmp(name, "net"))
		return CLONE_NEWNET;
	else if (!strcmp(name, "pid"))
		return CLONE_NEWPID;
	else if (!strcmp(name, "user"))
		return CLONE_NEWUSER;
	else if (!strcmp(name, "uts"))
		return CLONE_NEWUTS;

	/* If we don't recognise a name, fallback to 0. */
	return 0;
}


void join_namespaces(char *nslist)
{
	int num = 0, i;
	char *saveptr = NULL;
	char *namespace = strtok_r(nslist, ",", &saveptr);
	struct namespace_t {
		int fd;
		int ns;
		char type[PATH_MAX];
		char path[PATH_MAX];
	} *namespaces = NULL;

	if (!namespace || !strlen(namespace) || !strlen(nslist))
		perror("ns paths are empty");

	/*
	 * We have to open the file descriptors first, since after
	 * we join the mnt namespace we might no longer be able to
	 * access the paths.
	 */
	do {
		int fd;
		char *path;
		struct namespace_t *ns;

		/* Resize the namespace array. */
		namespaces = realloc(namespaces, ++num * sizeof(struct namespace_t));
		if (!namespaces)
			perror("failed to reallocate namespace array");
		ns = &namespaces[num - 1];

		/* Split 'ns:path'. */
		path = strstr(namespace, ":");
		if (!path)
			perror("failed to parse ");
		*path++ = '\0';

		fd =  open(path, O_RDONLY);
		if (fd < 0)
			perror("failed to open ");

		ns->fd = fd;
		ns->ns = nsflag(namespace);
		strncpy(ns->path, path, PATH_MAX - 1);
		ns->path[PATH_MAX - 1] = '\0';
	} while ((namespace = strtok_r(NULL, ",", &saveptr)) != NULL);

	/*
	 * The ordering in which we join namespaces is important. We should
	 * always join the user namespace *first*. This is all guaranteed
	 * from the container_linux.go side of this, so we're just going to
	 * follow the order given to us.
	 */
    /*
        xdushepherd 2019/11/15 10:47
        按照标准的配置文件，path为空，这里不操作
    */
	for (i = 0; i < num; i++) {
		struct namespace_t ns = namespaces[i];

		if (setns(ns.fd, ns.ns) < 0)
			perror("failed to setns to ");

		close(ns.fd);
	}

	free(namespaces);
}


void nsexec(void){


	/*
	 * If we don't have an init pipe, just return to the go routine.
	 * We'll only get an init pipe for start or exec.
	 * xdushepherd 2019/11/14 11:59
	 * 用于从环境变量中获取消息通信的pipenum
	 */
	int pipenum = initpipe();
	if (pipenum == -1){
		return;
    }
    /*
        xdushepherd 2019/11/21 14:40
        用于区分是否是是新建namespace还是setns
    */
    int exec = parseExec();


    jmp_buf env;
    struct config_t config = {0};
    nl_parse(pipenum,&config);

    char * container_root = config.container_root;

    if(container_root==NULL){
        perror("容器根目录为空");
    }



    /*
       xdushepherd 2019/11/18 9:15
       1. clone子进程，使得子进程可以进入新的命名空间
    */
    switch (setjmp(env)) {
        case 0:{
            if(exec){
                join_namespaces(config.namespaces);
            }else{
                int child_pid = clone_parent(&env,config.cloneflags,1);
                FILE* fp = fopen(config.pidpath, "w+");
                fprintf(fp, "%d",child_pid);
                fclose(fp);
                waitpid(child_pid, NULL, 0);
                exit(0);
            }
        }
        case 1:
            sethostname(config.container_id,10);
            if(!exec){
                if (mount("proc",strcat(container_root,"/proc"),  "proc", 0, NULL) !=0 ) {
                     perror("proc");
                }
                if (mount("sysfs", strcat(container_root,"/sys"), "sysfs", 0, NULL)!=0) {
                      perror("sys");
                }
                if (mount("none", strcat(container_root,"/tmp"), "tmpfs", 0, NULL)!=0) {
                    perror("tmp");
                }
                if (mount("udev", strcat(container_root,"/dev"), "devtmpfs", 0, NULL)!=0) {
                   perror("dev");
                }
                if (mount("devpts", strcat(container_root,"/dev/pts"), "devpts", 0, NULL)!=0) {
                    perror("dev/pts");
                }
                if (mount("shm", strcat(container_root,"/dev/shm"),"tmpfs", 0, NULL)!=0) {
                    perror("dev/shm");
                }
                if (mount("tmpfs", strcat(container_root,"/run"),"tmpfs", 0, NULL)!=0) {
                    perror("run");
                }
            }

            /* chroot 隔离目录 */
            if ( chdir(container_root) != 0 || chroot("./") != 0 ){
               perror("error  when  chdir/chroot");
            }

            setenv("PS1", "[\\u@\\H \\W] #", 0);
            execv(child_args[0], child_args);

            free(config.data);
            return;
        }
    }
}
