import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import { LocalService } from '../storage/local.service';

@Injectable({
  providedIn: 'root'
})
export class EnquiryFormService {

  httpHeaders: HttpHeaders;
  generalURL: string;

  constructor(
    private _http: HttpClient,
    private _constant: Constant,
    private localService: LocalService
  ) {
    this.generalURL = `${_constant.BASE_URL}/tenant/${_constant.TENANT_ID}`
  }

  //***************************************************Enquiry related api calls****************************************************
  // Add new enquiry
  addEnquiry(enquiry: any): Observable<any> {
    // let enquiryJson: string = JSON.stringify(enquiry);
    // console.log("enquiry" + enquiryJson)
    return this._http.post<any>(`${this.generalURL}/talent-enquiry-form`, enquiry);
  }
}
//**************************************ENQUIRY******************************************************* */
export interface IEnquiry {
  id?: string,

	// Personal Details.
  address: string,
  city: string,
  pinCode: number,
  state: IState,
  country: ICountry,
  code: string,
  firstName: string,
  lastName: string,
  email: string,
  contact: string,
  enquiryType: string,
  academicYear: number,
  resume: string,

	// Maps.
  technologies: ITechnology[],
  courses: ICourse[],

	// Related table IDs.
  sourceID: string,

	// Flags.
  isMastersAbroad: boolean,
  isExperience: boolean,

	// Child tables.
  academics: IAcademic[],
  experiences: IExperience[],
  mastersAbroad: IMastersAbroad,

	// Social media.
  facebookUrl: string,
  instagramUrl: string,
  githubUrl: string,
  linkedInUrl: string,
}

//**************************************EXPERIENCES******************************************************* */
export interface IExperience {
  id?: string,

	// Maps.
  technologies: ITechnology[],

	// Related table IDs.
  designationID: string,
  enquiryID: string,

	// Other fields.
  company: string,
  fromDate: string,
  toDate: string,
  package: number,
}

//**************************************ACADEMICS******************************************************* */
export interface IAcademic {
  id?: string,

	// Related table IDs.
  degreeID: string,
  collegeID: string,
  specializationID: string,
  enquiryID: string

	// Other fields.
  college: string
  percentage: number,
  passout: number,
}

//**************************************MASTERS ABROAD******************************************************* */
export interface IMastersAbroad {
  id?: string

	// Related table IDs.
  talentID?: string
  enquiryID?: string
  degreeID: string

	// Maps.
  countries: ICountry[]
  universities: IUniversity[]

	// Child tables.
  scores: IScore[]

  // Other fields.
  yearOfMS: number
}

export interface IScore {
  id?: string

	// Related table IDs.
  examinationID: string
  mastersAbroadID?: string

	// Other fields.
  marksObtained: number
}

export interface IExamination {
  id?: string
  name: string
  totalMarks: number
}

//**************************************OTHERS******************************************************* */
export interface ICollegeBranch {
  id?: string
  branchName: string
  code: string
}

export interface IUniversity {
  id?: string
  universityName: string
  countryID?: string
  countryName?: string
  isVisible?: boolean
  country: ICountry
}

export interface IState {
  id?: string,
  name: string,
  countryID?: string
}

export interface ICountry {
  id?: string,
  name: string
}

export interface IDegree {
  id?: string
  name: string
}

export interface ISpecialization {
  id?: string
  branchName: string,
  degreeID: string
}

export interface ISource {
  id?: string,
  name: string
  urlName: string
  description: string
}

export interface IDesignation {
  id?: string,
  position: string
}

export interface ITechnology {
  id?: string,
  language: string
}

export interface ICourse {
  id?: string,
  code: string
  name: string
}


