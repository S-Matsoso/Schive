export namespace models {
	
	export class Flashcard {
	    id: string;
	    moduleId: string;
	    question: string;
	    answer: string;
	    mastered: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Flashcard(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.moduleId = source["moduleId"];
	        this.question = source["question"];
	        this.answer = source["answer"];
	        this.mastered = source["mastered"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	    day: string;
	    startHour: number;
	    lengthHrs: number;
	    room: string;
	    module?: Module;
	
	    static createFrom(source: any = {}) {
	        return new ClassSlot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.moduleId = source["moduleId"];
	        this.day = source["day"];
	        this.startHour = source["startHour"];
	        this.lengthHrs = source["lengthHrs"];
	        this.room = source["room"];
	        this.module = this.convertValues(source["module"], Module);
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
	export class Folder {
	    id: string;
	    moduleId: string;
	    name: string;
	    color?: string;
	    // Go type: time
	    createdAt: any;
	    fileCount: number;
	    files?: File[];
	
	    static createFrom(source: any = {}) {
	        return new Folder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.moduleId = source["moduleId"];
	        this.name = source["name"];
	        this.color = source["color"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.fileCount = source["fileCount"];
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
	export class File {
	    id: string;
	    moduleId: string;
	    folderId?: string;
	    assignmentId?: string;
	    name: string;
	    path: string;
	    fileType: string;
	    fileSize: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.moduleId = source["moduleId"];
	        this.folderId = source["folderId"];
	        this.assignmentId = source["assignmentId"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.fileType = source["fileType"];
	        this.fileSize = source["fileSize"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class Module {
	    id: string;
	    code: string;
	    name: string;
	    color: string;
	    files?: File[];
	    folders?: Folder[];
	    classSlots?: ClassSlot[];
	    assignments?: Assignment[];
	    flashcards?: Flashcard[];
	
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
	        this.folders = this.convertValues(source["folders"], Folder);
	        this.classSlots = this.convertValues(source["classSlots"], ClassSlot);
	        this.assignments = this.convertValues(source["assignments"], Assignment);
	        this.flashcards = this.convertValues(source["flashcards"], Flashcard);
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
	export class Assignment {
	    id: string;
	    moduleId: string;
	    title: string;
	    description: string;
	    assignmentNumber: number;
	    status: string;
	    priority: string;
	    progress: number;
	    // Go type: time
	    dueDate?: any;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    module?: Module;
	    files?: File[];
	
	    static createFrom(source: any = {}) {
	        return new Assignment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.moduleId = source["moduleId"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.assignmentNumber = source["assignmentNumber"];
	        this.status = source["status"];
	        this.priority = source["priority"];
	        this.progress = source["progress"];
	        this.dueDate = this.convertValues(source["dueDate"], null);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.module = this.convertValues(source["module"], Module);
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
	
	
	
	

}

