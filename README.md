# RealDirect
A Market that never SLEEPS

Setup Instructions

Step 0: Download code to root directory

    -- Open a Terminal
    -- $ cd ~
    -- git clone https://github.com/sriram-kakarala/RealDirect.git
 

*Try the ChainCode Way*
Use this if you want to see the Demo based of Chaincode and Express. Step 0 is mandatory

Step 1. Generate New Crypto Assets. Run generate_crypto_assets.sh - If you want to generate new Genesis Blocks.

    -- This needs to be run once, or if you are cloning this repo, you can start with setup.sh
    -- Once you run, update the Private Key in docket-compose.yaml
    -- If you just want HAPPINESS go to Step 2
Step 2. Happiness

    -- Open a terminal and navigate to source path
    -- ./setup_and_initiate_network.sh
    -- Will fire up all the instances, create the channel, deploy and install the chaincode.
    -- Start Couch at http://localhost:5984/_utils/#/_all_dbs

Step 3. Node App Happiness
    -- Open a terminal and navigate to source path i.e uptil RealDirect    
    -- cd into api/rdapp i.e $ cd api/rdapp
    -- Fire away the instance $npm start
    -- Ahoy!!!
    
*Try the Composer Way*
Use this if you want to see the componser Demo. Step 0 is mandatory

Step 1. Open a Terminal, We assume that you have the basic fabric-dev-servers setup

    -- $ cd ~/fabric-dev-servers
    -- $ ./startFabric.sh
    -- $./createPeerAdminCard.sh

Step 2. Open a new Terminal Ctrl + Shift + T or Cmd + T

     -- $ cd ~/RealDirect/RealDirect-Composer
     -- $ composer network install --card PeerAdmin@hlfv1 --archiveFile bna-files/hedge-fund-network.bna
     -- $ composer network start --networkName hedge-fund-network1 --networkVersion 0.0.1 --networkAdmin admin --       networkAdminEnrollSecret adminpw --card PeerAdmin@hlfv1
     -- $ composer card import --file admin@hedge-fund-network1.card

Step 3. Setup the Rest Server

    -- $ composer-rest-server
    -- Use admin@hedge-fund-network1 as card name
    -- Wait for the Rest Server to start

Step 4. Open new Tab to start the UI App

    -- $ cd ~/RealDirect/RealDirect-Composer/app/ReaLDirect
    -- $ npm start
