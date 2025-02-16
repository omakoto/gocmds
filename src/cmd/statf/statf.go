package main

import (
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/fileinput"
	"strings"
)

func main() {
	common.RunAndExit(RealMain)
}

func RealMain() int {
	for line := range fileinput.FileInput() {
		c := strings.Split(strings.TrimSpace(line), " ")
		fmt.Printf(
			"pid=%s"+
				" comm=%s"+
				" state=%s"+
				" ppid=%s"+
				" pgrp=%s"+
				" session=%s"+
				" tty=%s"+
				" tpgid=%s"+
				" flags=%s"+
				" minflt=%s"+
				" cminflt=%s"+
				" majflt=%s"+
				" cmajflt=%s"+
				" utime=%s"+
				" stime=%s"+
				" cutime=%s"+
				" cstime=%s"+
				" priority=%s"+
				" nice=%s"+
				" num_threads=%s"+
				" itrealvalue=%s"+
				" starttime=%s"+
				" vsize=%s"+
				" rss=%s"+
				" rsslim=%s"+
				" startcode=%s"+
				" endcode=%s"+
				" startstack=%s"+
				" kstkesp=%s"+
				" kstkeip=%s"+
				" signal=%s"+
				" blocked=%s"+
				" sigignore=%s"+
				" sigcatch=%s"+
				" wchan=%s"+
				" nswap=%s"+
				" cnswap=%s"+
				" exit_signal=%s"+
				" processor=%s"+
				" rt_priority=%s"+
				" policy=%s"+
				" delayacct_blkio_ticks=%s"+
				" guest_time=%s"+
				" cguest_time=%s"+
				" start_data=%s"+
				" end_data=%s"+
				" start_brk=%s"+
				" arg_start=%s"+
				" arg_end=%s"+
				" env_start=%s"+
				" env_end=%s"+
				" exit_code=%s\n",
			c[0],
			c[1],
			c[2],
			c[3],
			c[4],
			c[5],
			c[6],
			c[7],
			c[8],
			c[9],
			c[10],
			c[11],
			c[12],
			c[13],
			c[14],
			c[15],
			c[16],
			c[17],
			c[18],
			c[19],
			c[20],
			c[21],
			c[22],
			c[23],
			c[24],
			c[25],
			c[26],
			c[27],
			c[28],
			c[29],
			c[30],
			c[31],
			c[32],
			c[33],
			c[34],
			c[35],
			c[36],
			c[37],
			c[38],
			c[39],
			c[40],
			c[41],
			c[42],
			c[43],
			c[44],
			c[45],
			c[46],
			c[47],
			c[48],
			c[49],
			c[50],
			c[51],
		)
	}
	return 0
}
