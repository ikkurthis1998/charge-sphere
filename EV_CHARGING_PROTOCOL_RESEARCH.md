# EV Charging Network Roaming Hub - Protocol Research

## Executive Summary

For creating the world's largest electric vehicle charging network **roaming hub** where charge point owners and vehicle owners pool together to use each other's chargers, **OCPI (Open Charge Point Interface) 2.3.0** is the recommended primary protocol.

**Important Clarification:** This research focuses on **roaming protocols** (network-to-network communication), NOT charger management protocols like OCPP (charger-to-server). Your hub will connect existing charging networks together, not manage individual charging stations directly.

---

## Understanding EV Roaming Architecture

### The Three Layers of EV Charging

1. **Vehicle ↔ Charger** (ISO 15118)
   - Communication between the EV and charging station
   - Out of scope for your hub

2. **Charger ↔ Backend Server** (OCPP)
   - Individual charging station management
   - Handled by Charge Point Operators (CPOs)
   - Out of scope for your hub

3. **Network ↔ Network** (OCPI, OICP, eMIP, OCHP) ← **YOUR HUB OPERATES HERE**
   - Roaming between different charging networks
   - Enables cross-network charging
   - **This is what you need!**

---

## 1. OCPI (Open Charge Point Interface) - 🏆 RECOMMENDED PRIMARY PROTOCOL

### Overview
OCPI is specifically designed for **automated roaming between EV charging networks**. It's the most widely adopted open protocol for network-to-network communication and peer-to-peer roaming.

### Latest Version: OCPI 2.3.0

### Key Purpose
Enables seamless communication between:
- **Charge Point Operators (CPOs)** - who own and operate charging stations
- **E-Mobility Service Providers (eMSPs)** - who provide charging services to EV drivers

### How It Works
OCPI supports **direct peer-to-peer roaming** where networks can connect directly to each other without requiring a central intermediary. This decentralized approach allows for:
- Direct business relationships between networks
- Lower transaction costs
- Greater autonomy
- Faster deployment

### Global Adoption (2025)
- **Most popular roaming protocol in North America**
- Widely adopted across Europe, Asia, and globally
- Even major hub operators (Hubject, GIREVE) are adopting OCPI:
  - **Hubject** announced native OCPI support in 2025 (dual-stack with OICP)
  - **GIREVE** has supported OCPI since 2018 alongside eMIP

### Key Capabilities
- ✅ Automated roaming authorization
- ✅ Real-time charge point information exchange
- ✅ Location and availability data
- ✅ Transaction event tracking
- ✅ Charge Detail Record (CDR) exchange
- ✅ Smart-charging command exchange
- ✅ Tariff and pricing information
- ✅ Booking module (OCPI 2.3) for reservations
- ✅ Token management and validation
- ✅ Session management

### Management & Licensing
- **Published and managed by**: EVRoaming Foundation
- **License**: Publicly available and **royalty-free**
- **Source**: Open standard with multiple vendor implementations
- **Documentation**: Comprehensive specifications on GitHub

### Why Choose OCPI?

✅ **Most widely adopted open roaming protocol globally**
✅ **Peer-to-peer architecture** (decentralized, no mandatory hub)
✅ **Fully open source and royalty-free**
✅ **Active development** with regular updates
✅ **Industry convergence** - even proprietary hubs adopting it
✅ **Perfect for P2P charging marketplaces**
✅ **Designed for multi-operator scenarios**
✅ **Scalable architecture**
✅ **Real-time data exchange**

### OCPI Adoption by Major Players
- Hubject (adding native support in 2025)
- GIREVE (supported since 2018)
- EVRoaming Foundation
- Hundreds of CPOs and eMSPs worldwide
- Major software platforms (AMPECO, Driivz, ChargeLab, Virta, Monta)

---

## 2. Existing Major Roaming Hubs

If you prefer to integrate with existing roaming hubs rather than building your own, here are the major players:

### 🌍 Hubject - The Largest Roaming Hub

**Protocol**: OICP (Open InterCharge Protocol) + OCPI (2025)

**Scale (2025)**:
- **2,750+ partner networks**
- **1+ million charging points**
- **70+ countries**
- Largest global coverage

**Protocol Evolution**:
- Originally used proprietary OICP protocol
- Open-sourced OICP in 2019
- **Announced native OCPI support in 2025** (dual-stack approach)
- General availability of OCPI support: late 2025

**Why This Matters**:
Even the world's largest roaming hub is adopting OCPI, signaling industry-wide convergence on this standard.

**Geographic Focus**: Global, strong in Europe and North America

---

### 🇪🇺 GIREVE - European Leader

**Protocol**: eMIP (eMobility Interoperation Protocol) + OCPI

**Scale (2025)**:
- **659,000 charging points** (October 2025)
- **~30 European countries**
- Pan-European coverage
- Hundreds of partners

**Protocol Details**:
- eMIP is GIREVE's proprietary protocol (but free to use with registration)
- **Also supports OCPI since 2018**
- Requires certification to connect to GIREVE platform

**Founded by**: EDF, Renault, CNR, and Caisse des Dépôts

**Geographic Focus**: France and Southern Europe (expanding pan-European)

---

### 🇩🇪 e-clearing.net - German/European Hub

**Protocol**: OCHP (Open Clearing House Protocol)

**Scale**:
- **1,200+ partners**
- **445,000 charging points**
- Primarily European coverage

**Protocol Approach**:
OCHP uses a **centralized clearing house model** where all roaming communications flow through a central intermediary platform.

**Latest Version**: OCHP 1.4

**Key Principles**:
- **Transparency**: Clearing house is invisible to end users
- **Independence**: Business models not influenced by hub
- **Anonymity**: Minimal private user data required

**Geographic Focus**: Europe, especially Germany

---

## 3. Roaming Protocol Comparison

### OCPI vs OICP vs eMIP vs OCHP

| Feature | OCPI | OICP | eMIP | OCHP |
|---------|------|------|------|------|
| **Architecture** | Peer-to-peer (decentralized) | Hub-based (Hubject) | Hub-based (GIREVE) | Hub-based (e-clearing.net) |
| **License** | Open source, royalty-free | Open source (since 2019) | Free with registration | Open source |
| **Management** | EVRoaming Foundation | Hubject | GIREVE | e-clearing.net |
| **Adoption** | Global, most widely adopted | 2,750+ partners, 70+ countries | 30 EU countries | 1,200+ partners, Europe |
| **Latest Version** | 2.3.0 | 2.3 | 0.7.4 | 1.4 |
| **Geographic Strength** | Global, esp. North America | Global | France, Southern Europe | Germany, Europe |
| **Requires Hub Membership** | No (direct P2P) | Yes (Hubject) | Yes (GIREVE) | Yes (e-clearing.net) |
| **Certification Required** | No | No | Yes | No |
| **Real-time** | Yes | Yes | Yes (but supports async) | Yes |
| **Smart Charging** | Yes | Yes | Yes | Limited |
| **Reservations** | Yes (2.3) | Yes | Yes | Limited |
| **2025 Trend** | Being adopted by hubs | Adding OCPI support | Already supports OCPI | Established |

---

## 4. Roaming Models Explained

### Peer-to-Peer Roaming (OCPI)
```
Network A ←→ Network B
Network A ←→ Network C
Network B ←→ Network C
```

**Advantages**:
- Direct business relationships
- No middleman fees
- Full control over partnerships
- Faster deployment
- Lower transaction costs

**Challenges**:
- Requires bilateral agreements
- More connections to manage
- Each network must implement protocol

---

### Hub-Based Roaming (OICP, eMIP, OCHP)
```
Network A ←→ HUB ←→ Network B
Network C ←→ HUB ←→ Network D
```

**Advantages**:
- Single connection to hub reaches all partners
- Hub manages all relationships
- Simplified integration
- Instant access to large network

**Challenges**:
- Hub membership fees
- Less control over partnerships
- Dependent on hub operator
- Potential higher transaction costs

---

### Hybrid Model (Your Opportunity!)
```
Network A ←→ YOUR HUB ←→ Network B
     ↓                        ↓
   OCPI                     OCPI
     ↓                        ↓
 Hubject Hub            GIREVE Hub
```

Build a hub that:
- Uses OCPI for direct peer-to-peer connections
- Also integrates with existing hubs (Hubject, GIREVE, e-clearing.net)
- Provides the largest possible network reach
- Offers both P2P and hub-based roaming

---

## 5. Successful P2P Charging Platforms

These platforms demonstrate successful peer-to-peer charging sharing implementations:

### Powerly
- Decentralized EV charging network
- Charger owners list their chargers for public use
- EV drivers find and book available chargers
- Uses blockchain for decentralized transactions

### GoPlugable (UK)
- P2P EV charger rental platform
- Stripe payment integration
- "Turn your driveway into a passive income stream"
- Mobile app for discovery and booking

### JustCharge (UK)
- Partners with Zap Map for charger discovery
- Residential charger booking
- App-based platform
- Growing UK network

### EVmatch (US)
- P2P charging marketplace
- Strong presence on US East and West coasts
- Mobile app platform
- Homeowner and driver focused

### Co Charger (UK)
- Growing UK P2P network
- Community-focused approach
- Residential charger sharing

### Common Features Across P2P Platforms:
- Mobile app for charger discovery
- Real-time availability display
- Online booking and scheduling
- Payment processing integration
- Owner schedule management
- User reviews and ratings
- Pricing flexibility for owners

---

## 6. RECOMMENDATIONS FOR YOUR ROAMING HUB

### Primary Strategy: Build on OCPI

**🏆 Implement OCPI 2.3.0 as your core protocol**

### Why This Is The Best Choice:

1. **Industry Convergence**: Even proprietary hubs (Hubject, GIREVE) are adopting OCPI
2. **Open & Free**: No licensing fees or vendor lock-in
3. **Most Widely Adopted**: Global standard for roaming
4. **Future-Proof**: Active development and industry support
5. **Peer-to-Peer Native**: Perfect for your "pooling" concept
6. **Scalable**: Proven at massive scale
7. **Flexible**: Can do both P2P and hub-based models

### Architecture Recommendation

**Phase 1: Core OCPI Hub**
- Build OCPI 2.3.0 compliant hub platform
- Enable peer-to-peer connections between charging networks
- Implement all core OCPI modules:
  - Locations (charger information)
  - Sessions (charging sessions)
  - CDRs (billing records)
  - Tariffs (pricing)
  - Tokens (authentication)
  - Commands (start/stop charging)
  - Credentials (network authentication)

**Phase 2: Advanced Features**
- Implement OCPI 2.3 booking module for reservations
- Add smart charging capabilities
- Real-time availability updates
- Dynamic pricing support

**Phase 3: Hub Integrations** (Optional - for maximum reach)
- Integrate with Hubject using OICP or OCPI
- Integrate with GIREVE using eMIP or OCPI
- Integrate with e-clearing.net using OCHP
- This gives you access to their existing networks

**Phase 4: Consumer Applications**
- Mobile app for EV drivers
- Web portal for charge point owners
- Real-time charger discovery
- Booking and payment processing
- Reviews and ratings

### Technology Stack Recommendations

**Backend Platforms Supporting OCPI**:
- **AMPECO** - Full OCPI support, white-label platform
- **Driivz** - Enterprise platform with OCPI
- **ChargeLab** - North America focused
- **Virta** - European focused
- **Monta** - Modern platform with OCPI
- **Custom Build** - Using OCPI specification

**Open Source OCPI Implementations**:
- EVRoaming Foundation GitHub repositories
- Community-maintained libraries in various languages
- Reference implementations available

---

## 7. Implementation Roadmap

### Step 1: Platform Selection (Month 1-2)
- Choose between white-label platform vs custom build
- Evaluate OCPI-compliant platforms (AMPECO, Driivz, etc.)
- Set up development environment

### Step 2: Core OCPI Implementation (Month 3-6)
- Implement OCPI 2.3.0 server
- Build admin dashboard for network management
- Set up database for locations, sessions, CDRs
- Implement authentication and security

### Step 3: Network Partner Onboarding (Month 6-9)
- Recruit initial Charge Point Operators (CPOs)
- Establish bilateral roaming agreements
- Test OCPI connections with partners
- Set up pricing and billing

### Step 4: Driver Application (Month 9-12)
- Build mobile app (iOS/Android)
- Integrate charger discovery
- Implement booking and payment
- Launch beta testing

### Step 5: Scale & Integrate (Month 12+)
- Onboard more CPOs and eMSPs
- Consider integration with Hubject/GIREVE/e-clearing.net
- Expand geographic coverage
- Add advanced features (smart charging, V2G)

---

## 8. Business Model Considerations

### Revenue Streams

**Transaction Fees**:
- Small percentage of each charging session
- Competitive with existing hubs
- Transparent pricing for partners

**Membership Fees**:
- Monthly/annual fees for CPOs to join network
- Tiered pricing based on number of chargers
- Free tier for small operators to encourage adoption

**Premium Features**:
- Advanced analytics for CPOs
- Priority support
- Custom integration assistance
- White-label solutions

**Data & Insights**:
- Anonymized charging patterns
- Network utilization analytics
- Market insights for partners

### Competitive Advantages

✅ **Most Open Protocol**: OCPI is truly open and free
✅ **Peer-to-Peer Native**: Lower costs than hub-based models
✅ **No Lock-in**: Partners can connect to multiple platforms
✅ **Future-Proof**: Industry converging on OCPI
✅ **Scalable**: Designed for massive networks
✅ **Community-Driven**: Supported by EVRoaming Foundation

---

## 9. Key Technical Requirements

### Security
- TLS/SSL encryption for all communications
- OAuth 2.0 or token-based authentication
- Secure credential management
- Regular security audits
- GDPR compliance for user data

### Scalability
- Cloud-based architecture (AWS, Azure, GCP)
- Microservices design
- Load balancing
- Database sharding for high volume
- CDN for global performance

### Reliability
- 99.9%+ uptime SLA
- Redundant systems
- Automated failover
- Real-time monitoring
- Incident response procedures

### API Performance
- Sub-second response times
- Rate limiting to prevent abuse
- Caching for frequently accessed data
- Webhook support for real-time updates
- Comprehensive API documentation

---

## 10. Success Metrics

### Network Growth
- Number of connected CPO networks
- Total charging points in network
- Geographic coverage
- Number of active eMSPs

### Transaction Volume
- Daily/monthly charging sessions
- Revenue processed
- Cross-network roaming percentage
- Average session value

### User Satisfaction
- Driver app ratings
- CPO partner satisfaction scores
- Support ticket resolution time
- Platform uptime

### Market Position
- Market share in target regions
- Competitive pricing analysis
- Brand recognition
- Partnership quality

---

## FINAL RECOMMENDATION

### For Building the World's Largest EV Charging Roaming Hub:

**🏆 Primary Protocol: OCPI 2.3.0**

### Why OCPI Wins:

1. **✅ Correct Scope**: Network-to-network roaming (not charger management)
2. **✅ Most Widely Adopted**: Global open roaming standard
3. **✅ Industry Convergence**: Even Hubject and GIREVE adopting it
4. **✅ Truly Open**: Royalty-free, no vendor lock-in
5. **✅ Peer-to-Peer Native**: Perfect for your "pooling" concept
6. **✅ Future-Proof**: Active development, strong community
7. **✅ Proven at Scale**: Used by major platforms worldwide
8. **✅ Flexible**: Supports both P2P and hub-based models

### Optional Integrations:

Once your OCPI hub is established, you can integrate with existing major hubs to expand reach:
- **Hubject** (OICP or OCPI) - 1M+ chargers, 70+ countries
- **GIREVE** (eMIP or OCPI) - 659K chargers, 30 countries
- **e-clearing.net** (OCHP) - 445K chargers, Europe

### The Opportunity:

Build an **open, OCPI-based roaming hub** that:
- Enables direct peer-to-peer connections (lower costs)
- Bridges to existing major hubs (maximum reach)
- Supports individual charge point owners (true sharing economy)
- Uses the protocol the industry is converging on (future-proof)

**This is the path to creating the world's largest and most open EV charging network hub.**

---

## Resources

### Official Documentation
- **OCPI Specification**: https://github.com/ocpi/ocpi
- **EVRoaming Foundation**: https://evroaming.org/
- **Hubject**: https://www.hubject.com/
- **GIREVE**: https://www.gireve.com/
- **e-clearing.net**: https://www.ochp.eu/

### Community
- EVRoaming Foundation working groups
- OCPI GitHub discussions
- EV charging industry forums
- Regional EV associations

### Further Reading
- "Comparative analysis of standardized protocols for EV roaming" (evRoaming4EU project)
- OCPI Implementation Guides
- EV roaming best practices documentation
