# EV Charging Network Hub - Protocol Research

## Executive Summary

For creating the world's largest electric vehicle charging network hub where charge point owners and vehicle owners pool together, the **combination of OCPP and OCPI** is the recommended approach. OCPP is the most widely adopted protocol globally, while OCPI enables the roaming and peer-to-peer functionality needed for your hub concept.

---

## 1. OCPP (Open Charge Point Protocol) - RECOMMENDED PRIMARY PROTOCOL

### Overview
OCPP is the **most widely adopted** communication protocol in the EV charging world and is considered the de-facto global standard for charging infrastructure interoperability.

### Global Adoption (2025)
- **137 countries** use OCPP as of January 2025
- **85% of EU charging stations** were OCPP-certified in 2022
- Approved as **IEC 63584** international standard in 2024
- Published by **CENELEC** as European standard in 2025

### Latest Version: OCPP 2.1 (Released January 2025)
Key features include:
- Distributed energy resources control
- Vehicle-to-Grid (V2G) capabilities
- Backwards compatible with OCPP 2.0.1
- Enhanced smart charging capabilities
- Improved load balancing
- Enhanced compliance with energy market requirements

### Government Support
**United States:**
- Required by Federal Highway Administration under NEVI Program
- Mandatory in New York and California for charging infrastructure funding

**Europe:**
- Required under Alternative Fuels Infrastructure Regulation (AFIR)
- Widespread adoption across all EU member states

**Asia:**
- South Korea requires OCPP certification for public funding
- Widely adopted across Asian markets

### Key Capabilities
- Remote monitoring and management of charging stations
- User authorization and access control
- Hardware configuration and firmware updates
- Payment and billing system integration
- Vendor independence (no proprietary lock-in)
- Scalability for large networks
- Security features for modern infrastructure

### Why Choose OCPP?
✅ **Largest adoption worldwide** (137 countries)
✅ **Open source and royalty-free**
✅ **Government-mandated in major markets**
✅ **Internationally standardized** (IEC 63584)
✅ **Strong ecosystem** of compatible vendors
✅ **Future-proof** with ongoing development

---

## 2. OCPI (Open Charge Point Interface) - RECOMMENDED FOR ROAMING

### Overview
OCPI is specifically designed for **peer-to-peer EV charging network roaming** and is the most popular roaming protocol, especially in North America.

### Key Purpose
Enables seamless communication between:
- **Charge Point Operators (CPOs)** - who manage charging stations
- **E-Mobility Service Providers (eMSPs)** - who serve EV drivers

### Latest Version: OCPI 2.3.0

### Key Capabilities
- Automated roaming between charging networks
- User authorization across different networks
- Real-time charge point information exchange
- Transaction event tracking
- Charge detail record (CDR) exchange
- Smart-charging command exchange between parties
- Scalable network architecture

### Management
- Published and managed by **EVRoaming Foundation**
- Publicly available and **royalty-free**
- Multiple compatible vendor implementations
- Used worldwide with broad industry adoption

### Why Choose OCPI?
✅ **Perfect for P2P charging hubs**
✅ **Enables network roaming** (charge anywhere)
✅ **Most popular roaming protocol in North America**
✅ **Open and royalty-free**
✅ **Designed for multi-operator scenarios**
✅ **Real-time data exchange**

---

## 3. ISO 15118 - Vehicle-to-Grid Communication

### Overview
ISO 15118 defines the communication protocol **between the EV and the charging station**, enabling advanced features like Plug & Charge and V2G.

### Key Features

**Plug & Charge:**
- Automatic vehicle identification and authorization
- Seamless charging experience (no cards/apps needed)
- Secure authentication using PKI

**Vehicle-to-Grid (V2G):**
- Bidirectional energy transfer
- EVs can supply energy back to the grid
- Essential for grid balancing with renewable energy

**Smart Charging:**
- Individual charging schedules per vehicle
- Grid capacity optimization
- Real-time demand management

**Security:**
- Cryptographic security mechanisms
- Digital certificates and PKI
- Secure data exchange

### Current Version: ISO 15118-20
Includes:
- Enhanced V2G support
- Improved security
- Better interoperability
- Advanced smart charging features

### Vehicle Support (2025)
Growing adoption in premium and newer vehicles:
- Porsche Taycan
- Mercedes-Benz EQS
- Lucid Air
- Ford Mustang Mach-E
- BMW i4, i5, i7, iX
- Hyundai Ioniq 5 & 6

### Why Consider ISO 15118?
✅ **Premium user experience** (Plug & Charge)
✅ **Future-proof** with V2G capabilities
✅ **Enhanced security**
✅ **Growing vehicle support**
⚠️ Requires compatible vehicles and chargers

---

## 4. Existing P2P Charging Platforms

### Examples of Successful P2P Platforms:

**Powerly:**
- Decentralized EV charging network
- Owners list chargers for public use
- Drivers find and book available chargers

**GoPlugable:**
- P2P EV charger rental platform
- Stripe payment integration
- Turn driveways into income streams

**JustCharge (UK):**
- Partners with Zap Map
- Residential charger booking
- App-based discovery

**EVmatch (US):**
- Concentrated on US coasts
- P2P charging marketplace

**Co Charger (UK):**
- Growing UK P2P network

### Common Features:
- Mobile app for discovery and booking
- Payment processing integration
- Owner schedule management
- Driver reviews and ratings
- Real-time availability

---

## 5. Physical Connector Standards

### By Region:

**Europe:**
- **Type 2 (Mennekes)** - AC charging
- **CCS2 (Combined Charging System)** - DC fast charging
- EU mandated CCS2 in 2014

**North America:**
- **CCS1** - Current standard
- **NACS (North American Charging Standard)** - Tesla's standard, being adopted by major automakers
- Ford and GM vehicles with NACS rolling out in 2025

**China:**
- **GB/T** - China's national standard

**Japan:**
- **CHAdeMO** - Still common in 2025
- **CCS2** - Growing adoption

---

## RECOMMENDATIONS FOR YOUR CHARGING HUB

### Primary Technology Stack:

**1. OCPP 2.1 (Core Communication Layer)**
- Use OCPP 2.1 for all charging station management
- Enables vendor-neutral hardware choices
- Future-proof with V2G and smart charging
- Government compliance in major markets
- Largest ecosystem and support

**2. OCPI 2.3.0 (Roaming & P2P Layer)**
- Implement OCPI for network roaming
- Enables charge point owners to share chargers
- Allows EV drivers to use any charger in the network
- Real-time availability and booking
- Automated billing across operators

**3. ISO 15118 (Optional - Premium Features)**
- Consider for future enhancement
- Provides Plug & Charge experience
- Enables V2G for revenue opportunities
- Requires compatible vehicles

### Architecture Benefits:

✅ **Largest Possible Network:** OCPP used in 137 countries
✅ **True P2P Capability:** OCPI designed for roaming
✅ **Vendor Independence:** Open standards prevent lock-in
✅ **Government Support:** Compliance with NEVI, AFIR
✅ **Scalable:** Both protocols handle massive networks
✅ **Future-Proof:** Active development and standardization

### Implementation Strategy:

1. **Phase 1:** Deploy OCPP 2.1 backend for charge point management
2. **Phase 2:** Add OCPI integration for roaming and P2P sharing
3. **Phase 3:** Mobile app for owners and drivers
4. **Phase 4:** Consider ISO 15118 for premium features

### Ecosystem Partners to Consider:

**Backend Platforms Supporting OCPP + OCPI:**
- AMPECO
- Driivz
- ChargeLab
- Virta
- Monta

---

## Conclusion

For the world's largest electric vehicle charging network hub with peer-to-peer capabilities, implement:

**🏆 OCPP 2.1** as your primary protocol (most adopted globally)
**🏆 OCPI 2.3.0** for enabling roaming and P2P sharing

This combination provides:
- Maximum global reach (137 countries)
- True peer-to-peer capability
- Vendor independence
- Government compliance
- Proven at scale
- Open source and royalty-free
- Strong industry support

**OCPP + OCPI is the industry-standard approach for large-scale, interoperable EV charging networks with roaming capabilities.**
