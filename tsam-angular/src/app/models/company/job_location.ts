import { ICountry } from "../general/country";
import { IState } from "../general/state";

export interface IJobLocation{
    address: string
    city: string
    pinCode: number
    state: IState
    country: ICountry
  }