import {Asset} from './org.hyperledger.composer.system';
import {Participant} from './org.hyperledger.composer.system';
import {Transaction} from './org.hyperledger.composer.system';
import {Event} from './org.hyperledger.composer.system';
// export namespace org.acme.mynetwork{
   export class Lot extends Asset {
      lotId: string;
      securityName: string;
      quantity: number;
      price: number;
      owner: Client;
   }
   export class Client extends Participant {
      clientId: string;
      description: string;
   }
   export class Custodian extends Participant {
      custodianId: string;
      description: string;
   }
   export class Trader extends Participant {
      traderId: string;
      name: string;
   }
   export class Trade extends Transaction {
      trader: Trader;
      client: Client;
      lot: Lot;
   }
   export class NewTradeEvent extends Event {
      lot: Lot;
   }
// }
