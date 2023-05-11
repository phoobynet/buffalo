export namespace alpaca {
	
	export class Asset {
	    id: string;
	    class: string;
	    exchange: string;
	    symbol: string;
	    name: string;
	    status: string;
	    tradable: boolean;
	    marginable: boolean;
	    shortable: boolean;
	    easy_to_borrow: boolean;
	    fractionable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Asset(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.class = source["class"];
	        this.exchange = source["exchange"];
	        this.symbol = source["symbol"];
	        this.name = source["name"];
	        this.status = source["status"];
	        this.tradable = source["tradable"];
	        this.marginable = source["marginable"];
	        this.shortable = source["shortable"];
	        this.easy_to_borrow = source["easy_to_borrow"];
	        this.fractionable = source["fractionable"];
	    }
	}

}

export namespace calendar {
	
	export class Calendar {
	    date: string;
	    // Go type: time
	    open: any;
	    // Go type: time
	    close: any;
	    // Go type: time
	    sessionOpen: any;
	    // Go type: time
	    sessionClose: any;
	
	    static createFrom(source: any = {}) {
	        return new Calendar(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.open = this.convertValues(source["open"], null);
	        this.close = this.convertValues(source["close"], null);
	        this.sessionOpen = this.convertValues(source["sessionOpen"], null);
	        this.sessionClose = this.convertValues(source["sessionClose"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

export namespace marketdata {
	
	export class Bar {
	    // Go type: time
	    t: any;
	    o: number;
	    h: number;
	    l: number;
	    c: number;
	    v: number;
	    n: number;
	    vw: number;
	
	    static createFrom(source: any = {}) {
	        return new Bar(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.t = this.convertValues(source["t"], null);
	        this.o = source["o"];
	        this.h = source["h"];
	        this.l = source["l"];
	        this.c = source["c"];
	        this.v = source["v"];
	        this.n = source["n"];
	        this.vw = source["vw"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class Quote {
	    // Go type: time
	    t: any;
	    bp: number;
	    bs: number;
	    bx: string;
	    ap: number;
	    as: number;
	    ax: string;
	    c: string[];
	    z: string;
	
	    static createFrom(source: any = {}) {
	        return new Quote(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.t = this.convertValues(source["t"], null);
	        this.bp = source["bp"];
	        this.bs = source["bs"];
	        this.bx = source["bx"];
	        this.ap = source["ap"];
	        this.as = source["as"];
	        this.ax = source["ax"];
	        this.c = source["c"];
	        this.z = source["z"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class Trade {
	    // Go type: time
	    t: any;
	    p: number;
	    s: number;
	    x: string;
	    i: number;
	    c: string[];
	    z: string;
	    u: string;
	
	    static createFrom(source: any = {}) {
	        return new Trade(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.t = this.convertValues(source["t"], null);
	        this.p = source["p"];
	        this.s = source["s"];
	        this.x = source["x"];
	        this.i = source["i"];
	        this.c = source["c"];
	        this.z = source["z"];
	        this.u = source["u"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class Snapshot {
	    latestTrade?: Trade;
	    latestQuote?: Quote;
	    minuteBar?: Bar;
	    dailyBar?: Bar;
	    prevDailyBar?: Bar;
	
	    static createFrom(source: any = {}) {
	        return new Snapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.latestTrade = this.convertValues(source["latestTrade"], Trade);
	        this.latestQuote = this.convertValues(source["latestQuote"], Quote);
	        this.minuteBar = this.convertValues(source["minuteBar"], Bar);
	        this.dailyBar = this.convertValues(source["dailyBar"], Bar);
	        this.prevDailyBar = this.convertValues(source["prevDailyBar"], Bar);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

