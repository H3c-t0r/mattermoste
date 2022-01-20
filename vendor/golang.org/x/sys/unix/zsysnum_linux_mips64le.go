// go run linux/mksysnum.go -Wall -Werror -static -I/tmp/include /tmp/include/asm/unistd.h
// Code generated by the command above; see README.md. DO NOT EDIT.

//go:build mips64le && linux
// +build mips64le,linux

package unix

const (
	SYS_READ                    = 5000
	SYS_WRITE                   = 5001
	SYS_OPEN                    = 5002
	SYS_CLOSE                   = 5003
	SYS_STAT                    = 5004
	SYS_FSTAT                   = 5005
	SYS_LSTAT                   = 5006
	SYS_POLL                    = 5007
	SYS_LSEEK                   = 5008
	SYS_MMAP                    = 5009
	SYS_MPROTECT                = 5010
	SYS_MUNMAP                  = 5011
	SYS_BRK                     = 5012
	SYS_RT_SIGACTION            = 5013
	SYS_RT_SIGPROCMASK          = 5014
	SYS_IOCTL                   = 5015
	SYS_PREAD64                 = 5016
	SYS_PWRITE64                = 5017
	SYS_READV                   = 5018
	SYS_WRITEV                  = 5019
	SYS_ACCESS                  = 5020
	SYS_PIPE                    = 5021
	SYS__NEWSELECT              = 5022
	SYS_SCHED_YIELD             = 5023
	SYS_MREMAP                  = 5024
	SYS_MSYNC                   = 5025
	SYS_MINCORE                 = 5026
	SYS_MADVISE                 = 5027
	SYS_SHMGET                  = 5028
	SYS_SHMAT                   = 5029
	SYS_SHMCTL                  = 5030
	SYS_DUP                     = 5031
	SYS_DUP2                    = 5032
	SYS_PAUSE                   = 5033
	SYS_NANOSLEEP               = 5034
	SYS_GETITIMER               = 5035
	SYS_SETITIMER               = 5036
	SYS_ALARM                   = 5037
	SYS_GETPID                  = 5038
	SYS_SENDFILE                = 5039
	SYS_SOCKET                  = 5040
	SYS_CONNECT                 = 5041
	SYS_ACCEPT                  = 5042
	SYS_SENDTO                  = 5043
	SYS_RECVFROM                = 5044
	SYS_SENDMSG                 = 5045
	SYS_RECVMSG                 = 5046
	SYS_SHUTDOWN                = 5047
	SYS_BIND                    = 5048
	SYS_LISTEN                  = 5049
	SYS_GETSOCKNAME             = 5050
	SYS_GETPEERNAME             = 5051
	SYS_SOCKETPAIR              = 5052
	SYS_SETSOCKOPT              = 5053
	SYS_GETSOCKOPT              = 5054
	SYS_CLONE                   = 5055
	SYS_FORK                    = 5056
	SYS_EXECVE                  = 5057
	SYS_EXIT                    = 5058
	SYS_WAIT4                   = 5059
	SYS_KILL                    = 5060
	SYS_UNAME                   = 5061
	SYS_SEMGET                  = 5062
	SYS_SEMOP                   = 5063
	SYS_SEMCTL                  = 5064
	SYS_SHMDT                   = 5065
	SYS_MSGGET                  = 5066
	SYS_MSGSND                  = 5067
	SYS_MSGRCV                  = 5068
	SYS_MSGCTL                  = 5069
	SYS_FCNTL                   = 5070
	SYS_FLOCK                   = 5071
	SYS_FSYNC                   = 5072
	SYS_FDATASYNC               = 5073
	SYS_TRUNCATE                = 5074
	SYS_FTRUNCATE               = 5075
	SYS_GETDENTS                = 5076
	SYS_GETCWD                  = 5077
	SYS_CHDIR                   = 5078
	SYS_FCHDIR                  = 5079
	SYS_RENAME                  = 5080
	SYS_MKDIR                   = 5081
	SYS_RMDIR                   = 5082
	SYS_CREAT                   = 5083
	SYS_LINK                    = 5084
	SYS_UNLINK                  = 5085
	SYS_SYMLINK                 = 5086
	SYS_READLINK                = 5087
	SYS_CHMOD                   = 5088
	SYS_FCHMOD                  = 5089
	SYS_CHOWN                   = 5090
	SYS_FCHOWN                  = 5091
	SYS_LCHOWN                  = 5092
	SYS_UMASK                   = 5093
	SYS_GETTIMEOFDAY            = 5094
	SYS_GETRLIMIT               = 5095
	SYS_GETRUSAGE               = 5096
	SYS_SYSINFO                 = 5097
	SYS_TIMES                   = 5098
	SYS_PTRACE                  = 5099
	SYS_GETUID                  = 5100
	SYS_SYSLOG                  = 5101
	SYS_GETGID                  = 5102
	SYS_SETUID                  = 5103
	SYS_SETGID                  = 5104
	SYS_GETEUID                 = 5105
	SYS_GETEGID                 = 5106
	SYS_SETPGID                 = 5107
	SYS_GETPPID                 = 5108
	SYS_GETPGRP                 = 5109
	SYS_SETSID                  = 5110
	SYS_SETREUID                = 5111
	SYS_SETREGID                = 5112
	SYS_GETGROUPS               = 5113
	SYS_SETGROUPS               = 5114
	SYS_SETRESUID               = 5115
	SYS_GETRESUID               = 5116
	SYS_SETRESGID               = 5117
	SYS_GETRESGID               = 5118
	SYS_GETPGID                 = 5119
	SYS_SETFSUID                = 5120
	SYS_SETFSGID                = 5121
	SYS_GETSID                  = 5122
	SYS_CAPGET                  = 5123
	SYS_CAPSET                  = 5124
	SYS_RT_SIGPENDING           = 5125
	SYS_RT_SIGTIMEDWAIT         = 5126
	SYS_RT_SIGQUEUEINFO         = 5127
	SYS_RT_SIGSUSPEND           = 5128
	SYS_SIGALTSTACK             = 5129
	SYS_UTIME                   = 5130
	SYS_MKNOD                   = 5131
	SYS_PERSONALITY             = 5132
	SYS_USTAT                   = 5133
	SYS_STATFS                  = 5134
	SYS_FSTATFS                 = 5135
	SYS_SYSFS                   = 5136
	SYS_GETPRIORITY             = 5137
	SYS_SETPRIORITY             = 5138
	SYS_SCHED_SETPARAM          = 5139
	SYS_SCHED_GETPARAM          = 5140
	SYS_SCHED_SETSCHEDULER      = 5141
	SYS_SCHED_GETSCHEDULER      = 5142
	SYS_SCHED_GET_PRIORITY_MAX  = 5143
	SYS_SCHED_GET_PRIORITY_MIN  = 5144
	SYS_SCHED_RR_GET_INTERVAL   = 5145
	SYS_MLOCK                   = 5146
	SYS_MUNLOCK                 = 5147
	SYS_MLOCKALL                = 5148
	SYS_MUNLOCKALL              = 5149
	SYS_VHANGUP                 = 5150
	SYS_PIVOT_ROOT              = 5151
	SYS__SYSCTL                 = 5152
	SYS_PRCTL                   = 5153
	SYS_ADJTIMEX                = 5154
	SYS_SETRLIMIT               = 5155
	SYS_CHROOT                  = 5156
	SYS_SYNC                    = 5157
	SYS_ACCT                    = 5158
	SYS_SETTIMEOFDAY            = 5159
	SYS_MOUNT                   = 5160
	SYS_UMOUNT2                 = 5161
	SYS_SWAPON                  = 5162
	SYS_SWAPOFF                 = 5163
	SYS_REBOOT                  = 5164
	SYS_SETHOSTNAME             = 5165
	SYS_SETDOMAINNAME           = 5166
	SYS_CREATE_MODULE           = 5167
	SYS_INIT_MODULE             = 5168
	SYS_DELETE_MODULE           = 5169
	SYS_GET_KERNEL_SYMS         = 5170
	SYS_QUERY_MODULE            = 5171
	SYS_QUOTACTL                = 5172
	SYS_NFSSERVCTL              = 5173
	SYS_GETPMSG                 = 5174
	SYS_PUTPMSG                 = 5175
	SYS_AFS_SYSCALL             = 5176
	SYS_RESERVED177             = 5177
	SYS_GETTID                  = 5178
	SYS_READAHEAD               = 5179
	SYS_SETXATTR                = 5180
	SYS_LSETXATTR               = 5181
	SYS_FSETXATTR               = 5182
	SYS_GETXATTR                = 5183
	SYS_LGETXATTR               = 5184
	SYS_FGETXATTR               = 5185
	SYS_LISTXATTR               = 5186
	SYS_LLISTXATTR              = 5187
	SYS_FLISTXATTR              = 5188
	SYS_REMOVEXATTR             = 5189
	SYS_LREMOVEXATTR            = 5190
	SYS_FREMOVEXATTR            = 5191
	SYS_TKILL                   = 5192
	SYS_RESERVED193             = 5193
	SYS_FUTEX                   = 5194
	SYS_SCHED_SETAFFINITY       = 5195
	SYS_SCHED_GETAFFINITY       = 5196
	SYS_CACHEFLUSH              = 5197
	SYS_CACHECTL                = 5198
	SYS_SYSMIPS                 = 5199
	SYS_IO_SETUP                = 5200
	SYS_IO_DESTROY              = 5201
	SYS_IO_GETEVENTS            = 5202
	SYS_IO_SUBMIT               = 5203
	SYS_IO_CANCEL               = 5204
	SYS_EXIT_GROUP              = 5205
	SYS_LOOKUP_DCOOKIE          = 5206
	SYS_EPOLL_CREATE            = 5207
	SYS_EPOLL_CTL               = 5208
	SYS_EPOLL_WAIT              = 5209
	SYS_REMAP_FILE_PAGES        = 5210
	SYS_RT_SIGRETURN            = 5211
	SYS_SET_TID_ADDRESS         = 5212
	SYS_RESTART_SYSCALL         = 5213
	SYS_SEMTIMEDOP              = 5214
	SYS_FADVISE64               = 5215
	SYS_TIMER_CREATE            = 5216
	SYS_TIMER_SETTIME           = 5217
	SYS_TIMER_GETTIME           = 5218
	SYS_TIMER_GETOVERRUN        = 5219
	SYS_TIMER_DELETE            = 5220
	SYS_CLOCK_SETTIME           = 5221
	SYS_CLOCK_GETTIME           = 5222
	SYS_CLOCK_GETRES            = 5223
	SYS_CLOCK_NANOSLEEP         = 5224
	SYS_TGKILL                  = 5225
	SYS_UTIMES                  = 5226
	SYS_MBIND                   = 5227
	SYS_GET_MEMPOLICY           = 5228
	SYS_SET_MEMPOLICY           = 5229
	SYS_MQ_OPEN                 = 5230
	SYS_MQ_UNLINK               = 5231
	SYS_MQ_TIMEDSEND            = 5232
	SYS_MQ_TIMEDRECEIVE         = 5233
	SYS_MQ_NOTIFY               = 5234
	SYS_MQ_GETSETATTR           = 5235
	SYS_VSERVER                 = 5236
	SYS_WAITID                  = 5237
	SYS_ADD_KEY                 = 5239
	SYS_REQUEST_KEY             = 5240
	SYS_KEYCTL                  = 5241
	SYS_SET_THREAD_AREA         = 5242
	SYS_INOTIFY_INIT            = 5243
	SYS_INOTIFY_ADD_WATCH       = 5244
	SYS_INOTIFY_RM_WATCH        = 5245
	SYS_MIGRATE_PAGES           = 5246
	SYS_OPENAT                  = 5247
	SYS_MKDIRAT                 = 5248
	SYS_MKNODAT                 = 5249
	SYS_FCHOWNAT                = 5250
	SYS_FUTIMESAT               = 5251
	SYS_NEWFSTATAT              = 5252
	SYS_UNLINKAT                = 5253
	SYS_RENAMEAT                = 5254
	SYS_LINKAT                  = 5255
	SYS_SYMLINKAT               = 5256
	SYS_READLINKAT              = 5257
	SYS_FCHMODAT                = 5258
	SYS_FACCESSAT               = 5259
	SYS_PSELECT6                = 5260
	SYS_PPOLL                   = 5261
	SYS_UNSHARE                 = 5262
	SYS_SPLICE                  = 5263
	SYS_SYNC_FILE_RANGE         = 5264
	SYS_TEE                     = 5265
	SYS_VMSPLICE                = 5266
	SYS_MOVE_PAGES              = 5267
	SYS_SET_ROBUST_LIST         = 5268
	SYS_GET_ROBUST_LIST         = 5269
	SYS_KEXEC_LOAD              = 5270
	SYS_GETCPU                  = 5271
	SYS_EPOLL_PWAIT             = 5272
	SYS_IOPRIO_SET              = 5273
	SYS_IOPRIO_GET              = 5274
	SYS_UTIMENSAT               = 5275
	SYS_SIGNALFD                = 5276
	SYS_TIMERFD                 = 5277
	SYS_EVENTFD                 = 5278
	SYS_FALLOCATE               = 5279
	SYS_TIMERFD_CREATE          = 5280
	SYS_TIMERFD_GETTIME         = 5281
	SYS_TIMERFD_SETTIME         = 5282
	SYS_SIGNALFD4               = 5283
	SYS_EVENTFD2                = 5284
	SYS_EPOLL_CREATE1           = 5285
	SYS_DUP3                    = 5286
	SYS_PIPE2                   = 5287
	SYS_INOTIFY_INIT1           = 5288
	SYS_PREADV                  = 5289
	SYS_PWRITEV                 = 5290
	SYS_RT_TGSIGQUEUEINFO       = 5291
	SYS_PERF_EVENT_OPEN         = 5292
	SYS_ACCEPT4                 = 5293
	SYS_RECVMMSG                = 5294
	SYS_FANOTIFY_INIT           = 5295
	SYS_FANOTIFY_MARK           = 5296
	SYS_PRLIMIT64               = 5297
	SYS_NAME_TO_HANDLE_AT       = 5298
	SYS_OPEN_BY_HANDLE_AT       = 5299
	SYS_CLOCK_ADJTIME           = 5300
	SYS_SYNCFS                  = 5301
	SYS_SENDMMSG                = 5302
	SYS_SETNS                   = 5303
	SYS_PROCESS_VM_READV        = 5304
	SYS_PROCESS_VM_WRITEV       = 5305
	SYS_KCMP                    = 5306
	SYS_FINIT_MODULE            = 5307
	SYS_GETDENTS64              = 5308
	SYS_SCHED_SETATTR           = 5309
	SYS_SCHED_GETATTR           = 5310
	SYS_RENAMEAT2               = 5311
	SYS_SECCOMP                 = 5312
	SYS_GETRANDOM               = 5313
	SYS_MEMFD_CREATE            = 5314
	SYS_BPF                     = 5315
	SYS_EXECVEAT                = 5316
	SYS_USERFAULTFD             = 5317
	SYS_MEMBARRIER              = 5318
	SYS_MLOCK2                  = 5319
	SYS_COPY_FILE_RANGE         = 5320
	SYS_PREADV2                 = 5321
	SYS_PWRITEV2                = 5322
	SYS_PKEY_MPROTECT           = 5323
	SYS_PKEY_ALLOC              = 5324
	SYS_PKEY_FREE               = 5325
	SYS_STATX                   = 5326
	SYS_RSEQ                    = 5327
	SYS_IO_PGETEVENTS           = 5328
	SYS_PIDFD_SEND_SIGNAL       = 5424
	SYS_IO_URING_SETUP          = 5425
	SYS_IO_URING_ENTER          = 5426
	SYS_IO_URING_REGISTER       = 5427
	SYS_OPEN_TREE               = 5428
	SYS_MOVE_MOUNT              = 5429
	SYS_FSOPEN                  = 5430
	SYS_FSCONFIG                = 5431
	SYS_FSMOUNT                 = 5432
	SYS_FSPICK                  = 5433
	SYS_PIDFD_OPEN              = 5434
	SYS_CLONE3                  = 5435
	SYS_CLOSE_RANGE             = 5436
	SYS_OPENAT2                 = 5437
	SYS_PIDFD_GETFD             = 5438
	SYS_FACCESSAT2              = 5439
	SYS_PROCESS_MADVISE         = 5440
	SYS_EPOLL_PWAIT2            = 5441
	SYS_MOUNT_SETATTR           = 5442
	SYS_QUOTACTL_FD             = 5443
	SYS_LANDLOCK_CREATE_RULESET = 5444
	SYS_LANDLOCK_ADD_RULE       = 5445
	SYS_LANDLOCK_RESTRICT_SELF  = 5446
	SYS_PROCESS_MRELEASE        = 5448
)
