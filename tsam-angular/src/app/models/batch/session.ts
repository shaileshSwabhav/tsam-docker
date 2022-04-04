export interface IBatchSession {
    id: string
    name: string
    order: number
    studentOutput: string
    hours: string
    sessionID: string
    subSessions: IBatchSession[]
    courseID: string
    isChecked: boolean
    viewSubSessionClicked: boolean
    cardColumn: string
    materialIcon: string
  }