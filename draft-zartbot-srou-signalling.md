---
title: "Distributed KV Store based Routing protocol for SR over UDP(SRoU)"
abbrev: Distributed KV Store based Routing protocol for SR over UDP(SRoU)
docname: draft-zartbot-srou-control-00
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

This document defines the Distributed KV store based routing protocol for
Segment Routing over UDP.

--- middle

# Introduction

This draft provides a contol plane support for SRoU(Segment Routing over UDP).

Discussion of this work is encouraged to happen on GitHub repository which
contains the draft: <https://github.com/zartbot/draft-quic-sr>

## Specification of Requirements

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and
"OPTIONAL" in this document are to be interpreted as described in BCP 14
{{?RFC2119}} {{?RFC8174}} when, and only when,
they appear in all capitals, as shown here.

## Motivation

SRoU support udp transport session over internet, but it lack of reachability
detection and routing control, existing routing protocol like BGP-EVPN did not
provide Dynamic NAT traversal capability. 

This document provide a distributed KV store based routing protocol for SRoU.

## Overview

The routing protocol is based on source routing, each of the ingress node 
cloud get the overlay prefix and dest location mapping from distributed KV 
store, then the ingress node could fetch linkstate database from this KV store
and execute A* algorithm to search the candidate path which meet the SLA 
requirement.


# Node abstraction and registration

Each Node has the following attribute 

1. Role:  the system contains different node type, role attribute is a uint16
value which contains: 

| Type | Name        |Section                                                                                                                                                              |
|-----:|:------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|  0x0 | STUN        | This node is used as a STUN server to help other nodes discovery their public address.This node must deploy with a public internet address or behind static 1:1 NAT |
|  0x1 | Fabric      | This node type is used as a interim node to relay the SRoU traffic, this node MUST initial TWAMP link probe to other Fabric node and report linkstate to KV Store.  |
|  0x2 | Linecard    | This node type is used to connect existing network, it   could use TWAMP probe other Fabric Node or Linecard node                                                   |
|-----:|:------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
{: #node role title="Node Role"}

2. SiteID: uint32 number, defined the node which belongs to same site or 
Automomous System.
3. SystemName: unique string type to indicate a node.
4. Label: unique 24bit value, allocation algorithm is described in the following
section.
5. Location: Optional filed. It contains two float32 value(latitude and 
longitude) to indicate the Geo location.

## Node Label allocation

Each node initial TLS session to Distributed KV Store, and fetch a distributed
lock with key "/lock/systemlabel". The node will fetch prefix "/systemlabel" to
get all label mapping once it get the lock. Then it will assign the smallest 
unpresent int "X" in the list as it's system label, and register it to KV store 
by key="/systemlabel/X", then it could release the distributed lock. All of the
fabric node MUST listen the "/systemlabel" to update it's local node mapping 
table, Linecard node may fetch the "/systemlabel" key when it need to optimize
the local route.

This System Label could be used for cSID encoding or VPN based client linecard
node convert to it's tunnel address.

## Node registration

Each node will send Key="/node/role/systemName" and Value=" SiteID,SystemLabel,
Lat,Long" to the distributed KV store.


# SRoU Locator and Route
Each node may have multiple underlay socket which may behind the dynamic NAT, 
it MUST fetch the STUN list from "/node/stun" and "/service/stun" to get 
the STUN server address list, then send the SRoU OAM-STUN packet to the random
selected stun server to get the public address.

Once the socket get the public address, it will encode the udp socket info as a
SRoU Locator:

"SystemName/Color/LocalIP:Port/PublicIP:Port/LocalInterface/TXBW/RXBW"

If the local socket has public address and port information, it could be added
in the service list.

The node MUST update it local servicelist to distributed KV store by:
Key= "/service/role/systemName"
Value= "SRoULocator1,SRoULocator2"

# Node Keepalive

Each KV pair registration MUST have a leasetime and keepalive timer, Once the
Node out of service and disconnected, the KV store MUST withdraw the KV pair 
after lease timeout.

# Link State
Each Fabric Node must watch the "/service/fabric" key prefix to update its local
SRoU Service list database. It MUST initial TWAMP session over the service udp
socket to measure the link performance and reachablity.

Linkstate measurement result COULD send to the KV store to construct the 
linkstate Database by the following Key Value type:

Key="/stats/linkstate/SRC_SRoU_Locator->DST_SRoU_Locator"
value= TWAMP measured jitter/delay/loss result and underlay interface load.

The Node CPU,Memory usage also could be updated by:
Key="/stats/node/SystemName"
Value="CPULoad,MemoryUsage"

An telemetry analytics node could watch key prefix ="/stats" for assurance
and AIOps based routing optimization.

# Sercurity Key

Each node may update it node key or per socket key , or per session pair key to
the KV Store:

Key="/key/SystemName"
Value="Key1,Key2"

Key="/key/socket/SRoU_Locator"
Value="Key1,Key2"

Key="/key/session/SRC_SRoU_Locator->DST_SRoU_Locator"
Value="Key1,Key2"

During Rekey, the node must update both OldKey and newKey to the KV Store and 
accept both Key in a while to wait the entire system sync to the new key.

# Overlay Routing

RouteDistinguish could encode by SystemName + local VNID
The overlay routing prefix is encoded as below:

Type-2 EVPN Route
Key="/route/2/exportRT/RD/MAC/IP"
Value="VNID/SystemName/PolicyTag"

Type-5 EVPN Route

Key="/route/5/exportRT/RD/IPPrefix/IPMask"
Value="VNID/SystemName/PolicyTag"

Each of the linecard node could based on import RT list to watch key 
prefix ="/route/2/importRT" and "/route/5/importRT" to sync the routing table.

Each linecard node could selective fetch the "/stats/linkstate" to get the 
toplogy information and execute flexibile algorithm(SPF,A* search) to calculate
the candidate path, then enforce it to its forwarding table.

# Control Policy

## Route control
Inspired by BGP FlowSpec,Network operator could update the control policy to 
the entire system by using:

Key="/control/RT/2/SRC_MAC/SRC_IP/DST_MAC/DST_IP"
Key="/control/RT/5/SRC_Prefix/SRC_Mask/DST_Prefix/DST_Mask"
Value="Action" /"SR Locator list"

## Access Control

Each node may use the SRoU flowID field as a token based access control.
This token could grant or revoke by a policy engine.

Key="/token/permit/flowid"
Key="/token/block/flowid"

Each node could sync this table to execute the access control policy.

## User identity

Each of the endpoint may have it's identity or group policy tags, it could be
updated by

key="/identity/userid/user_device_id"
value="group policy tags"

Group policy could be updated and store in ETCD by

key="/policy/src_grp/dst_grp"
value="actions"


# Distributed KV Store
ETCD is used in our prototype, we deploy an etcd cluster in main datacenter
and place many of the proxy node on public cloud to make sure the node could 
be available connect to the entire system. In some on-prem deployment, each of
the nodes could act as a ETCD proxy to help other node register to KV store.

# Security Considerations

All of the control connection is TLS based and MUST validate the server and
client certification.


# IANA Considerations


# Acknowledgements
{:numbered="false"}

The following people provided substantial contributions to this document:

- Yijen Wang, Cisco Systems, Inc.



