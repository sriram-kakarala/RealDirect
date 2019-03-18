#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

docker-compose -f docker-compose.yaml down

docker-compose -f docker-compose.yaml up -d ca.realdirect.com orderer.realdirect.com peer0.nsdl.realdirect.com couchdb cli

# wait for Hyperledger Fabric to start
# incase of errors when running later commands, issue export FABRIC_START_TIMEOUT=<larger number>
export FABRIC_START_TIMEOUT=10
#echo ${FABRIC_START_TIMEOUT}
sleep ${FABRIC_START_TIMEOUT}

# Create the channel
docker exec -e "CORE_PEER_LOCALMSPID=RealDirectNSDLMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@nsdl.realdirect.com/msp" peer0.nsdl.realdirect.com peer channel create -o orderer.realdirect.com:7050 -c realdirectchannel -f /etc/hyperledger/configtx/channel.tx
# Join peer0.org1.example.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=RealDirectNSDLMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@nsdl.realdirect.com/msp"  peer0.nsdl.realdirect.com peer channel join -b realdirectchannel.block

# Install
docker exec -e "CORE_PEER_LOCALMSPID=RealDirectNSDLMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@nsdl.realdirect.com/msp" cli peer chaincode install -n realdirect3 -v 1.0 -p "github.com/" -l "golang"

# Initiate
docker exec -e "CORE_PEER_LOCALMSPID=RealDirectNSDLMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@nsdl.realdirect.com/msp" cli peer chaincode instantiate -o orderer.realdirect.com:7050 -C realdirectchannel -n realdirect3 -l "golang" -v 1.0 -c '{"Args":[]}' -P "OR ('RealDirectNSDLMSP.member','Org2MSP.member')"