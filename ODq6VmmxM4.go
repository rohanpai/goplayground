// Interface to both live and offline pcap parsing.
package pcap

// TODO: Look at windows x64 winpcap

/*
#cgo linux LDFLAGS: -lpcap
#cgo darwin LDFLAGS: -lpcap
#cgo windows CFLAGS: -I C:/WpdPack/Include
#cgo windows LDFLAGS: -L C:/WpdPack/Lib -lwpcap
#include <stdlib.h>
#include <pcap.h>

// Workaround for not knowing how to cast to const u_char**
int hack_pcap_next_ex(pcap_t *p, struct pcap_pkthdr **pkt_header,
                      u_char **pkt_data) {
    return pcap_next_ex(p, pkt_header, (const u_char **)pkt_data);
}
*/
import "C"