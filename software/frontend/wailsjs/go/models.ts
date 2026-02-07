export namespace application {
	
	export class Result_bool_ {
	    data: boolean;
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result_bool_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	}
	export class Result_string_ {
	    data: string;
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result_string_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	}

}

