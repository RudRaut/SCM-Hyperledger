# **Chaincode**

This repository contains chaincode for supply chain management system using Hyperledger-Fabric 

# **Objective**:
This chaincode is developed for building a supply chain management application using blockchain. We have used Hyperledger fabric because of its enterprise grade capabilities. Since the blockchain is  transperant, immutable and secure decentralized system, it enables us to build an effective supply chain system. Stakeholders can keep a track of their assets in real time. It facilitaate efficient data sharing among all the stakeholders and enables them to build and maintain trust.

# **Functions in chaincode:**
- createUser
- signIn
- createProduct
- updateProduct
- toSupplier
- toTransporter
- sellToCustomer
- QueryAsset
- QueryAll
- orderProduct
- Init
- Invoke

# **Changes**
Initially, chaincode was implemented using the ShimAPI. Chnaged it to ContractAPI. 
  
# **Functions Updated**:
InitLedger and Sign along with some helper functions

# **NOTE**:
There will be errors in the chaincode not all of the functions are updated to use the ContractAPI.
