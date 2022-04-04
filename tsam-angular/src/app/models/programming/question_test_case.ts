export interface IProgrammingQuestionTestCase{
    id?: string
    programmingQuestionID: string
    input: string
    output: string
    explanation: string
    isActive: boolean
    isHidden: boolean
  }