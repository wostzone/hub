import {reactive} from "vue";

import {AccountRecord} from "@/data/AccountStore";
import MqttClient from "@/data/MqttClient";
import DirectoryClient from "@/data/DirectoryClient";
import AuthClient from "@/data/AuthClient";

// Account connection status
export interface IConnectionStatus {
  readonly accountID: string        // ID of the account
  connected: boolean        // authenticated and at least one service connected
  authenticated: boolean    // authentication was successful
  directory: boolean        // the directory is obtained
  messaging: boolean        // message bus connection is established
  statusMessage: string     // human description of connection status
}
export type AccountConnection = {
    readonly accountID: string
    name: string
    authClient: AuthClient|null
    mqttClient: MqttClient|null
    dirClient: DirectoryClient|null
    state: IConnectionStatus
}

// Manage account connections to MQTT and Directory services
export class ConnectionManager {
    // accounts that are watched
    // active mqtt broker clients by account ID
    private connections: Map<string, AccountConnection>
    private started: boolean
    // active connection state
    private status: IConnectionStatus

    // Create a new connection manager
    constructor() {
        this.connections = new Map<string, AccountConnection>()
        this.status = reactive(<IConnectionStatus>{
            accountID:"",
            connected:false,
            authenticated:false,
            directory:false,
          messaging: false,
          statusMessage: "Not connected"
        })
        this.started = false
    }

    // Active connection status
    get connectionStatus(): IConnectionStatus {
        return this.status
    }

    // Nr of authenticated connections
    get connectionCount(): number {
        let count = 0
        this.connections.forEach((connection: AccountConnection) => {
            if (connection.authClient && connection.authClient.IsConnected()) {
                count++
            }
        })
        return count
    }

    // Connect or reconnect.
    // This provides both a promise for the initial connection result and a callback that is invoked each time
    // the connection changes. If a connection fails, then a new connection attempt will be made periodically.
    //
    //  @param account is the account to connect to authentication, mqtt and directory services
    //  @param onConnectChanged optional callback when the connection status changes
    // This returns a promise for async operation of the first connection attempt.
    async Connect(account: AccountRecord,
                  onConnectChanged: ((record:AccountRecord, connected:boolean, error:Error|null)=>void)|undefined =undefined ) {
        this.Disconnect(account.id);

        let ac = this.GetAccountConnection(account)
        ac.authClient = new AuthClient(account.address,account.authPort)
        ac.mqttClient = new MqttClient(
            account.id,
            account.address,
            account.mqttPort,
            this.handleMqttConnected.bind(this),
            this.handleMqttDisconnected.bind(this),
            this.handleMqttMessage.bind(this))

        ac.dirClient = new DirectoryClient(account.address, account.directoryPort)

        console.log("ConnectionManager.Connect: Connecting to", account.address, "as", account.loginName)
        this.connections.set(ac.accountID, ac)
        // after login, connect to mqtt and directory client using the access token
        let password = "user1" // FIXME - for testing

        return ac.authClient.ConnectWithLoginID(account.loginName, password)
            .then((accessToken: string) => {
                console.log("ConnectionManager.Connect: Authentication successful. Connecting to mqtt and directory")
                ac.state.authenticated = true
              ac.state.statusMessage = "Authentication successful"
                this.status.authenticated = true
                // this.status.connected = true

                // FIXME support login to mqtt with access token
                if (ac.mqttClient) {
                    ac.mqttClient.Connect(account.loginName, password)
                }
                if (ac.dirClient) {
                    ac.dirClient.Connect(accessToken)
                }
                if (onConnectChanged) {
                    onConnectChanged(account, true, null)
                }
            })
            .catch((err: Error) => {
                ac.state.authenticated = false
              ac.state.statusMessage = "Failed to authenticate: " + err.message
              console.error("ConnectionManager.Connect: failed to authenticate: ", err)
                if (onConnectChanged) {
                    onConnectChanged(account, false, err)
                }
                throw(err.message)
            });
    }

    // Disable the connection
    Disconnect(accountId:string) {
        let connection = this.connections.get(accountId)
        if (connection) {
            console.log("AccountManager.Disconnect: Disconnecting account:", connection.name)
            // if (connection.authClient) {
            //     connection.authClient.Disconnect()
            // }
            if (connection.dirClient) {
                connection.dirClient.Disconnect();
            }
            if (connection.mqttClient) {
                connection.mqttClient.Disconnect();
            }
        }
    }

    // Close all existing connections (modify the map as we go)
    // DisconnectAll() {
    //     for(let key of Array.from( this.connections.keys()) ) {
    //         this.Disconnect(key)
    //     }
    // }


    // Re-connect all enabled accounts
    // async ConnectAll(accounts: Array<AccountRecord>,
    //                  onConnectChanged: (record:AccountRecord, connected:boolean, err:Error|null) => void) {
    //     let p: Promise<Awaited<void[][number]>>;
    //
    //     p = Promise.any(
    //         accounts.map((item: AccountRecord) => {
    //             if (item.enabled) {
    //                 this.Connect(item, onConnectChanged)
    //             }
    //         })
    //     );
    //     return p
    // }


    // Get the connection of an account or create on if it doesn't exist
    // If an account connection isn't yet known, it will be created
    //
    // account to get the status of
    // return the connection instance
    protected GetAccountConnection(account: AccountRecord): AccountConnection {
        let connection = this.connections.get(account.id)
        if (!connection) {
            connection = {
                name: account.name,
                accountID: account.id,
                authClient: null,
                mqttClient: null,
                dirClient: null,
              state: reactive(<IConnectionStatus>{
                statusMessage: "Not connected"
              })
            }
            this.connections.set(account.id, connection)
        }
        return connection
    }

    // Get the reactive connection status of an account.
    // The result is reactive and can be used directly in a UI
    // This returns an empty status object if the account is not known
    // If no accountID is specified, then return the currently active account status
    //
    // accountID to get or "" for the active account status
    // return the connection status object
    GetConnectionStatus(account:AccountRecord): IConnectionStatus {
        let connection = this.GetAccountConnection(account)
        return connection.state
    }

    // Handle an incoming MQTT message
    handleMqttMessage(_accountID:string, topic: string, payload:Buffer, _retain: boolean): void {
        console.log("handleMqttMessage. topic:",topic, "Message size:", payload.length)
    }

    // track the MQTT account connection status
    handleMqttConnected(accountID:string) {
        let connection = this.connections.get(accountID)
        if (connection) {
            connection.state.messaging = true
            connection.state.connected = true
          connection.state.statusMessage = "Connected to message bus"
        }
        this.updateStatus()
    }
    // track the MQTT account connection status
    handleMqttDisconnected(accountID:string) {
        let connection = this.connections.get(accountID)
        if (connection) {
            connection.state.messaging = false
            connection.state.connected = false
          connection.state.statusMessage = "Disconnected from messaging"
        }
        this.updateStatus()
    }

    // update the aggregate connection status of all accounts
    updateStatus() {
      let messaging = false, connected = false, directory = false, authenticated = false;
        this.connections.forEach( c => {
            messaging = messaging || c.state.messaging
            connected = connected || c.state.connected
            directory = directory || c.state.directory
          authenticated = authenticated || c.state.authenticated
        })
        this.status.connected = connected
        this.status.messaging = messaging
        this.status.directory = directory
      this.status.authenticated = authenticated

      let newMessage = "Not connected"
      if (authenticated) {
        newMessage = "The user is authenticated"
        if (directory) {
          newMessage += ", the directory of Things is retrieved"
        }
        if (messaging) {
          newMessage += " and message bus connection is established"
        }
        if (!directory && !messaging) {
          newMessage = "Authenticated but not connected"
        }
      } else {
        newMessage = "Not authenticated"
      }
      this.status.statusMessage = newMessage

    }
}

// the global connection manager
const cm = new ConnectionManager()
export default cm