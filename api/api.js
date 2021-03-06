/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);
const argv = require('yargs').argv

async function main() {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('realdirectchannel');

        // Get the contract from the network.
        const contract = network.getContract('realdirect');

        // Submit the specified transaction.
        if(argv.target == 'reguser') {
            console.log('attempting the transaction');
            await contract.submitTransaction('inituser', argv.name, argv.email, argv.password);
            console.log('Transaction has been submitted');
        }
        else if(argv.target == 'initasset') {
            console.log('attempting the transaction');
            await contract.submitTransaction('initasset', argv.name, "" + argv.quantity, argv.owner, "" + argv.price);
            console.log('Transaction has been submitted');
        } else if(argv.target == 'readasset'){
            const result = await contract.evaluateTransaction('readasset', argv.name);
            console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        } else if(argv.target == 'delete'){
            await contract.submitTransaction('delete', argv.name);
            console.log('Transaction has been submitted');
        } else if(argv.target == 'transferasset'){
            await contract.submitTransaction('transferasset', argv.name, argv.owner);
            console.log('Transaction has been submitted');
        } else {
            console.log('Invalid Transaction!! Invoke with --target=[createCar|changeCarOwner]');
        }

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
