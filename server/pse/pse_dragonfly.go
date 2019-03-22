// Copyright 2015-2016 Apcera Inc. All rights reserved.

package pse

/*
#include <sys/types.h>
#include <sys/sysctl.h>
#include <sys/user.h>
//#include <sys/resource.h> // for fscale
#include <stddef.h>
#include <unistd.h>

long pagetok(long size)
{
    int pageshift, pagesize;

    pagesize = getpagesize();
    pageshift = 0;

    while (pagesize > 1) {
        pageshift++;
        pagesize >>= 1;
    }

    return (size << pageshift);
}

int getusage(double *pcpu, unsigned int *rss, unsigned int *vss)
{
    int mib[4], ret, fscale;
    size_t len, oldlen;
    struct kinfo_proc kp;

    len = 4;
    sysctlnametomib("kern.proc.pid", mib, &len);

    mib[3] = getpid();
    len = sizeof(kp);

    ret = sysctl(mib, 4, &kp, &len, NULL, 0);
    if (ret != 0) {
        return (errno);
    }

    oldlen = sizeof(fscale);
    if (sysctlbyname("kern.fscale", &fscale, &oldlen, NULL, 0) < 0) {
        return (errno);
    }
    
#define	fxtofl(fixpt)	((double)(fixpt) / fscale)
#define KI_PROC(ki, f) (ki.kp_ ## f)
#define KI_LWP(ki, f) (ki.kp_lwp.kl_ ## f)

//    if (KI_PROC(kp, swtime) == 0 || (KI_PROC(kp, flags) & P_SWAPPEDOUT)) {
    if (kp.kp_swtime == 0 || (kp.kp_flags & P_SWAPPEDOUT)) {
        *pcpu = 0.0;
    } else {
        *pcpu = (100.0 * fxtofl(kp.kp_lwp.kl_pctcpu));
        //*pcpu = (100.0 * fxtofl(KI_LWP(kp, pctcpu)));
    }

    *rss = pagetok(kp.kp_vm_rssize);
    *vss = kp.kp_vm_map_size;
    //*pcpu = fscale;//kp.kp_lwp.kl_pctcpu / 100.0; //kp.ki_pctcpu;

    return 0;
}

*/
import "C"

import (
	"syscall"
)

// CPU and memory statistics
func ProcUsage(pcpu *float64, rss, vss *int64) error {
	var r, v C.uint
	var c C.double

	if ret := C.getusage(&c, &r, &v); ret != 0 {
		return syscall.Errno(ret)
	}

	*pcpu = float64(c)
	*rss = int64(r)
	*vss = int64(v)

	return nil
}
