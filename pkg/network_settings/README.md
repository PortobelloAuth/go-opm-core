# Network Settings

Network settings are the key concept allowing gradual adoption of Open Person
Matching, new versions, and other supported technologies for implementing a
privacy preserving records exchange. This is intended to make Open Person
Matching networks inherently upgradable. Requestors, Responders, and
Facilitators all may participate on multiple networks at once, enabling them to
start with the most recently implemented settings and fall back to networks
using older, more broadly supported standards.

For example, a Facilitator could participate in three Open Person Matching
networks with these (conceptual) settings:

1. TEFCA-compatible (not privacy preserving):
   - Type: tefca2026
   - Bridge to XPDQ / Carequality query through QHIN and existing frameworks
     Demographics and credentials passed unmodified without matching first
   - Third-party verified patient consent required for use
   - Established expiration date
   - IAL2 identity for IAS / Federated FHIR
   - Facilitators without their own client credentials for a QHIN will proxy
     requests to one that does
   - Requestors, Responders, and Facilitators authenticate with appropriate KERI
     AID
   - No salts
   - No hashing
2. CMS-Aligned-2026-07
   - Type: hash-routed-challenge
   - CMS Aligned Network matching rules v3.2.2
   - CMS Aligned Network matching rule derived routing keys
   - 24 character prefix
   - PBKDF2-XR2 with 8-bit challenge
   - IAL2 identity with CMS Aligned consent artifact (or SMART Permission
     ticket)
   - Requestors, Responders, and Facilitators authenticate with appropriate KERI
     AID (details TBD)
   - Paxos / raft negotiated salts
   - Asymmetric encryption: HKDF
   - Symmetric Encryption: JWE
3. CMS-Aligned-SEDI
   - Type: hash-routed-bloom-filter
   - CMS Aligned Network matching rules v4
   - CMS Aligned Network matching rule derived routing keys
   - 24 character prefix
   - SEDI / IAL2 identity with CMS Aligned consent in ACDC artifact
   - Requestors, Responders, and Facilitators authenticate with appropriate KERI
     AID
   - Paxos / raft negotiated salts
   - Asymmetric encryption: HKDF (or a newer post-quantum safe algorithm)
   - Symmetric Encryption: JWE
