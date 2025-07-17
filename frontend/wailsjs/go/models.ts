export namespace main {
	
	export class HistoryItem {
	    id: number;
	    method: string;
	    url: string;
	    fullCommand: string;
	    statusCode: number;
	    responseBody?: string;
	    responseHeaders?: string;
	    responseType?: string;
	    requestSize: number;
	    responseSize: number;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new HistoryItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.method = source["method"];
	        this.url = source["url"];
	        this.fullCommand = source["fullCommand"];
	        this.statusCode = source["statusCode"];
	        this.responseBody = source["responseBody"];
	        this.responseHeaders = source["responseHeaders"];
	        this.responseType = source["responseType"];
	        this.requestSize = source["requestSize"];
	        this.responseSize = source["responseSize"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RequestData {
	    method: string;
	    url: string;
	    headers: Record<string, string>;
	    body: string;
	
	    static createFrom(source: any = {}) {
	        return new RequestData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.method = source["method"];
	        this.url = source["url"];
	        this.headers = source["headers"];
	        this.body = source["body"];
	    }
	}
	export class ResponseTiming {
	    totalTime: number;
	    connectTime: number;
	    transferTime: number;
	
	    static createFrom(source: any = {}) {
	        return new ResponseTiming(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalTime = source["totalTime"];
	        this.connectTime = source["connectTime"];
	        this.transferTime = source["transferTime"];
	    }
	}
	export class ResponseData {
	    statusCode: number;
	    headers: Record<string, string>;
	    body: string;
	    formattedBody?: string;
	    responseType: string;
	    timing: ResponseTiming;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ResponseData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.statusCode = source["statusCode"];
	        this.headers = source["headers"];
	        this.body = source["body"];
	        this.formattedBody = source["formattedBody"];
	        this.responseType = source["responseType"];
	        this.timing = this.convertValues(source["timing"], ResponseTiming);
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

