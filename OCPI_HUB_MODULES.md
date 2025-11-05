# OCPI Hub - Required Modules & Implementation Order

## Understanding Hub vs CPO vs eMSP

### What is a Hub?
A **roaming hub** acts as an **intermediary** that connects multiple CPOs (Charge Point Operators) and eMSPs (e-Mobility Service Providers) together. The hub doesn't own charging stations itself.

**Hub Role:**
```
eMSP A ←→ HUB ←→ CPO X
eMSP B ←→ HUB ←→ CPO Y
eMSP C ←→ HUB ←→ CPO Z
```

The hub needs to support **BOTH sides**:
- **CPO-facing**: Receive locations, sessions, CDRs from CPOs
- **eMSP-facing**: Receive tokens from eMSPs, forward CDRs to eMSPs

---

## Required OCPI Modules for a Hub

### ✅ CRITICAL Modules (Must Have - Phase 1)

#### 1. **Credentials Module** (MUST BE FIRST!)
**Priority:** 🔴 CRITICAL - Start here!

**Purpose:** Partner registration and authentication

**What it does:**
- Partners (CPOs and eMSPs) register with the hub
- Exchange credentials and API tokens
- Establish trust relationship
- Version negotiation (OCPI 2.3.0)

**Hub must implement:**
- `POST /credentials` - Register new partner
- `GET /credentials` - Get partner credentials
- `PUT /credentials` - Update credentials
- `DELETE /credentials` - Unregister partner

**Why first?**
Without this, no partner can connect to your hub! This is the handshake.

**Data stored:**
```json
{
  "partner_id": "DE-ABC",
  "type": "CPO|EMSP",
  "token": "encrypted_token",
  "url": "https://partner.com/ocpi/2.3",
  "roles": [
    {
      "role": "CPO",
      "party_id": "ABC",
      "country_code": "DE"
    }
  ]
}
```

---

#### 2. **Locations Module** (Second Priority)
**Priority:** 🔴 CRITICAL - Implement after Credentials

**Purpose:** Manage charging station information

**Hub Role:**
- **FROM CPOs:** Receive location data (CPOs push their charging stations to hub)
- **TO eMSPs:** Provide location data (eMSPs query hub to find chargers)

**Hub must implement:**

**Receiving from CPOs (Sender Interface):**
- `PUT /locations/{location_id}` - CPO creates/updates location
- `PATCH /locations/{location_id}` - CPO partially updates location
- `PUT /locations/{location_id}/evses/{evse_id}` - CPO updates specific EVSE
- `PATCH /locations/{location_id}/evses/{evse_id}` - CPO updates EVSE status
- `PUT /locations/{location_id}/evses/{evse_id}/connectors/{connector_id}` - Update connector

**Providing to eMSPs (Receiver Interface):**
- `GET /locations` - eMSP lists all locations (with pagination, filtering)
- `GET /locations/{location_id}` - eMSP gets specific location

**Why second?**
eMSPs need to know where chargers are! This is the core of roaming.

**Key Features:**
- Geospatial queries (find chargers near me)
- Real-time availability updates
- Filter by connector type, power, etc.

---

#### 3. **Tokens Module** (Third Priority)
**Priority:** 🔴 CRITICAL - Needed for authentication

**Purpose:** Manage EV driver authentication tokens (RFID cards, app users)

**Hub Role:**
- **FROM eMSPs:** Receive token whitelist (eMSPs push their customer tokens)
- **TO CPOs:** Validate tokens (CPOs query hub to check if token is authorized)

**Hub must implement:**

**Receiving from eMSPs (Sender Interface):**
- `PUT /tokens/{country_code}/{party_id}/{token_uid}` - eMSP creates/updates token
- `PATCH /tokens/{country_code}/{party_id}/{token_uid}` - eMSP updates token

**Providing to CPOs (Receiver Interface):**
- `GET /tokens/{country_code}/{party_id}/{token_uid}` - CPO checks token validity
- `POST /tokens/{token_uid}/authorize` - CPO validates token in real-time

**Why third?**
Before starting a charging session, CPOs need to validate that the driver's token is authorized.

**Authorization Flow:**
```
1. Driver presents RFID card at charger (CPO network)
2. CPO asks hub: "Is this token valid?"
3. Hub checks which eMSP owns this token
4. Hub returns: "Yes, belongs to eMSP-X, authorized to charge"
5. CPO starts charging session
```

---

#### 4. **Sessions Module** (Fourth Priority)
**Priority:** 🟡 HIGH - Needed for active session tracking

**Purpose:** Track active (ongoing) charging sessions

**Hub Role:**
- **FROM CPOs:** Receive session updates (CPO reports session start/progress/end)
- **TO eMSPs:** Provide session status (eMSP can query their customer's active sessions)

**Hub must implement:**

**Receiving from CPOs (Sender Interface):**
- `PUT /sessions/{country_code}/{party_id}/{session_id}` - CPO creates/updates session
- `PATCH /sessions/{country_code}/{party_id}/{session_id}` - CPO updates session progress

**Providing to eMSPs (Receiver Interface):**
- `GET /sessions` - eMSP lists their active sessions
- `GET /sessions/{country_code}/{party_id}/{session_id}` - eMSP gets specific session

**Why fourth?**
Once charging starts, both CPO and eMSP need to track the session status.

**Session Lifecycle:**
```
1. CPO starts session → PUT /sessions (sends to hub)
2. Hub stores session and notifies eMSP (webhook)
3. Session updates (kWh, duration) → PATCH /sessions
4. Session ends → Final PUT /sessions with end_datetime
5. CDR generation triggered
```

---

#### 5. **CDRs Module** (Fifth Priority - CRITICAL for billing!)
**Priority:** 🔴 CRITICAL - Money depends on this!

**Purpose:** Charge Detail Records - billing information for completed sessions

**Hub Role:**
- **FROM CPOs:** Receive CDRs (CPO sends billing data after session ends)
- **TO eMSPs:** Forward CDRs (Hub sends CDR to eMSP who owns the token)

**Hub must implement:**

**Receiving from CPOs (Sender Interface):**
- `POST /cdrs` - CPO sends completed CDR
- `GET /cdrs/{cdr_id}` - CPO checks if CDR was received

**Providing to eMSPs (Receiver Interface):**
- `GET /cdrs` - eMSP lists their CDRs
- `GET /cdrs/{cdr_id}` - eMSP gets specific CDR

**Why fifth (but critical)?**
This is how billing happens! Without CDRs, nobody gets paid. Must be 100% reliable.

**CDR Flow:**
```
1. Session ends
2. CPO generates CDR with costs
3. CPO sends CDR to hub (POST /cdrs)
4. Hub validates and stores CDR
5. Hub forwards CDR to eMSP (who bills their customer)
6. Revenue sharing calculated
```

**Reliability is KEY:**
- Use Temporal workflow for reliable CDR processing
- Retry delivery to eMSP if it fails
- Store CDR forever (financial records!)

---

### ✅ IMPORTANT Modules (Should Have - Phase 2)

#### 6. **Tariffs Module**
**Priority:** 🟡 HIGH - Needed for pricing transparency

**Purpose:** Pricing information (how much does charging cost?)

**Hub Role:**
- **FROM CPOs:** Receive tariff structures
- **TO eMSPs:** Provide pricing info so drivers know costs upfront

**Hub must implement:**
- `PUT /tariffs/{country_code}/{party_id}/{tariff_id}` - CPO creates/updates tariff
- `GET /tariffs` - eMSP lists tariffs
- `GET /tariffs/{country_code}/{party_id}/{tariff_id}` - Get specific tariff

**Why important?**
Drivers want to know the price before charging. Also needed for CDR calculation.

---

#### 7. **Commands Module**
**Priority:** 🟡 HIGH - Enables remote control

**Purpose:** Remote operations (start/stop charging, unlock connector, reserve)

**Hub Role:**
- **FROM eMSPs:** Receive commands (eMSP wants to remote start/stop)
- **TO CPOs:** Forward commands (Hub tells CPO to execute command)

**Hub must implement:**

**Receiving from eMSPs:**
- `POST /commands/START_SESSION` - Start charging remotely
- `POST /commands/STOP_SESSION` - Stop charging remotely
- `POST /commands/UNLOCK_CONNECTOR` - Unlock connector
- `POST /commands/RESERVE_NOW` - Reserve a charger

**Providing to CPOs:**
- Forward commands via CPO's receiver interface
- Track command status (pending, accepted, rejected)

**Async Flow:**
```
1. eMSP sends command to hub
2. Hub immediately returns: "Command received" (202 Accepted)
3. Hub forwards command to CPO
4. CPO executes and responds
5. Hub updates command status
6. Hub notifies eMSP via webhook
```

**Why important?**
Enables app-based charging (no RFID card needed).

---

### ⚪ OPTIONAL Modules (Nice to Have - Phase 3)

#### 8. **Charging Profiles Module** (OCPI 2.3)
**Priority:** 🟢 OPTIONAL - For smart charging

**Purpose:** Set charging speed/schedule (smart charging)

**Use Case:**
- Charge slower during peak hours
- Charge faster when renewable energy available
- Schedule charging for night-time

---

#### 9. **ChargingSessions Module** (OCPI 2.3)
**Priority:** 🟢 OPTIONAL - Alternative to Sessions

**Note:** Similar to Sessions but with some differences. Can skip if implementing Sessions.

---

#### 10. **Reservations Module** (OCPI 2.3 Booking)
**Priority:** 🟢 OPTIONAL - For advance booking

**Purpose:** Reserve a charger in advance

**Use Case:**
- Driver books charger for specific time
- Charger is held for them
- Guarantees availability

---

## Implementation Order & Timeline

### **Phase 1: Foundation (Weeks 1-4)** ← START HERE

**Week 1-2: Project Setup + Credentials**
- ✅ Set up Go project structure
- ✅ MongoDB + Redis + Temporal setup
- ✅ Basic HTTP server with Gin
- ✅ **Credentials Module** - Partner registration
- ✅ Token authentication middleware
- ✅ Partner management (create, update, delete)

**Deliverable:** Partners can register and get API tokens

---

**Week 3-4: Locations Module**
- ✅ **Locations Module** - Full CRUD
- ✅ MongoDB geospatial indexes
- ✅ Sender interface (CPOs can push locations)
- ✅ Receiver interface (eMSPs can query locations)
- ✅ Real-time availability updates
- ✅ Search by distance, connector type, power

**Deliverable:** CPOs can add charging stations, eMSPs can search for chargers

---

### **Phase 2: Core Roaming (Weeks 5-8)**

**Week 5-6: Tokens Module**
- ✅ **Tokens Module** - Token management
- ✅ eMSPs can register their customer tokens
- ✅ CPOs can validate tokens (authorization)
- ✅ Real-time token authorization endpoint
- ✅ Token whitelist management

**Deliverable:** CPOs can validate driver tokens before starting sessions

---

**Week 7: Sessions Module**
- ✅ **Sessions Module** - Active session tracking
- ✅ CPOs report session start/progress/end
- ✅ eMSPs can query active sessions
- ✅ Real-time session updates
- ✅ Webhook notifications for session events

**Deliverable:** Live charging sessions are tracked

---

**Week 8: CDRs Module + Temporal Workflows**
- ✅ **CDRs Module** - Billing records
- ✅ CPOs send CDRs after session completion
- ✅ Hub forwards CDRs to eMSPs
- ✅ **Temporal workflows:**
  - `ProcessCDRWorkflow` - Reliable CDR delivery
  - `SessionManagementWorkflow` - Session lifecycle
  - `WebhookDeliveryWorkflow` - Retry failed webhooks

**Deliverable:** Complete billing flow works reliably

---

### **Phase 3: Advanced Features (Weeks 9-12)**

**Week 9-10: Tariffs + Commands**
- ✅ **Tariffs Module** - Pricing information
- ✅ **Commands Module** - Remote control
- ✅ Handle async command flow
- ✅ Command status tracking

**Deliverable:** Drivers can see prices and remote start/stop charging

---

**Week 11-12: Testing + Optimization**
- ✅ Integration testing with test partners
- ✅ Load testing
- ✅ Performance optimization
- ✅ Security audit
- ✅ Documentation

**Deliverable:** Production-ready hub

---

## Module Dependencies

```
Credentials (NO dependencies - START HERE!)
    ↓
Locations (depends on: Credentials)
    ↓
Tokens (depends on: Credentials)
    ↓
Sessions (depends on: Credentials, Locations, Tokens)
    ↓
CDRs (depends on: Sessions, Tariffs)
    ↓
Tariffs (depends on: Credentials)
    ↓
Commands (depends on: Credentials, Sessions)
```

---

## Which Modules Do We Build First?

### **RECOMMENDED START:**

**Sprint 1 (Week 1-2):**
```
1. Credentials Module (MUST BE FIRST!)
   - Partner registration
   - Token authentication
   - Basic partner CRUD
```

**Sprint 2 (Week 3-4):**
```
2. Locations Module
   - CPO can push locations
   - eMSP can search locations
   - Geospatial queries
```

**Sprint 3 (Week 5-6):**
```
3. Tokens Module
   - eMSP registers tokens
   - CPO validates tokens
   - Real-time authorization
```

**Sprint 4 (Week 7-8):**
```
4. Sessions Module
5. CDRs Module
   - Complete charging flow
   - Billing works end-to-end
```

After this, you have a **working roaming hub**! eMSPs and CPOs can do basic roaming.

---

## My Recommendation

**START WITH:**
1. **Credentials Module** (Week 1-2)

This is the foundation. Without it, partners can't even connect to your hub.

**What we'll build:**
- Partner registration endpoint
- Token generation and storage
- Authentication middleware
- Partner management API

**Once Credentials is done, we move to Locations.**

---

## Next Steps

**Should I start building the Credentials Module?**

I'll create:
1. Go project structure
2. MongoDB connection and partner schema
3. Credentials API endpoints
4. Token authentication
5. Partner registration flow

**Ready to start?** 🚀
