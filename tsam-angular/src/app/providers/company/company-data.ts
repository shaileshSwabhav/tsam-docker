import { Injectable } from "@angular/core"
import { ICompanyRequirement } from 'src/app/service/company/company.service'

@Injectable({
    providedIn: 'root'
})
export class CompanyData {
    // public talentSearchParams: ISearchTalentParams
    public companyRequirement: ICompanyRequirement

    constructor() { }
}
