export namespace common {
	
	export class TimeTaskInput {
	    description: string;
	    billable: string;
	    taskid: string;
	    projectid: string;
	    worktypeid: string;
	
	    static createFrom(source: any = {}) {
	        return new TimeTaskInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.description = source["description"];
	        this.billable = source["billable"];
	        this.taskid = source["taskid"];
	        this.projectid = source["projectid"];
	        this.worktypeid = source["worktypeid"];
	    }
	}

}

