// Constants for use by applications
import {TLSSocket} from "tls";

const DefaultServiceName = "thingdir"
const DefaultPort = 8886

// paths with REST commands
const RouteThings  = "/things"           // list or query path
const RouteThingID = "/things/{thingID}" // for methods get, post, patch, delete

// query parameters
const ParamOffset = "offset"
const ParamLimit  = "limit"
const ParamQuery  = "queryparams"

const DefaultLimit = 100
const MaxLimit = 1000

// Thing description document
interface ThingTD {
    id: string
    props: Map<string, object>
}

// Client for connecting to a Hub Directory service
export default class DirectoryClient {
    private hostport: string = ""
    private tlsClient: TLSSocket|null = null

    constructor() {
        this.hostport = ""
        this.tlsClient = null
    }

    /* Close the connection to the directory server
     */
    async Close() {
    }

    /* ConnectWithClientCert opens the connection to the directory server using a client certificate for authentication
     *  @param clientCertFile  client certificate to authenticate the client with the broker
     *  @param clientKeyFile   client key to authenticate the client with the broker
     */
    // public ConnectWithClientCert(tlsClientCert: tls.Certificate):void {
    // }

    /* ConnectWithLoginID open the connection to the directory server using a login ID and password for authentication
     */
    // ConnectWithLoginID(loginID: string, password: string): Error {
    //     return null
    // }

    // Connect open the connection to the directory server using an access token
    // @param address
    // @param port, 0 or undefined uses the default port
    // @param accessToken
    async Connect(address: string, port: number|undefined, accessToken: string) {
        if (port == 0 || port == undefined) {
            port = DefaultPort
        }
    }

    /* Delete a TD
     * @param id of the Thing Description document
     */
    // Delete(id: string) :void {
    // }

    // GetTD the ThingTD with the given ID
    //  id is the ThingID whose ThingTD to get
    // GetTD(id: string): ThingTD|undefined {
    //     return undefined
    // }

    /* ListTDs
     * Returns a list of TDs starting at the offset. The result is limited to the nr of records provided
     * with the limit parameter. The server can choose to apply its own limit, in which case the lowest
     * value is used.
     * @param offset of the list to query from
     * @param limit result to nr of TDs. Use 0 for default.
     */
    async ListTDs(offset: number, limit: number): Promise<Array<ThingTD>> {
        return Array<ThingTD>()
    }

    /* PatchTD changes a TD with the attributes of the given TD
     */
    async PatchTD(id: string, td: ThingTD) {
    }

    /* QueryTDs with the given JSONPATH expression
     * Returns a list of TDs matching the query, starting at the offset. The result is limited to the
     * nr of records provided with the limit parameter. The server can choose to apply its own limit,
     * in which case the lowest value is used.
     * @param jsonPath with the query expression
     * @param offset is the start index of the list to query from
     * @param limit limits the result to nr of TDs. Use 0 for default.
     */
    async QueryTDs(jsonpath: string, offset:number, limit: number): Promise<Array<ThingTD>> {
        return Array<ThingTD>()
    }

    /* UpdateTD fully replaces the TD with the given ID, eg create/update
     */
    async UpdateTD(id: string, td: ThingTD)  {
    }
}

