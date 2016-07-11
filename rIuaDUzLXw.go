/*
 * Copyright (c) 2013 Landon Fuller &lt;landonf@mac68k.info&gt;
 * All rights reserved.
 */

/* Select-based polling support for the pcap API */
package pcap

/*
#include &lt;pcap/pcap.h&gt;

#include &lt;sys/select.h&gt;

#include &lt;stdlib.h&gt;
#include &lt;unistd.h&gt;
#include &lt;string.h&gt;
#include &lt;stdint.h&gt;

// Indirection required to use Go callbacks
extern void mac68k_pcap_dispatchCallback (unsigned char *user, struct pcap_pkthdr *h, unsigned char *bytes);
static void pcap_dispatch_cb_handler (u_char *user, const struct pcap_pkthdr *h, const u_char *bytes) {
    mac68k_pcap_dispatchCallback((unsigned char *) user, (struct pcap_pkthdr *) h, (unsigned char *) bytes);
}

// cgo gets upset when we use &#39;select&#39;
static int my_select (int nfds, fd_set *readfds, fd_set *writefds, fd_set *errorfds, struct timeval *timeout) {
    return select(nfds, readfds, writefds, errorfds, timeout);
}

// cgo refuses to resolve the FD_* macros.
static void MY_FD_ZERO (fd_set *fdset) {
    FD_ZERO(fdset);
}
static void MY_FD_SET (int fd, fd_set *fdset) {
    FD_SET(fd, fdset);
}
static int MY_FD_ISSET (int fd, fd_set *fdset) {
    return FD_ISSET(fd, fdset);
}
*/
import &#34;C&#34;
import (
	&#34;errors&#34;
	&#34;fmt&#34;
	&#34;unsafe&#34;
)

// Server that manages polling the backing file descriptor
type pollServer struct {
	// Channel that may be used to read packets
	packets chan []byte

	// The pcap file descriptor
	pcapfd C.int

	// If readable, the poll server should stop. Can be made readable
	// by writing to the signalFD
	waitfd C.int

	// Write to this file descriptor to make waitFD readable.
	signalfd C.int

	// Backing capture source
	source *captureSource
}

func max(x C.int, y C.int) C.int {
	if x &gt; y {
		return x
	}

	return y
}

//export mac68k_pcap_dispatchCallback
func mac68k_pcap_dispatchCallback(user *C.uchar, h *C.struct_pcap_pkthdr, bytes *C.uchar) {
	// TODO
}

func (server *pollServer) selector() {
	/* Determine the maxfd */
	var maxfd C.int
	maxfd = 0
	maxfd = max(server.pcapfd, maxfd)
	maxfd = max(server.waitfd, maxfd)
	maxfd &#43;= 1

	/* Configure fd sets */
	var master_readset C.fd_set
	C.MY_FD_ZERO(&amp;master_readset)

	C.MY_FD_SET(server.pcapfd, &amp;master_readset)
	C.MY_FD_SET(server.waitfd, &amp;master_readset)

	for {
		readset := master_readset

		ret, err := C.my_select(maxfd, &amp;readset, nil, nil, nil)
		if ret == -1 {
			// Shouldn&#39;t happen!
			fmt.Println(&#34;Unexpected select error&#34;, err)
		}

		/* The select timed out */
		if ret == 0 {
			continue
		}

		/* Check for completion */
		if C.MY_FD_ISSET(server.waitfd, &amp;readset) != 0 {
			fmt.Println(&#34;Cleaning up&#34;)
			C.close(server.waitfd)
			C.close(server.signalfd)
			break
		}

		/* Check for pcap readability */
		if C.MY_FD_ISSET(server.waitfd, &amp;readset) != 0 {
			/* Dispatch a read */
			C.pcap_dispatch(server.source.cptr, -1, unsafe.Pointer(C.pcap_dispatch_cb_handler), nil)
		}
	}
}

// Create a new poll server for the given capture source
func newPollServer(source *captureSource) (*pollServer, error) {
	server := new(pollServer)
	server.source = source

	/* Set up our error buffer */
	errbuf := (*C.char)(C.calloc(C.PCAP_ERRBUF_SIZE, 1))
	defer C.free(unsafe.Pointer(errbuf))

	/* Mark source non-blocking */
	if ret := C.pcap_setnonblock(source.cptr, 1, errbuf); ret != 0 {
		return nil, errors.New(C.GoString(errbuf))
	}

	/* Configure the fd-based signaling mechanism */
	var fds [2]C.int
	if ret, err := C.pipe(&amp;fds[0]); ret != 0 {
		return nil, fmt.Errorf(&#34;Failed to create signal pipe: %v&#34;, err)
	}

	server.waitfd = fds[0]
	server.signalfd = fds[1]

	server.pcapfd = C.pcap_get_selectable_fd(source.cptr)

	/* Fire off the background handler */
	go server.selector()

	return server, nil
}
