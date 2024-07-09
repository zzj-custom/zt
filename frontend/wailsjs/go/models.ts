export namespace app {
	
	export class LoginRequest {
	    email: string;
	    captcha: number;
	
	    static createFrom(source: any = {}) {
	        return new LoginRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.email = source["email"];
	        this.captcha = source["captcha"];
	    }
	}
	export class ProcessRequest {
	    flag: string;
	    pType: number;
	    outPath: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.flag = source["flag"];
	        this.pType = source["pType"];
	        this.outPath = source["outPath"];
	    }
	}

}

export namespace response {
	
	export class Reply {
	    code: number;
	    msg: string;
	    result?: any;
	
	    static createFrom(source: any = {}) {
	        return new Reply(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.result = source["result"];
	    }
	}

}

