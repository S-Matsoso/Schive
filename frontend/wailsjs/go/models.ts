export namespace models {
	
	export class File {
	    id: string;
	    moduleId: string;
	    assignmentId?: string;
	    name: string;
	    path: string;
	    fileType: string;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.moduleId = source["moduleId"];
	        this.assignmentId = source["assignmentId"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.fileType = source["fileType"];
	    }
	}
	export class Assignment {
	    id: string;
	    moduleId: string;
	    assignmentNumber: number;
	    // Go type: time
	    dueDate?: any;
	    files?: File[];
	
	    static createFrom(source: any = {}) {
	        return new Assignment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.moduleId = source["moduleId"];
	        this.assignmentNumber = source["assignmentNumber"];
	        this.dueDate = this.convertValues(source["dueDate"], null);
	        this.files = this.convertValues(source["files"], File);
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
	export class ClassSlot {
	    id: string;
	    moduleId: string;
	    "day;; not null": string;
	    startHour: number;
	    lengthHrs: number;
	    room: string;
	
	    static createFrom(source: any = {}) {
	        return new ClassSlot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.moduleId = source["moduleId"];
	        this["day;; not null"] = source["day;; not null"];
	        this.startHour = source["startHour"];
	        this.lengthHrs = source["lengthHrs"];
	        this.room = source["room"];
	    }
	}
	
	export class Module {
	    id: string;
	    code: string;
	    name: string;
	    color: string;
	    files?: File[];
	    classSlots?: ClassSlot[];
	    assigments?: Assignment[];
	
	    static createFrom(source: any = {}) {
	        return new Module(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.code = source["code"];
	        this.name = source["name"];
	        this.color = source["color"];
	        this.files = this.convertValues(source["files"], File);
	        this.classSlots = this.convertValues(source["classSlots"], ClassSlot);
	        this.assigments = this.convertValues(source["assigments"], Assignment);
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

