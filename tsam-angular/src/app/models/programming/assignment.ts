import { IProgrammingAssignmentSubTask } from "./assignment_sub_task";
import { IProgrammingQuestion } from "./question";

export interface IProgrammingAssignment {
	id?: string
	title: string
	taskDescription: string
	timeRequired: string
	complexityLevel: number
	score: number
	additionalComments: string
	sourceURL: string
	source?: string
  programmingAssignmentType: string
	programmingQuestion: IProgrammingQuestion[]
	programmingAssignmentSubTask: IProgrammingAssignmentSubTask[]
}
