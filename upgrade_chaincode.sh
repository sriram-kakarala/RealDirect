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

# Install
docker exec -e "CORE_PEER_LOCALMSPID=RealDirectNSDLMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@nsdl.realdirect.com/msp" cli peer chaincode install -n realdirect -v $1 -p "github.com/" -l "golang"

docker exec -e "CORE_PEER_LOCALMSPID=RealDirectNSDLMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@nsdl.realdirect.com/msp" cli peer chaincode upgrade -n realdirect -C realdirectchannel -c '{"Args":[]}' -v $1 -p "github.com/" -l "golang"

