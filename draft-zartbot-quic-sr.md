---
title: "Segment Routing over QUIC(QUIC-SR)"
abbrev: Segment Routing over QUIC
docname: draft-zartbot-quic-sr-00
date: {DATE}
category: exp
ipr: trust200902
area: Transport
workgroup: QUIC

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


normative:

  MICRO-SEG-DEF:
   title: "Network Micro-Segmentation vs. Traditional Network Segmentation"
   author:
    - ins: Elena Garrett
   date: 2018-02
   target: "https://www.linkedin.com/pulse/adaptive-network-micro-segmentation-elena-garrett/"


--- abstract

This document defines the Segment Routing Header{{!RFC8754}} extension
in QUIC transport protocol.

--- middle

# Introduction

The QUIC Transport Protocol {{!I-D.ietf-quic-transport}} provides a secure,
multiplexed connection for transmitting reliable streams of application data.
CONNECTION-ID provides IP address independently transportation, thus allowing
multipath transportation.

Enable Segment routing for QUIC(QUIC-SR) provides more flexibility for endpoint
to select multipath to reduce latency and avoid packet drop.

This document defines a new QUIC-SR packet types.

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

QUIC provides connection migration and CONNECTION-ID definition which could
provide IP address independent forwarding mechanism and cloud be easily use
for NAT traversal and easy for load balance {{!I-D.ietf-quic-load-balancers}}.

QUIC also supports reliable transportation for general purpose applications and
unreliable transportation {{!I-D.ietf-quic-datagram}} for low-latency
applications.

QUIC provides a secure transportation and low-latency connection establishment.
It could be the best general-purpose transport protocol for SD-WAN.

Add Segment Routing  as a new packet type in QUIC will provide more
flexibility and general-purpose programmability for multi-vendor and multi-cloud
deployment.

# QUIC-SR Packet Type

The Segment Routing information MUST be present as clear-text in data payload.
A new QUIC-SR packet type is required.

## QUIC-SR Packet

QUIC-SR Packet means a QUIC packet type that could be encapsulated
in UDP payload.

~~~
QUIC-SR Packet {
  Header Form (1) = 1,
  Fixed Bit (1) = 1,
  Long Packet Type (2) = 0,
  -----------------------
  QUIC-SR Flag(1) =  1,
  Unused (3),
  -----------------------
  Version (32),
  DCID Length (8),
  Destination Connection ID (0..160),
  SCID Length (8),
  Source Connection ID (0..160),
  QUIC-SR Header (..),
}
~~~
{: #sr-pkt-format title="Segment Routing QUIC Packet Type"}

QUIC Long Header packet type(2bits) is fully used, Need to work with QUIC-WG
to extend one bit for Segment Routing Packet type.

## QUIC-SR Header

QUIC-SR Header is defined in the following figure:

~~~
  0                   1                   2                   3
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | Segment Type  |  SR Hdr Len   | Last Entry    | Segments Left |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |            Segment List[0] (32-bit/ 64-bit / 128-bit )        |
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
 |            Segment List[n] (32-bit/ 64-bit / 128-bit )        |
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 //                                                             //
 //         Optional Type Length Value objects (variable)       //
 //                                                             //
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
~~~
{: #srh-format title="QUIC-SR Header"}

Segment Type:
      Defined 4 type of segment list
  
| Type | Name                          | Len |Section                 |
|-----:|:------------------------------|:----|:-----------------------|
|  0x0 | IPv4 Locator Only             | 32b |{{ipv4-locator}}        |
|  0x1 | 64bit Segment Mode            | 64b |{{type1-64b-locator}}   |
|  0x2 | IPv6 SRH                      | 128b|{{srv6-locator}}        |
{: #segment-types title="Segment Types"}

SR Hdr Len:
: QUIC-SR Header length

Last Entry:
: contains the index(zero based), in the Segment List, of
  the last element of the Segment List.

Segments Left:
: 8-bit unsigned integer. Number of route segemnts remaining,
  i.e., number of explicitly listed intermediate nodes still
  to be visited before reaching the final destination.

Segment List[0..n]:
: 32-bit/64-bit/128-bit addresses to represnet the SR Policy.
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
H1---+-----R1------I1------R3------R4------H2
     |                              |
     |-----R2------I2---------------|

I1,I2: Interim Node that support QUIC-SR
R1~R4: Traditional Router
H1,H2: Host
~~~
{: #pp-topology title="Topology for packet proccesing"}

## Type:0x0, IPv4 Locator Mode {#ipv4-locator}

In this mode, the endpoint could directly insert the interim node IPv4
addresses into the segment-list.

For example, H1 intend to send packet to H2 via R1-->I1-->R3-->R4,

H1 send packet with SRH as below:

~~~
IP/UDP Header {
  Source IP: H1,
  Destination IP: I1,
  Source Port: srcport,
  Destination Port: QUIC-SR,
}
QUIC-SR Header {
  SegmentType = 0x0,
  SR_HDR_Len = 2,
  Last Entry = 0x1,
  Segment Left = 1,
  SegmenetList[0] = H2,
  SegmenetList[1] = I1,
}
~~~
{: #type-0-h1-I1 title="Type:0x0 H1-->I1 Packet Header"}

I1 is QUIC-SR enabled Node,It will swap the IP Header and reduce
segment left field, then forward packet to H2.

~~~
IP/UDP Header {
  Source IP: H1,
  Destination IP: H2,
  Source Port: srcport,
  Destination Port: QUIC-SR,
}
QUIC-SR Header {
  SegmentType = 0x0,
  SR_HDR_Len = 2,
  Last Entry = 0x1,
  Segment Left = 0,
  SegmenetList[0] = H2,
  SegmenetList[1] = I1,
}
~~~
{: #type-0-I1-H2 title="Type:0x0 I1-->h2 Packet Header"}

## Type:0x1, 64bit Segment mode {#type1-64b-locator}

64bit segment mode is encoded as below:

~~~
  0                   1                   2                   3
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | Function                      | Args          | T | Reserved  |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                            Locator                            |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

 T: defined locator type, 0x0=IPv4,0x1=MPLS,0x2=VNID

~~~
{: #type-1-format title="64bit SR Header in QUIC"}

Currently no plan to implement this type. But this type could be used in the
following scenario:

 - MPLS-SR with SDWAN : combined IPv4 + SID locators
 - IPv4 Locator(32b) + Function(16b) + Args(8b)

## Type:0x2, IPv6 SRH {#srv6-locator}

In this mode, If the QUIC-SR header does not contain optional TLV,Interim Node
could copy the SRH to IPv6 Header, and remove the QUIC-SR Packet in UDP payload.

If the QUIC-SR contains optional TLV, the segment forwarding behaviour is same
as {{RFC8754}}, but the last node must proccesing optional TLV.

## NAT Travasal and dynamic port mapping

QUIC is encapsulated in standard UDP packet, to ensure the packet send to
correct quic socket.2 options are listed for future usage

### Option.1 STUN Server mode

Each of the QUIC-SR enabled device use QUIC-SR port listen INADDR_ANY, then use
QUIC-SR port as source port to send packet to STUN{{!RFC5389}} Server to
discovery the NATted public ip port. Then it must send this mapping to a K-V
store for other QUIC-SR enabled device query.

A pre-allocated Segment-ID could be used as key:

| SID    | private IP:Port            |Public IP:Port              |
|-------:|:---------------------------|:---------------------------|
|  111   | 192.168.1.2:QUIC-SR        | 1.1.1.1:12345              |
|  222   | 192.168.1.2:QUIC-SR        | 2.2.2.2:45312              |
{: #quic-node-sid title="QUIC-SR node SID and pub/priv mapping"}

Then the QUIC-SR header could use uSID{{!I-D.filsfils-spring-net-pgm-extension-srv6-usid}} 
for path encoding to reduce overhead.

### Option.2 Peer Cache mode

When device recieve upstream packet on QUIC-SR socket, it MUST cache the
CONNECTION-ID and source IP/Port, then resend by it's own socket to nexthop.


## Optional TLV

### SR Integrity TLV {#sr-integrity}
 SR Integrity Tag to validate the SRH. All fields in the SRH except
 Segments Left fields need to be checked.

### Micro-segmentation(uSeg) {#useg-policy}
Micro-segmentation is a cyber security technique that segments the network based
on a diverse set of variables to describe various security zones. In addition to
utilizing IPs and VLANs, micro-segmentation allows security zones(microsegments)
to be described using various host-centric attributes (such as OS, hardware, 
behavior) or application-driven attributes (such as where the workload is 
residing and how it is behaving) in order to improve visibility and security 
around those attributes and behaviors.{{MICRO-SEG-DEF}}  

Option-TLV could defined Sub-TLV to support Micro-segmentation

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
  
This policy tag could be added by interim network device, and the src/dst
group id and device id could be encoded in CONNECTION-ID to reduce packet 
overhead.


###  End.PacketInfo {#end-packet-info}
This optional TLV defines extened packet info and Segment-end packet edit
function. Sub-TLV defines as below:

#### Type:0x0, VPN-ID
 The SDWAN Router could use {{!I-D.ietf-quic-datagram}} as VPN tunnel, This
 Sub-TLV defined the VPN-ID inside the tunnel.

 If QUIC-SR header has this sub-TLV, the device MUST decrypt inner payload and
 use the VPN-ID for inner packet destination lookup.

#### Type:0x1, Edit Destination Port
 The interim QUIC-SR enabled router need to proccess packet on dedicated port,
 this sub-TLV is used to store the original quic socket destination port.

 If Segment Left = 0, the Device must replace it's destination port based on
 the sub-TLV defined value.

 For exmaple, topology is shown as {{pp-topology}}, H1 must send traffic to I1
 with Destination Port = QUIC-SR.

~~~
IP/UDP Header {
  Source IP: H1,
  Destination IP: H2,
  Source Port: srcport,
  Destination Port: QUIC-SR,  //sr-quic-port
}
QUIC-SR Header {
  SegmentType = 0x0,
  SR_HDR_Len = 2,
  Last Entry = 0x1,
  Segment Left = 1,
  SegmenetList[0] = H2,
  SegmenetList[1] = I1,
  OptionTLV {
    0x1, PacketInfo{
        0x1, DestinationPort: PortH2
    }
  }
}
~~~
{: #pmap-H1-I1 title="PortMapping H1-->I1 Packet Header"}

When I1 recieved the packet, and reduce the Segment Left field to 0.
It will trigger I1 proccess option TLV to change the dest port to PortH2

~~~
IP/UDP Header {
  Source IP: H1,
  Destination IP: H2,
  Source Port: srcport,
  Destination Port: PortH2,  //h2 port
}
QUIC-SR Header {
  SegmentType = 0x0,
  SR_HDR_Len = 2,
  Last Entry = 0x1,
  Segment Left = 0,
  SegmenetList[0] = H2,
  SegmenetList[1] = I1,
  OptionTLV {
    0x1, PacketInfo{
        0x1, DestinationPort: PortH2
    }
  }
}
~~~
{: #pmap-I1-H2 title="PortMapping I1-->H2 Packet Header"}

# Usage

## Traffic engineering over Internet

~~~
Client-------R1------------Internet--------------R2-----------Server
             |                                    |
             |                                    |
             R3----V1----PubliCloud--------V2-----|
~~~
{: #use-1 title="Traffic Engineering over internet"}

network traffic will go through R1 and R2, but it's congested or link failured.
Cloud provider or network admin cloud setup QUIC-SR enabled NFV on public cloud.
Client could use V1/V2 forward traffic and avoid congestion.

## Multipath forwarding

Same as previously topoloy{{use-1}}, customer cloud ask server transmit packet
over different path, two path have same CONNECTION-ID, and could be shared with
different QUIC stream.

In this scenario, V1,V2 could implement Source address/port based NAT to make
sure the server-->client forwarding path are symmetric.

## Micro Segmentation

Same as previously topoloy{{use-1}}, the interim Router: R1/R2/R3, V1,V2 could
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

QUIC-SR also could be used for container network interface, especially for 
service-mesh sidecar. The sidecar could aware the Datacenter underlay topology
by BGP-LinkState, and use SRH select best path to avoid congestion. At the same
time, all traffic are encrypted by {{!I-D.ietf-quic-tls}}.


## Standard SDWAN Encapsulation

In SDWAN era, different vendor defined different encapsulation format, it's very
hard for interop accross multiple cloud. QUIC-SR gives a standard encapsulation
proposal.

~~~
QUIC-SR Header {
  SegmentType = 0x0,
  SR_HDR_Len = 2,
  Last Entry = 0x1,
  Segment Left = 0,
  SegmenetList[0] = H2,
  SegmenetList[1] = I1,
  OptionTLV {
    0x1, PacketInfo{
        0x0, VPN-ID: VPN-ID
    }
  }
}
~~~
{: #std-sdwan title="Standard SDWAN"}

QUIC transport protocol provides the following benefit for SDWAN :

 * Stream multiplexing
 * Stream and connection-level flow control
 * Low-latency connection establishment
 * Connection migration and resilience to NAT rebinding
 * Authenticated and encrypted header and payload

QUIC-SR add traffic-engineering and VPN capabilites for SDWAN.
Many existing SDWAN features could gain the benefits like:

 * TCP optimization
 * Packet duplication

## MPLS-SR with SDWAN 

~~~
S1---INET(ipv4)----PE1------MPLS------PE2----S2

S1,S2: SDWAN Router
PE1,PE2: SR enabled MPLS PE
~~~
{: #sr-sdwan-topology title="MPLS-SR with SDWAN"}

S1 will setup QUIC socket with S2 for end-to-end encryption,
And it will use BSID between PE1--PE2 for traffic engineering.

MPLS based BSID and IPv4 based locator could be encoded in uSID.A distributed
mapping table could be used to translate uSID to packet action.

~~~
IP/UDP Header {
  Source IP: H1,
  Destination IP: PE1,
  Source Port: srcport,
  Destination Port: QUIC-SR,
}
QUIC-SR Header {
  SegmentType = 0x1,
  SR_HDR_Len = 2,
  Last Entry = 0x0,
  Segment Left = 0,
  SegmenetList[0] = uSID: FC0:2222:3333:4444::
}
~~~
{: #type-1-s1-pe1 title="Type:0x1 S1-->PE1 Packet Header"}


# IANA Considerations

## QUIC-SR Packet
This document registers a new value in the QUIC Long Long Packet Header:

Original:

~~~
Long Header Packet {
  Header Form (1) = 1,
  Fixed Bit (1) = 1,
  ----------------------
  Long Packet Type (2),
  Unused (4),
  ----------------------
  ...omit...
}
~~~

Changed if this document is approved:

~~~
QUIC-SR Packet {
  Header Form (1) = 1,
  Fixed Bit (1) = 1,
  ----------------------
  Long Packet Type (2) = 0,
  QUIC-SR Flag(1) =  1,
  Unused (3),
  ----------------------
  ...omit...
}
~~~

## QUIC-SR Port

A dedicated UDP port need to be allocated to QUIC-SR

# Acknowledgements
{:numbered="false"}

The following people provided substantial contributions to this document:

- Bin Shi, Cisco Systems, Inc.
- Yijen Wang, Cisco Systems, Inc.
- Pix Xu, Cisco Systems, Inc.
- Xing James Jiang, Cisco Systems, Inc.
