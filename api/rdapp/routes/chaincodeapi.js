var express = require('express');
var chainCodeRouter = express.Router();

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve("", 'public/fabric', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

chainCodeRouter.get('/signup', function(req, res, next) {
    res.render('signup', { title: 'RealDirect - Sign Up' });
  });
  
chainCodeRouter.get('/signin', function(req, res, next) {
    res.render('signin', { title: 'RealDirect - Sign In' });
  });

chainCodeRouter.get('/dashboard', function(req, res, next) {
    var email = req.query.email;
    res.render('dashboard', { title: 'Welcome to RealDirect', username:  email});
  });

chainCodeRouter.post('/signup', function(req, res, next) {
  console.log(req.body)
  signup(req.body.name, req.body.email, req.body.password, res)
});

async function signup(name, email, password, res) {
    try {
  
          // Create a new file system based wallet for managing identities.
          const walletPath = path.resolve("", 'public/fabric', 'wallet');
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
  
          
          console.log('attempting the transaction');
          await contract.submitTransaction('inituser', name, email, password);
          console.log('Transaction has been submitted');
          // Disconnect from the gateway.
          await gateway.disconnect();
          res.render('signin', { title: 'RealDirect - Sign In', msg: 'Sign Up Success' });
      } catch (error) {
          console.error(`Failed to submit transaction: ${error}`);
          if(error.message.includes('300')) {
            res.render('signup', { title: 'RealDirect - Sign Up', err: 'User already exists' });
          } else if(error.message.includes('301')) {
            res.render('signup', { title: 'RealDirect - Sign Up', err: 'Server is Busy, please try again' });
          } else if(error.message.includes('302')){
            res.render('signup', { title: 'RealDirect - Sign Up', err: 'Server is Busy, please try again' });
          }
      }
  }

chainCodeRouter.post('/signin', function(req, res, next) {
    console.log(req.body)
    signin(req.body.email, req.body.password, res)
});

async function signin(email, password, res) {
    try {
  
        // Create a new file system based wallet for managing identities.
        const walletPath = path.resolve("", 'public/fabric', 'wallet');
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

        
        console.log('attempting the transaction');
        await contract.submitTransaction('signinuser', email, password);
        console.log('Transaction has been submitted');
        // Disconnect from the gateway.
        await gateway.disconnect();
        res.redirect('/dashboard?email=' + email);
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        if(error.message.includes('300')) {
            res.render('signin', { title: 'RealDirect - Sign In', err: 'Please check your Email Id/Password' });
        } else if(error.message.includes('301')) {
            res.render('signin', { title: 'RealDirect - Sign In', err: 'Server is Busy, please try again' });
        } else if(error.message.includes('302')){
            res.render('signin', { title: 'RealDirect - Sign In', err: 'Please check your Email Id/Password' });
        }
    }
}

module.exports = chainCodeRouter;
