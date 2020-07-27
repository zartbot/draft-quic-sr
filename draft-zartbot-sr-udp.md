---
title: "Segment Routing over UDP(SRoU)"
abbrev: Segment Routing over UDP(SRoU)
docname: draft-zartbot-sr-udp-00
date: {DATE}
category: exp
ipr: trust200902
area: Transport
workgroup: SPRING

stand_alone: yes
pi: [toc, sortrefs, symrefs, docmapping]

author:
  -
    ins: K. Fang
    name: Kevin Fang
    org: Cisco Systems, Inc.
    email: zartbot.ietf@gmail.com
  -
    ins: Y. Li
    name: Yinghao Li
    org: Google, Inc.
    email: liyinghao@gmail.com
  -
    ins: F. Cai
    name: Feng Cai
    org:  Cisco Systems, Inc.
    email: fecai@cisco.com
  -
    ins: X. Jiang
    name: Xing Jiang
    org:  Cisco Systems, Inc.
    email: jamjiang@cisco.com

--- abstract

This document defines the Segment Routing Header{{!RFC8754}} extension
in UDP transport protocol with Network Address Translation Traversal.

--- middle

# Introduction

Many UDP based transport protocol(eg, IPSec/DTLS/QUIC) could provide a secure
transportation layer to handle overlay traffic. How ever it does not flexible
for source based path enforcement.

This document defines a new Segment Routing Header in UDP payload to enable 
segment routing over UDP(SRoU) for IPSec/DTLS/QUIC or other UDP based traffic.

Segment Routing over UDP(SRoU) interworking with QUIC could provide a generic 
programmable and secure transport layer for next generation applications.

Discussion of this work is encouraged to happen on GitHub repository which
contains the draft: <https://github.com/zartbot/draft-quic-sr>

## Specification of Requirements

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and
"OPTIONAL" in this document are to be interpreted as described in BCP 14
{{?RFC2119}} {{?RFC8174}} when, and only when,
they appear in all capitals, as shown here.

## Motivation

Segment Routing provides source-based path enforcement and transportation level
programmability but lacks of IPv4 support for transport over internet.

MPLS-over-UDP{{!RFC7510}} and MPLS Segment Routing over IP{{!RFC8663}}
defined SR-MPLS over IPv4 network, but it lacks of NAT traversal capabilities.

Many SDWAN vendors defined their private protocols for routing control over
multiple public cloud and internet, it’s hard for interop with multi-vendors.

Many applications may require intelligence traffic steering(CDN/LB case), 
SRoU with QUIC could be used in these cases.


# SR over UDP(SRoU) Packet encapsulation

The SRoU defined a generic segment routing enabled transport layer,the SR Header
insert in UDP payload.


~~~
  
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                       IP Header                         |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                       UDP Header                        |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                       SRoU Header                       |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                         |
 |                       Payload                           |
 |                                                         |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
~~~
{: #srou-encap title="SRoU encapsulation"}

## SR over UDP(SRoU) Header

SR over UDP must be present at the head of UDP payload. 

~~~
  0                   1                   2                   3
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | Magic Number  |  SRoU Length  | Flow ID Length| Protocol-ID   |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |                 Flow ID( Variable length)                     |
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                        Source Address                         |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |      Source Port              |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | Segment Type  |  SR Hdr Len   | Last Entry    | Segments Left |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |            Segment List[0] (length based on segment type)     |
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |                                                               |
                               ...
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |            Segment List[0] (length based on segment type)     |
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 //                                                             //
 //         Optional Type Length Value objects (variable)       //
 //                                                             //
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
~~~
{: #srh-format title="SRoU Header"}
 
Magic Number:
      1 Byte field
      For QUIC: could set to ALL ZERO to diffenciate with original header.
      For IPSec: could set to 0xFE value and avoid SPI allocation in
                 this range.
                 *0x00 may conflict with NON-ESP HEADER 
                 *0xFF may conflict with KeepAlive Message

SRoU Length:
      1 Byte, The byte length of a SRoU header.

FlowID Length:
      1 Byte, The byte length of FlowID field.

Protocol-ID:
  
| Type | Name        |Section                              |
|-----:|:------------|:------------------------------------|
|  0x0 | OAM         | for Link state probe and other OAM  |
|  0x1 | IPv4        | Indicate inner payload is IPv4 pkt  |
|  0x2 | IPv6        | Indicate inner payload is IPv6 pkt  |
{: #protocol-id title="Protocol ID field"}

Source Address:
       Protocol-ID = 1, this field is 4-Bytes IPv4 address
       Protocol-ID = 2, this field is 16-Bytes IPv6 address

Source Port:
       Source UDP Port Number

Segment Type:

| Type | Name                          | Len  |Section                 |
|-----:|:------------------------------|:-----|:-----------------------|
|  0x0 | Reserved                      |      |                        |
|  0x1 | IPv4 Address+ Port            | 48b  |{{ipv4-locator}}        |
|  0x2 | SRv6                          | 128b |{{srv6-locator}}        |
|  0x3 | Compressed Segment List       | 128b |{{cSID}}                |
{: #segment-types title="Segment Types"}

SR Hdr Len:
: SR Header length, include the SR Header flags  Segment-List and Optional TLV.
  
Last Entry:
: contains the index(zero based), in the Segment List, of
  the last element of the Segment List.

Segments Left:
: 8-bit unsigned integer. Number of route segemnts remaining,
  i.e., number of explicitly listed intermediate nodes still
  to be visited before reaching the final destination.

Segment List[0..n]:
: 128-bit/48-bit/144-bit addresses to represent the SR Policy.
  Detailed forwarding behavior will be defined in {{pkt-proccessing}}

TLV:
: Opptional TLV used for future extension.currently only defined
  the following TLV.

| Type | Value               | Len      |Section                 |
|-----:|:--------------------|:---------|:-----------------------|
|  0x0 | SR Integrity        | 32b      |{{sr-integrity}}        |
|  0x1 | Micro Segment Policy| variable |{{useg-policy}}         |
|  0x2 | End.PacketInfo      | variable |{{end-packet-info}}     |
{: #optional-tlv-types title="Optional TLV"}


# Packet Processing {#pkt-proccessing}

This section describe the packet proccessing procedure. The following
topology will be used in this section.

~~~
H1---R1----------I1------R3----------+---R4---H2
     |                               |
     |-----------R2------------------|
                 |
                 |
                 I2

I1,I2: Interim Node that support SRoU
R1~R4: Traditional Router
H1,H2: Host
~~~
{: #pp-topology title="Topology for packet proccesing"}

| Host | Address                  | SRoU Port| Post NAT         |
|-----:|:-------------------------|:---------|:-----------------|
|  H1  | 192.168.1.2              | 5111     |  10.1.1.1:23456  |
|  R1  | 192.168.1.1/10.1.1.1     |          |                  |
|  R2  | 10.1.2.2                 |          |                  |
|  R3  | 10.1.3.3                 |          |                  |
|  R4  | 10.1.4.4                 |          |                  |
|  H2  | 10.99.2.2                | 443      |  10.1.4.4:443    |
|  I1  | 10.99.1.1                | 8811     |                  |
|  I2  | 192.168.99.2             | 8822     |  10.1.2.2:12345  |
{: #ipv4-addr title="IP address table"}


## Type:0x1, IPv4 Locator Mode {#ipv4-locator}

In this mode, the endpoint could directly insert the interim node IPv4
addresses and port into the segment-list.

For example, H1 intend to send packet to H2 via R1-->I2---->H2, 
In this case SRoU packet will be NATed twice to show the NAT traversal workflow.
I2's public address could use STUN{{!RFC5389}} protocol detected and sync to all 
SRoU enabled devices.

H1 send packet with SRoU Header as below, H1 could use STUN detect it's source
public address, but consider the simplicity, the 1st hop SRoU forwarder cloud
update the source ip/port field in SRoU header.

~~~
IP/UDP Header {
  Source IP: 192.168.1.2,
  Destination IP: 10.1.2.2(SegmentList[1],I2 Pre-NAT public address),
  Source Port: 5111,
  Destination Port: 12345(SegmentList[1],I2 Pre-NAT public port),
}
SRoU Header {
  Magic Num = 0x0
  SRoU Length = 29
  FlowID Length = 0x3
  Protocol-ID = 0x1(IPv4),
  FlowID =  0x123,  
  Source Address = 192.168.1.2,
  Source Port = 5111,
  Segment Left = 0x1,
  Last Entry   = 0x1,
  SegmenetList[0] = 10.1.4.4:443(H2),
  SegmenetList[1] = 10.1.2.2:12345(I2),
}
~~~
{: #type-1-h1-i2 title="Type:0x1 H1-->I2 Packet Header"}


R1 is a NAT Device it will change the Source IP/Port to 10.1.1.1:23456.
But this router may not have ALG function to modify SRoU Header.Then packet
will send to 10.1.2.2:12345. It will be NAT again to I2.

After twice NAT, I2 Recieved packet as below:

~~~
IP/UDP Header {
  Source IP: 10.1.1.1(H1 post NAT addr),
  Destination IP: 192.168.99.2(I2 private addr),
  Source Port: 23456(H1 post NAT port),
  Destination Port: 8822(I2 private port),
}
SRoU Header {
  Magic Num = 0x0
  SRoU Length = 29
  FlowID Length = 0x3
  Protocol-ID = 0x1(IPv4),
  FlowID =  0x123,  
  Source Address = 192.168.1.2,
  Source Port = 5111,
  Segment Left = 0x1,
  Last Entry   = 0x1,
  SegmenetList[0] = 10.1.4.4:443(H2),
  SegmenetList[1] = 10.1.2.2:12345(I2),
}
~~~
{: #type-1-i2-recieved title="Type:0x1 H1-->I2, I2 Recieved Packet Header"}


if the (LastEntry == Segment Left) indicate I2 is the 1st hop SRoU forwarder,
It MUST apply ALG to update the Source Address/Port field by the IP/UDP header.
Then it will execute Segment Left - 1, and copy SegmentList[0] to DA/Dport.
Consider some interim router like R2 has URPF checking, the SA/Sport will also
updated to I2 SRoU socket address.

I2--->H2 packet:

~~~
IP/UDP Header {
  Source IP: 192.168.00.2(I2 Private),
  Destination IP: 10.1.4.4(SegmentList[0]),
  Source Port: 8822(I2 Private),
  Destination Port: 443(SegmentList[0]),
}
SRoU Header {
  Magic Num = 0x0
  SRoU Length = 29
  FlowID Length = 0x3
  Protocol-ID = 0x1(IPv4),
  FlowID =  0x123,  
  Source Address = 10.1.1.1(update by I2 ALG),
  Source Port = 23456(update by I2 ALG),
  Segment Left = 0x0(SL--),
  Last Entry   = 0x1,
  SegmenetList[0] = 10.1.4.4:443(H2),
  SegmenetList[1] = 10.1.2.2:12345(I2),
}
~~~
{: #type-1-i2h2 title="Type:0x1 I2-->H2 Packet Header"}

H2 will recieve the packet, and if the segment left == 0, it MUST copy the 
Source Address and Port into IP/UDP Header and strip out the SRoU Header and
send to udp socket. It may cache the reversed segmentlist for symmetric routing.

H2 send to UDP socket

~~~
IP/UDP Header {
  Source IP: 10.1.1.1(Copied from SRoU Src field),
  Destination IP: 10.99.2.2(Static NAT by R4),
  Source Port: 23456(Copied from SRoU Src field),
  Destination Port: 443(SegmentList[0]),
}
UDP Payload {
}
~~~
{: #type-1-h2tx title="Type:0x1 H2 Send to UDP socket"}

## Type:0x2, SRv6 format {#srv6-locator}

IPv6 does not need to consider the NAT traversal case, In this mode almost 
forwarding action is same as SRv6. This is only used for application driven
traffic steering(like CDN/LB usecase.). It has some benefit interworking with
QUIC, the pure userspace implementation could provide additional flexibility.

For example some IOT sensor with legacy kernel stack does not support SRv6 could
use SRoU insert SRH in UDP payload, the 1st hop SRoU forwarder could convert it
to standard SRv6 packet.

## Type:0x3, Compressed Segment List {#cSID}

### Service Registration & Mapping
I1,I2 use SRoU port as source port to inital STUN{{!RFC5389}} session to SR
mapping server, the mapping server could detect the Post NAT address and assign
SID for each host, and distribute IP/port--SID mapping database to all the SRoU
enabled host.

|Host  |  Socket                | SID      |
|-----:|:-----------------------|:---------|
|  I1  | 10.99.1.1:8811         | 1111     |
|  I2  | 10.1.2.2:12345         | 2222     |
{: #sid_map title="sid mapping"}

In this mode the socket information could combined with IPv4 and IPv6. 

## Optional TLV

### SR Integrity TLV {#sr-integrity}
 SR Integrity Tag to validate the SRH. All fields in the SRH except
 Segments Left fields need to be checked.

### Micro-segmentation(uSeg) {#useg-policy}
Option-TLV could defined Sub-TLV to support Micro-segmentation Security policy

~~~
  OptionTLV {
    0x1, uSeg{
        0x0, SRC_GROUP_ID,
        0x1, DST_GROUP_ID,
        0x2, APP_GROUP_ID,
        0x3, SRC_DEVICE_ID,
        0x4, DST_DEVICE_ID,
        0x5, APP_ID,
    }
  }
~~~
  
Customer also could encode this microsegment policy header in flowID field.

###  End.PacketInfo {#end-packet-info}
This optional TLV defines extened packet info and Segment-end packet edit
function. Sub-TLV defines as below:

#### Type:0x0, VPN-ID
 The SDWAN Router could use {{!I-D.ietf-quic-datagram}} as VPN tunnel, This
 Sub-TLV defined the VPN-ID inside the tunnel.

 If SRoU header has this sub-TLV, the device MUST decrypt inner payload and
 use the VPN-ID for inner packet destination lookup.

#### Type:0x1, Orginal Destination Address/Port
In SR Type 0x3, The original destination address/port cloud not encode in 128bit
field, it could be store in option TLV.




# OAM 

SRoU OAM Packet format is defined as below:

~~~
  0                   1                   2                   3
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | Magic Number  |  SRoU Length  | Flow ID Length|  P-ID  =0x0   |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |                 Flow ID( Variable length)                     |
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | OAM-Type      |   OAM Payload(Variable Length based on Type)  |
 +-+-+-+-+-+-+-+-+                                               +
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
~~~
{: #oam-format title="SRoU OAM Header"}

OAM-Type:

|ID    | Type                   | Usage                                       |
|-----:|:----------------------:|:--------------------------------------------|
|  0   | LinkState              | KeepAlive / Latency Measurement             |
|  1   | IPv4 STUN Request      |                                             |
|  2   | IPv4 STUN Response     |                                             |
|  3   | IPv6 STUN Request      | *Reserved for NAT66 Case(Not implement yet) |
|  4   | IPv6 STUN Response     | *Reserved for NAT66 Case(Not implement yet) |
{: #oam_type title="oam message type"}

## Link State

Each enpoint could initial this OAM message to its peer with local generated
sequence number and timestamp. This payload recommend to use private key 
encrypted. 

~~~
  0                   1                   2                   3
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+
 | LinkStateType |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                         Sequence Number                       |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                          TimeStamp                            |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
~~~
{: #oam-ls-format title="SRoU OAM Link State Header"}


The initiator send packet with LinkStateType = 0, The responder will modify this
flag to 1 and echo the entire packet back. Other type may defined for Two-way 
latency measurement which will be defined in later rfc version.

LinkStateType:

|ID    | Type                   | Usage                                       |
|-----:|:----------------------:|:--------------------------------------------|
|  0   | RTT_Request            |                                             |
|  1   | RTT_Response           |                                             |
{: #oam_linkstate_type title="oam linkstate message type"}


## STUN Service

SRoU forwarding endpoint may stay behind NAT, it request STUN service to 
discover the public network address.

Initiator send address and port with ALL-ZERO to STUN Server, STUN server
copy the recieve source address and port in this payload, and generate HMAC.
The STUN Server's key could be propogate to initiator by a out-of-band channel.

~~~
  0                   1                   2                   3
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                       IP Address                              |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |    Port                       |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                          HMAC                                 |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
~~~
{: #oam-stun-format title="SRoU OAM STUN Header"}


# Usage

## Traffic engineering over Internet

~~~
Client-------R1------------Internet--------------R2-----------Server
             |                                    |
             |                                    |
             R3----V1----PubliCloud--------V2-----|
~~~
{: #use-1 title="Traffic Engineering over internet"}

Many video/conferencing application requires traffic engineering over IPv4 
Internet, Webex/Zoom/Teams may setup V1,V2 in public cloud, The client and 
server could encode the V1/V2 information in SRoU header for traffic engineering

## Multipath forwarding

Same as previously topoloy {{use-1}}, customer cloud ask server transmit packet
over different path, two path have same Flow-ID, QUIC could be used in this case
to provide multistream/multihoming support.

## Micro Segmentation

Same as previously topoloy {{use-1}}, the interim Router: R1/R2/R3, V1,V2 could
insert uSeg Sub-TLV based on client and server uSeg identity, and other interim
network equipment could based on this sub-TLV implement security policy or QoS
policy.

## Container Network

~~~
C1----SideCar1-----L1-----S1------L2----SideCar2-------C2
                   |               |
                   |------S2-------|
C1,C2: Container
L1,L2: Leaf switch
S1,S2: Spine switch
~~~
{: #use-3 title="Service-Mesh & Container Network"}

SRoU with QUIC also could be used for container network interface, especially 
for service-mesh sidecar. The sidecar could aware the Datacenter underlay 
topology by BGP-LinkState, and use SRH select best path to avoid congestion. 
At the same time, all traffic are encrypted by {{!I-D.ietf-quic-tls}}.

## MPLS-SR with SDWAN 

~~~
S1---INET(ipv4)----PE1------MPLS------PE2----S2

S1,S2: SDWAN Router
PE1,PE2: SR enabled MPLS PE
~~~
{: #sr-sdwan-topology title="MPLS-SR with SDWAN"}

S1 will setup IPSec SA with S2 for end-to-end encryption,
And it will use BSID between PE1--PE2 for traffic engineering.

MPLS based BSID and IPv4 based locator could be encoded in uSID.A distributed
mapping table could be used to translate uSID to packet action.

~~~
IP/UDP Header {
  Source IP: H1,
  Destination IP: PE1,
  Source Port: srcport,
  Destination Port: IPSec,
}
SRoU Header {
  SegmentType = 0x1,
  SR_HDR_Len = 2,
  Last Entry = 0x0,
  Segment Left = 0,
  SegmenetList[0] = uSID: FC0:2222:3333:4444::
}
~~~
{: #type-1-s1-pe1 title="Type:0x1 S1-->PE1 Packet Header"}

## Cloud Native Network platform

Each of the SRoU forwarder only rely on a UDP socket, it could be implement
by a container. Customer could deploy such SRoU enable container in multiple
cloud to provide a cloud-angonostic solution. All containers could be managed
by K8S.

A distributed K-V store could be used for SRoU forwarder service registration,
routing(announce prefix), all the SRoU forwarder could measue peer's 
reachability/jitter/loss and update link-state to the K-V store. forwarding 
policy also could be sync by the K-V store. Detailed information will be
provided in another I.D(ETCD based disaggregated SDN control plane).

SRoU forwarder also could be implement by BPF for container communication. It 
will provide host level traffic engineering for massive scale datacenter to 
reduce the West-East traffic congestion.

The best practice for SRoU is working with QUIC.
SRoU with QUIC transport protocol provides the following benefit for SDWAN :

 * Stream multiplexing
 * Stream and connection-level flow control
 * Low-latency connection establishment
 * Connection migration and resilience to NAT rebinding
 * Authenticated and encrypted header and payload

SRoU add traffic-engineering and VPN capabilites for SDWAN.
Many existing SDWAN features could gain the benefits like:

 * TCP optimization
 * Packet duplication

# Security Considerations

The SRoU forwarder must validate the packet, FlowID could be used for source
validation. It could be a token based solution, this token could be assigned
by controller with a dedicated expire time. Source/Dest device ID and group 
cloud encode in flowid and signed by controller, just like JWT.

A blacklist on controller k-v store could be implemented to block device when 
the token does not expire.

# IANA Considerations

## SRoU with QUIC

The magic number in SRoU must be ZERO to distiguish with QUIC Long/Short 
packet format.


# Acknowledgements
{:numbered="false"}

The following people provided substantial contributions to this document:

- Bin Shi, Cisco Systems, Inc.
- Yijen Wang, Cisco Systems, Inc.
- Pix Xu, Cisco Systems, Inc.



