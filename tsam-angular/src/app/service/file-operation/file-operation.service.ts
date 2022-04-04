import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { v4 as uuidv4 } from 'uuid';
import { Observable } from 'rxjs';
import { Constant } from '../constant';
import * as XLSX from 'xlsx';
import { environment } from "src/environments/environment";

@Injectable({
  providedIn: 'root'
})
export class FileOperationService {

  // path where file is stored.
  private readonly fileUploadLocation = environment.FILE_UPLOAD_LOACTION
  // path to php file that handles file write on host side.
  private readonly uploadApiPath = environment.UPLOAD_API_PATH
  // path with specified file name that can be used to fetch uploaded file.
  private readonly deleteAPIPath = environment.DELETE_API_PATH

  // Faculty folder
  private readonly FACULTY_VOICE_NOTES_FOLDER: string = "faculty/voice-notes"

  private readonly FOLDER_KEY: string = "folderName"
  private readonly FILE_KEY: string = "file"
  private readonly PREVIEW: string = "Preview"

  // Folders
  private readonly RESUME_FOLDER: string = "resumes/"
  private readonly RESOURCES_FOLDER: string = "resources/"
  private readonly PREVIEW_FOLDER: string = "resources/preview/"
  // Profile Image Folder.
  private readonly TALENT_IMAGE_FOLDER: string = "profile-image/talent/"

  readonly LOGO_FOLDER: string = "logos/"
  readonly BROCHURE_FOLDER: string = "brochures/"
  readonly PROJECT_FOLDER: string = "projects/"

  // Batch Folers.
  readonly BATCH_FOLDER: string = "batch/"

  // Course Folder.
  readonly COURSE_FOLDER: string = "course/"

  // Company Folders.
  readonly COMPANY_FOLDER: string = "company/"
  private readonly ONE_PAGER_FOLDER: string = "company/one-pagers/"
  private readonly TERMS_AND_CONDITION_FOLDER: string = "company/terms-and-conditions/"
  private readonly COMPANY_LOGO_FOLDER: string = "company/logos/"
  readonly OFFER_LETTER_FOLDER: string = "company/sample_offer_letters/"
  readonly BLOG_BANNER_IMAGE: string = "blog/banner-image/"
  readonly BLOG_IMAGES: string = "blog/images/"
  readonly PROGRAMMING_QUESTION_IMAGES: string = "programming-question/"

  // Module Folders.
  readonly MODULE_FOLDER: string = "module/"

  // Event Image
  readonly EVENT_IMAGE: string = "event/images/"

  // Programming concept logo.
  readonly PROGRAMMING_CONCEPT_LOGO: string = "programming-concept/logo/"

  // Talent submission images.
  readonly TALENT_SUBMISSION_IMAGE: string = "talent/assignment-submission/"
  readonly TALENT_POJECT_SUBMISSION_IMAGE: string = "talent/project-submission/images"
  readonly TALENT_POJECT_SUBMISSION_UPLOAD: string = "talent/project-submission/project-files"


  private readonly DOCUMENT_EXTENSIONS = ["application/pdf", "application/msword",
    "application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "application/vnd.openxmlformats-officedocument.presentationml.presentation", "text/plain"]

  private readonly AUDIO_EXTENSIONS = ["audio/mpeg", "audio/x-m4a", "audio/wav"]

  private readonly IMAGE_EXTENSIONS = ["image/png", "image/jpeg", "image/bmp", "image/gif"]

  private readonly ZIP_EXTENSIONS = ["application/zip", "application/rar", "application/zar",
    "application/x-zip-compressed"]


  readonly KB = 1024
  readonly MB = 1024 * this.KB
  readonly SIZE_LIMIT = 10 * this.MB

  // Download excel.
  worksheet: XLSX.WorkSheet
  workbook: XLSX.WorkBook
  excelBuffer: any
  private readonly EXCEL_TYPE = 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=UTF-8';
  private readonly EXCEL_EXTENSION = '.xlsx';

  constructor(
    private http: HttpClient,
    private constant: Constant
  ) { }

  isFileValid(file: File): Error {
    // file extension should be of word or pdf.
    let fileType = file.type
    if (!(fileType == "application/pdf") && !(fileType.includes("word"))) {
      return Error("Invalid file type.\nAllowed file types( .pdf .doc .docx)")
    }
    // 10 mb file size limit.
    if (file.size > this.SIZE_LIMIT) {
      return Error("Size should be less than 10 MB.")
    }
    return null
  }

  // Checks if file with proper document extension is being uploaded.
  isDocumentFileValid(file: File): Error {
    // file extension should be of word or pdf or ppt.
    //  let fileType = file.type
    console.log(file.type)

    if (!this.DOCUMENT_EXTENSIONS.includes(file.type)) {
      return Error("Invalid file type.\nAllowed file types( .pdf .doc .docx .ppt .pptx .txt)")
    }

    // 10 mb file size limit.
    if (file.size > this.SIZE_LIMIT) {
      return Error("Size should be less than 10 MB.")
    }
    return null
  }

  // Checks if file with proper image extension is being uploaded.
  isImageFileValid(file: File): Error {
    // file extension should be of word or pdf or ppt.
    //  let fileType = file.type

    if (!this.IMAGE_EXTENSIONS.includes(file.type)) {
      return Error("Invalid file type.\nAllowed file types( .jpeg .png)")
    }

    // 10 mb file size limit.
    if (file.size > this.SIZE_LIMIT) {
      return Error("Size should be less than 10 MB.")
    }
    return null
  }

  // Checks if file with proper audio extension is being uploaded.
  isAudioFileValid(file: File): Error {
    // file extension should be of word or pdf.
    // let fileType = file.type

    if (!this.AUDIO_EXTENSIONS.includes(file.type)) {
      return Error("Invalid file type.\nAllowed file types( .mpeg .mp3)")
    }
    // 10 mb file size limit.
    if (file.size > this.SIZE_LIMIT) {
      return Error("Size should be less than 10 MB.")
    }
    return null
  }

  uploadResume(file: File): Observable<any> {
    return this.uploadFile(file, this.RESUME_FOLDER)
  }

  uploadFacultyVoiceNote(file: File): Observable<any> {
    return this.uploadFile(file, this.FACULTY_VOICE_NOTES_FOLDER)
  }

  // Upload talent image*************************************** change path*********************************.
  uploadTalentImage(file: File): Observable<any> {
    return this.uploadFile(file, this.TALENT_IMAGE_FOLDER)
  }

  uploadBlogBannerImage(file: File): Observable<any> {
    return this.uploadFile(file, this.BLOG_BANNER_IMAGE)
  }

  uploadProgrammingConceptLogo(file: File): Observable<any> {
    return this.uploadFile(file, this.PROGRAMMING_CONCEPT_LOGO)
  }

  uploadTalentSubmissionImage(file: File): Observable<any> {
    return this.uploadFile(file, this.TALENT_SUBMISSION_IMAGE)
  }

  uploadTalentProjectSubmissionImage(file: File): Observable<any> {
    return this.uploadFile(file, this.TALENT_POJECT_SUBMISSION_IMAGE)
  }

  uploadOnePager(file: File): Observable<any> {
    return this.uploadFile(file, this.ONE_PAGER_FOLDER)
  }

  uploadResource(file: File, fileType?: string): Observable<any> {
    let folder_name = this.RESOURCES_FOLDER
    if (fileType == this.PREVIEW) {
      folder_name = this.PREVIEW_FOLDER
    }
    return this.uploadFile(file, folder_name)
  }

  uploadBrochure(file: File, folderPath: string): Observable<any> {
    return this.uploadFile(file, folderPath)
  }

  uploadTermsAndCondition(file: File): Observable<any> {
    return this.uploadFile(file, this.TERMS_AND_CONDITION_FOLDER)
  }

  uploadLogo(file: File, folderPath: string): Observable<any> {
    console.log(folderPath);
    return this.uploadFile(file, folderPath)
  }

  uploadOfferLetter(file: File): Observable<any> {
    return this.uploadFile(file, this.OFFER_LETTER_FOLDER)
  }

  uploadImage(file: File, folderPath: string): Observable<any> {
    return this.uploadFile(file, folderPath)
  }

  uploadFile(file: File, folderPath: string): Observable<any> {
    return new Observable((observer) => {
      let ext = this.getFileExtension(file.name)
      let fileName = uuidv4() + ext;
      console.log("File name:", fileName);

      let formData = new FormData();
      let fileURL = this.fileUploadLocation + folderPath + "/" + fileName;
      console.log("File URL:", fileURL);
      formData.append(this.FILE_KEY, file, fileName);
      // Appending the folder name.
      formData.append(this.FOLDER_KEY, folderPath)

      // calling the API present in phpPath(server)
      this.http.post<any>(this.uploadApiPath, formData).subscribe((data: any) => {
        //  for delete # need to check later #niranjan
        console.log("Server data :", data)
        if (!data.isUploadSuccessful) {
          observer.error("Upload failed")
          return
        }
        observer.next(fileURL)
      }, (error) => {
        console.log("upload file error :", error)
        if (error.statusText == "Unknown Error") {
          observer.error("Check internet connection.")
          return
        }
        observer.error("Some error occurred.")
      })
    });
  }

  uploadProject(file: File): Observable<any> {
    return this.uploadFile(file, this.TALENT_POJECT_SUBMISSION_UPLOAD)
  }

  // Checks if file with proper document extension is being uploaded.
  isProjectFileValid(file: File): Error {
    // file extension should be of word or pdf or ppt.
    //  let fileType = file.type
    console.log(file.type)

    if (!this.ZIP_EXTENSIONS.includes(file.type)) {
      return Error("Invalid file type.\nAllowed file types( .zip .rar .zar)")
    }

    // 10 mb file size limit.
    if (file.size > this.SIZE_LIMIT) {
      return Error("Size should be less than 10 MB.")
    }
    return null
  }

  // ==================================================================================================================

  // fileName is a mandatory parameter. Needs to be changed #Niranjan
  deleteUploadedFile(fileName?: string): Observable<any> {
    return new Observable((observer) => {
      if (fileName == "" || !fileName) {
        observer.error("no file")
        return
      }
      let strArray = fileName.split(this.fileUploadLocation)
      if (strArray.length < 1) {
        observer.error("Improper path to file")
        return
      }
      let filePath = strArray[1]
      console.log("filepath", filePath)
      this.http.delete<any>(this.deleteAPIPath, {
        params: { "filePath": filePath }
      }).subscribe((data: any) => {
        observer.next(data)
      }, (err) => {
        observer.error(err)
      })
    });
  }

  uploadExcel(file: File): Observable<any> {

    let workBook: XLSX.WorkBook
    let jsonData = null

    const reader = new FileReader();
    return new Observable<any>((observer) => {
      if (!file.type.includes('sheet') && !file.type.includes('excel')) {
        observer.error("Invalid file type. Only excel files(.xlsx, .xls etc) are allowed")
      }
      reader.onload = () => {
        const data = reader.result;
        // will try type string
        try {
          workBook = XLSX.read(data, { type: 'binary' })
          jsonData = XLSX.utils.sheet_to_json(workBook.Sheets[workBook.SheetNames[0]]);
        } catch {
          observer.error("Couldn't read file.")
        }
        observer.next(jsonData)
        observer.complete()
      }
      // if failed
      reader.onerror = (err: any): void => {
        console.error(err)
        observer.error("Some error occured. Try again.");
      }
      reader.readAsBinaryString(file)
    })
  }

  // Create excel workbook.
  createExcelWorkBook(): void {
    this.workbook = XLSX.utils.book_new()
  }

  // Create excel worksheet.
  createExcelWorkSheet(json: any[]): void {
    this.worksheet = XLSX.utils.json_to_sheet(json)
  }

  // Append to excel worksheet.
  appendToExcelWorkSheet(json: any[]): void {
    XLSX.utils.sheet_add_json(this.worksheet, json, { skipHeader: true, origin: -1 })
  }

  // Save the excel file.
  saveAsExcelFile(fileName: string): void {
    XLSX.utils.book_append_sheet(this.workbook, this.worksheet, 'talent')
    XLSX.writeFile(this.workbook, fileName + this.EXCEL_EXTENSION)
  }

  getFileExtension(filename: string): string {
    var parts = filename.split('.');
    return "." + parts[parts.length - 1];
  }

}
/*
reader.readAsArrayBuffer((<any>ev.target).files[0]);
 jsonData = workBook.SheetNames.reduce((initial: any, name: any) => {
          const sheet = workBook.Sheets[0];
          initial[name] = XLSX.utils.sheet_to_json(sheet, { blankrows: true, raw: true });
          return initial;
        }, {}); */

export interface IUploadStatus {
  status: string
  message: string
  filePath: string
}
