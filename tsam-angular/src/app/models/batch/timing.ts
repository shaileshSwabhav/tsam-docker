import { IDay } from "../general/day";

export interface IBatchTiming {
    id?: string
    batchID: string
    day: IDay
    fromTime: string
    toTime: string
  }
  