# draft-quic-sr
RFC draft for segment routing over UDP/QUIC

## Abstract

This document defines the Segment Routing Header(RFC8754) extension 
in QUIC transport protocol.
It will provide a new general purpose transportation layer with the following features:
* Secure [QUIC TLS]
* Reliable [QUIC Transport]
* Programmable [Segment Routing]

## Presentation

There is easy for understanding presentation available in:

<https://github.com/zartbot/draft-quic-sr/tree/master/slides>


## Contribution
Discussion of this work is encouraged to happen on GitHub repository which
contains the draft: 

Issue and PRs are welcome:
<https://github.com/zartbot/draft-quic-sr/pulls>

## Use case

1. Traffic Engineering over IPv4 internet
2. Client-less VPC access
3. CNI(Container Network Interface)
4. Wire and Wireless Converged Access
5. Cloud native network service platform

## Prototype
A working IPv4 based QUIC-SR application avaiable at
<https://github.com/zartbot/draft-quic-sr/tree/master/example_apps>

We just did some hack on quic-go to provide userspace quic support.

//quic-go create session
session, err := quic.DialAddr(*remoteSock, tlsConf, config)

//update QUIC-SR segmentlist and it could be runtime modified.
session.SetQUICSR([]string{1.1.1.1:2345,2.2.2.2:4567}, []byte{0x1, 0x2, 0x3})
