## **Supply-chain-Management
In the traditional supply-chain system, stakeholders face credibility, transperancy and quality issues. Blockchain based supply-chain-management system improves efficiency, transparency, and traceability of supply chain processes.This application gives access to real-time, tamper-proof data about products to all the stakeholders. 



# **Chaincode**
This repository contains chaincode for supply chain management system using Hyperledger-Fabric 

# **Objective**:
This chaincode is developed for building a supply chain management application using blockchain. We have used Hyperledger fabric because of its enterprise grade capabilities. Since the blockchain is  transperant, immutable and secure decentralized system, it enables us to build an effective supply chain system. Stakeholders can keep a track of their assets in real time. It facilitaate efficient data sharing among all the stakeholders and enables them to build and maintain trust.

# **Functions in chaincode:**
- signIn
- createUser
- createProduct
- updateProduct
- toSupplier
- toTransporter
- sellToCustomer
- QueryProduct
- QueryAllProducts
- InitLedger

# **Changes**
Initially, chaincode was implemented using the ShimAPI. Chnaged it to ContractAPI. 
  

# **NOTE**:
There will be errors in the chaincode not all of the functions are updated to use the ContractAPI.
