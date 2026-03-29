export namespace config {
	
	export class AppSettings {
	    mode: string;
	    socks5_port: number;
	    http_port: number;
	    api_port: number;
	    auto_connect: boolean;
	    last_server_id?: string;

	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.socks5_port = source["socks5_port"];
	        this.http_port = source["http_port"];
	        this.api_port = source["api_port"];
	        this.auto_connect = source["auto_connect"];
	        this.last_server_id = source["last_server_id"];
	    }
	}
	export class ConnectionStatus {
	    connected: boolean;
	    server_id?: string;
	    server_name?: string;
	    mode: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connected = source["connected"];
	        this.server_id = source["server_id"];
	        this.server_name = source["server_name"];
	        this.mode = source["mode"];
	    }
	}
	export class ServerConfig {
	    id: string;
	    name: string;
	    address: string;
	    port: number;
	    uuid: string;
	    transport: string;
	    tls: boolean;
	    sni?: string;
	    path?: string;
	    flow?: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.port = source["port"];
	        this.uuid = source["uuid"];
	        this.transport = source["transport"];
	        this.tls = source["tls"];
	        this.sni = source["sni"];
	        this.path = source["path"];
	        this.flow = source["flow"];
	    }
	}

}

