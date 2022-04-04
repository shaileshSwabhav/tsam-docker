export interface IProgrammingQuestionOption {
    id?: string
    programmingQuestionID: string
    option: string
    isCorrect: boolean
    isActive: boolean
    order: number
    optionClass: string
  }