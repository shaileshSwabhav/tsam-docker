import { ICountry } from "./country";

export interface IUniversity {
    id?: string
    universityName: string
    countryID?: string
    countryName?: string
    isVisible?: boolean
    country: ICountry
}
  