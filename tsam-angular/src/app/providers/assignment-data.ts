import { Injectable } from "@angular/core";
import { IAssignmentSubmission as IAssignmentSubmission } from "../models/assignment-submission";
import { ITalent } from "../service/talent/talent.service";

@Injectable({
    providedIn: 'root'
})
export class AssignmentData {
    public allAssignmentSubmissions: IAssignmentSubmission[]
    public batchTalents: ITalent[]
    public allProjectSubmissions: any[]

    constructor() { }
}
