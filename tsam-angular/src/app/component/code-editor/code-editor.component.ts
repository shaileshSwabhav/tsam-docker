import { ChangeDetectorRef, Component, ElementRef, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges, ViewChild, ViewRef } from '@angular/core';
import { FormGroup, FormControl, Validators, FormBuilder, Form } from '@angular/forms';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;

import * as ace from 'ace-builds';
import { CodeEditorService } from 'src/app/service/code-editor/code-editor.service';
// import { select } from 'd3';

// Constants for theme.
const TWILIGHT_THEME = 'ace/theme/twilight';
const GIT_THEME = 'ace/theme/github';
const NORD_DARK_THEME = 'ace/theme/nord_dark';

// Constants for language mode.
const LANG_C_CPP = "ace/mode/c_cpp"
const LANG_PHP = "ace/mode/php"
const LANG_JAVA = "ace/mode/java"
const LANG_GOLANG = "ace/mode/golang"
const LANG_PYTHON = "ace/mode/python"
const LANG_JAVASCRIPT = "ace/mode/javascript"

// For testing.
const CODE = `func main(){fmt.Println("HELLO")}`

// Constant for code compiler.
const RUNNING = "running"
const GUEST = "guest"

// Constant for prgramming language array.
const LANGUAGE_ARRAY = [
  {
    name: "C",
    path: "ace/mode/c_cpp",
    startCode:
      `#include <stdio.h> 

        int main() 
          {
            printf("Hello");
            return 0;
          }`,
    nameForCompiler: "c",
    tempCode: "",
    sphereCompilerID: 11
  },
  {
    name: "C++",
    path: "ace/mode/c_cpp",
    startCode:
      `#include <iostream>

      int main(){
        cout<<"Hello";
        return 0;
      }`,
    nameForCompiler: "cpp",
    tempCode: "",
    sphereCompilerID: 41
  },
  {
    name: "PHP",
    path: "ace/mode/php",
    startCode:
      `<?php 
        echo "Hello"
        ?>`,
    nameForCompiler: "php",
    tempCode: "",
    sphereCompilerID: 29
  },
  {
    name: "Python",
    path: "ace/mode/python",
    startCode: `print("Hello")`,
    nameForCompiler: "python",
    tempCode: "",
    sphereCompilerID: 99
  },
  {
    name: "Java",
    path: "ace/mode/java",
    startCode:
      `public class Main{
        public static void main(String[] args) {
        System.out.println("Hello");
        }
      }`,
    nameForCompiler: "java",
    tempCode: "",
    sphereCompilerID: 10
  },
  {
    name: "Go",
    path: "ace/mode/golang",
    startCode:
      `package main 
      
    import "fmt"
    
    func main(){
      fmt.Println("Hello")
    }`,
    nameForCompiler: "go",
    tempCode: "",
    sphereCompilerID: 114
  },
  {
    name: "Java Script",
    path: "ace/mode/javascript",
    startCode: `console.log("Hello");`,
    nameForCompiler: "javascript",
    tempCode: "",
    sphereCompilerID: 56
  },
]

@Component({
  selector: 'app-code-editor',
  templateUrl: './code-editor.component.html',
  styleUrls: ['./code-editor.component.css']
})
export class CodeEditorComponent implements OnInit {

  // Definite Assignment Assertion
  @ViewChild('codeEditor') codeEditorElmRef!: ElementRef
  @ViewChild('hiddenText') hiddenText!: ElementRef
  private codeEditor!: ace.Ace.Editor
  private editorBeautify!: any
  public submittedCode!: string

  // Event emitter for sending values to parent component.
  @Output() childToParentEvent = new EventEmitter<ILanguageAndCode>();

  // Programming language.
  selectedLanguage: any
  languageForm: FormGroup

  // Values from parent component.
  @Input() valuesFromParent: any

  // Flags.
  isReadOnly: boolean
  isSolution: boolean
  isCodeRunning: boolean
  showInputForm: boolean

  // Spinner.



  // Output.
  output: string
  isErrors: boolean

  // Input.
  inputForm: FormGroup

  // Test Case.
  testCaseJudgeForm: FormGroup
  testCases: any[]

  constructor(
    private formBuilder: FormBuilder,
    private changeDetectorRef: ChangeDetectorRef,
    private spinnerService: SpinnerService,
    private codeEditorService: CodeEditorService,
  ) {

    // Create language form and set the first element of array for the form.
    this.createLanguageForm()

    // Create input form.
    this.createInputForm()
    this.createTestCaseJudgeForm()

    // Flags.
    this.isReadOnly = false
    this.isSolution = false
    this.isCodeRunning = false
    this.showInputForm = true

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Output"


    // Output.
    this.output = ""
    this.isErrors = false
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  ngAfterViewInit(): void {
    // console.log("oninit", this.codeEditorElmRef)

    const element = this.codeEditorElmRef.nativeElement
    const editorOptions = this.getEditorOptions()

    // This has to be provided in drop down & will change dynamically.
    ace.config.set("fontSize", "14px")
    ace.config.set(
      "basePath",
      "https://unpkg.com/ace-builds@1.4.12/src-noconflict"
    );

    // ace.edit will not create a new instance of code editor if a prev instance is passed.
    this.codeEditor = ace.edit(element, editorOptions);

    // Set initial theme and language 
    this.codeEditor.setTheme(TWILIGHT_THEME);
    this.codeEditor.session.setMode(LANG_PYTHON);

    // change detection and printing value
    this.codeEditor.on("change", () => {
    });
    this.codeEditor.setShowFoldWidgets(true); // for the scope fold feature
    this.editorBeautify = ace.require('ace/ext/beautify')

  }

  // =========================================CREATE FORMS=================================================

  // Create language form.
  createLanguageForm(): void {
    this.languageForm = this.formBuilder.group({
      programmingLanguage: new FormControl(null, [Validators.required]),
    })
  }

  // Create input form.
  createInputForm(): void {
    this.inputForm = this.formBuilder.group({
      input: new FormControl(null),
    })
  }

  // Create test case form.
  createTestCaseJudgeForm(): void {
    this.testCaseJudgeForm = this.formBuilder.group({
      testCaseJudge: new FormControl(null),
    })
  }

  // =========================================FORMAT FUNCTIONS=================================================

  // Format the code according to language.
  public beautifyContent() {
    if (this.codeEditor && this.editorBeautify) {
      const session = this.codeEditor.getSession()
      this.editorBeautify.beautify(session)
    }
  }

  // missing property on EditorOptions 'enableBasicAutocompletion'
  private getEditorOptions(): Partial<ace.Ace.EditorOptions> & { enableBasicAutocompletion?: boolean; } {
    const basicEditorOptions: Partial<ace.Ace.EditorOptions> = {
      highlightActiveLine: true,
      minLines: 14,
      maxLines: 1000,
      enableAutoIndent: true
    }

    const extraEditorOptions = {
      enableBasicAutocompletion: true,
      enableLiveAutocompletion: true
    }
    const margedOptions = Object.assign(basicEditorOptions, extraEditorOptions)
    return margedOptions
  }

  // Set the language mode for formatting and compilation.
  setLanguageModeOfCodeEditor(): void {
    this.codeEditor.getSession().setMode(this.languageForm.get('programmingLanguage').value.path)
  }

  // Set the theme on theme change.
  onThemeChange(event: any) {
    this.codeEditor.setTheme(event.target.value)
  }

  // =========================================ON CLICK FUNCTIONS=================================================

  // On clicking submit button.
  validateLanguageAndCode(): void {
    if (this.languageForm.invalid) {
      this.languageForm.markAllAsTouched()
      return
    }
    if (!this.codeEditor.getValue() || this.codeEditor.getValue() == "") {
      alert("Please type your answer before submitting")
      return
    }
    this.sendMessageToParent(this.languageForm.get('programmingLanguage').value.name, this.codeEditor.getValue())
  }

  // On clearing the editor content fill it with a basic code.
  clear() {
    this.codeEditor.setValue(this.languageForm.get('programmingLanguage').value.startCode)
  }

  // Compile and run the code.
  run() {
    let languageAndCode: any = {
      code: this.codeEditor.session.getValue(),
      language: this.languageForm.get('programmingLanguage').value.nameForCompiler,
      compilerId: this.languageForm.get('programmingLanguage').value.sphereCompilerID,
    }

    // If language is go, javascript or php then use paiza compiler, else use codex compiler.
    // if (languageAndCode.language == "go" || languageAndCode.language == "javascript" || languageAndCode.language == "php"){
    //   this.sendCode(languageAndCode)
    // }
    // else{
    //   this.sendAndReceive(languageAndCode)
    // }
    this.sendCodeSphere(languageAndCode)
    // this.getStatusSphere()

    // this.submittedCode = this.codeEditor.session.getValue()
    // let inner = this.codeEditorElmRef.nativeElement.innerHTML
  }

  // =========================================OTHER FUNCTIONS=================================================

  // Get the string from the editor.
  private getCode() {
    const code = this.codeEditor.getValue()
  }

  // Set the language in language form and code in code editor.
  setLanguageAndCode(): void {

    // If selected problem has changed then reset the temp code af all languages in language array.
    if (this.valuesFromParent.isSelectedProblemChanged) {
      for (let i = 0; i < this.languageArray.length; i++) {
        this.languageArray.tempCode = ""
      }
      // for (let i = 0; i < this.languageArray.length; i++){
      //   console.log(this.languageArray[i].tempCode)
      // }
    }

    // Set the language.
    for (let i = 0; i < this.languageArray.length; i++) {
      if (this.languageArray[i].name == this.valuesFromParent.languageName) {
        this.languageForm.get('programmingLanguage').setValue(this.languageArray[i])
        this.selectedLanguage = this.languageArray[i]
        this.setLanguageModeOfCodeEditor()
        this.clear()
        break
      }
    }

    // If isReadOnly is false and code is null then set the code editor to clear.
    if (!this.valuesFromParent.isReadOnly && !this.valuesFromParent.code) {

      // If there is temp code in selected language then set it to the editor, else set the editor value as start code.
      if (this.selectedLanguage.tempCode != "" && !this.valuesFromParent.isSelectedProblemChanged) {
        this.codeEditor.setValue(this.selectedLanguage.tempCode)
      } else {
        this.codeEditor.setValue(this.languageForm.get('programmingLanguage').value.startCode)
      }
      this.codeEditor.getSession().setMode(this.languageForm.get('programmingLanguage').value.path)
      this.beautifyContent()
      this.codeEditor.setReadOnly(false)
      this.isReadOnly = this.valuesFromParent.isReadOnly
      this.isSolution = this.valuesFromParent.isSolution
      this.isCodeRunning = false
    }

    // If isReadOnly is false and code is not null then set the answer.
    if (!this.valuesFromParent.isReadOnly && this.valuesFromParent.code) {
      this.codeEditor.setValue(this.valuesFromParent.code)
      this.beautifyContent()
      this.codeEditor.setReadOnly(false)
      this.isReadOnly = this.valuesFromParent.isReadOnly
      this.isSolution = this.valuesFromParent.isSolution
      this.isCodeRunning = false
    }

    // If isReadOnly is true and code is not null then set the answer.
    if (this.valuesFromParent.isReadOnly && this.valuesFromParent.code) {
      this.codeEditor.setValue(this.valuesFromParent.code)
      this.beautifyContent()
      this.codeEditor.setReadOnly(true)
      this.isReadOnly = this.valuesFromParent.isReadOnly
      this.isSolution = this.valuesFromParent.isSolution
      this.isCodeRunning = true
    }

    // To detect changes only after isAnswered value changes and component is not destroyed.
    if (this.changeDetectorRef && !(this.changeDetectorRef as ViewRef).destroyed) {
      this.changeDetectorRef.detectChanges()
    }
  }

  // On values from parent are changed, invoked by parent component.
  valuesFromParentChange(languageAndCode: ILanguageAndCode): void {

    // Set the valuesFromParent object with the values received from parent component.
    this.valuesFromParent = languageAndCode
    this.testCases = languageAndCode['testCases']

    // Set language and code.
    this.createLanguageForm()
    this.setLanguageAndCode()
  }

  // On changing language in language form.
  onLanguageChange() {

    // Reset the input form.
    this.createInputForm()

    // Set the temp code for the previous selected language in language form.
    for (let i = 0; i < this.languageArray.length; i++) {
      if (this.languageArray[i].name == this.selectedLanguage.name) {
        this.languageArray[i].tempCode = this.codeEditor.getValue()
      }
    }
    this.output = ""  // #Niranjan
    this.sendMessageToParent(this.languageForm.get('programmingLanguage').value.name, null)
  }

  // Sends a message to parent component.
  sendMessageToParent(languageName: string, code: string) {

    let languageAndAnswer: ILanguageAndCode = {
      languageName: languageName,
      code: code,
    }
    this.childToParentEvent.emit(languageAndAnswer)
  }

  // Toggle custom input and test case input.
  toggleCustomInput(): void {
    this.showInputForm = !this.showInputForm
  }

  // =========================================EXTERNAL COMPILER FUNCTIONS=================================================

  // **************************** Paiza Compiler *****************************************
  // Send code to compiler API.
  sendCode(languageNameAndCode: any): void {
    this.spinnerService.loadingMessage = "Getting Output"


    let queryParams: any = {
      source_code: languageNameAndCode.code,
      language: languageNameAndCode.language,
      api_key: GUEST
    }
    this.codeEditorService.sendCode(queryParams).subscribe((response) => {
      if (response.status == RUNNING) {
        this.getStatus(response.id)
      }
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Gets status of compilation.
  getStatus(sessionID: string): void {
    this.spinnerService.loadingMessage = "Getting Output"


    let queryParams: any = {
      id: sessionID,
      api_key: GUEST
    }
    this.codeEditorService.getStatus(queryParams).subscribe((response) => {
      if (response.status == RUNNING) {
        this.getStatus(sessionID)
        return
      }
      this.receiveOutput(sessionID)
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Receice output from compiler.
  receiveOutput(sessionID: string): void {
    this.spinnerService.loadingMessage = "Getting Output"


    let queryParams: any = {
      id: sessionID,
      api_key: GUEST
    }
    this.codeEditorService.receiveOutput(queryParams).subscribe((response) => {
      if (response.build_stderr) {
        this.output = response.build_stderr
        this.isCodeRunning = false
      }
      if (response.stderr) {
        this.output = response.stderr
        this.isCodeRunning = false
      }
      if (response.stdout) {
        this.output = response.stdout
        this.isCodeRunning = true
      }
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // **************************** Codex Compiler *****************************************

  // Send code and receive output from compiler API.
  sendAndReceive(languageNameAndCode: any): void {
    this.spinnerService.loadingMessage = "Getting Output"


    let queryParams: any = {
      code: languageNameAndCode.code,
      language: languageNameAndCode.language,
      input: ""
    }
    this.codeEditorService.sendCode(queryParams).subscribe((response) => {
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // **************************** Sphere Compiler *****************************************

  // Send code to compiler API.
  sendCodeSphere(languageNameAndCode: any): void {
    this.spinnerService.loadingMessage = "Getting Output"


    let sourceCodeAndInput: any = {
      compilerId: languageNameAndCode.compilerId,
      source: languageNameAndCode.code,
      input: this.inputForm.get('input').value
    }
    this.codeEditorService.sendCodeSphere(sourceCodeAndInput).subscribe((response) => {
      this.getOutputShpere(response.id)
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Gets output of compilation.
  getOutputShpere(outputID: number): void {
    this.spinnerService.loadingMessage = "Getting Output"


    let queryParams: any = {
      access_token: this.codeEditorService.tempAccessToken
    }
    this.codeEditorService.getOutputShpere(outputID, queryParams).subscribe((response) => {
      // If code is still not in accepted state then call the get output API again.
      if (response.result.status.code >= 0 && response.result.status.code <= 3) {

        setTimeout(() => { this.getOutputShpere(outputID) }, 1000)
      }
      else {

        if (response.result.streams.error) {
          this.getResultShpere(response.result.streams.error.uri)
          return
        }
        if (response.result.streams.cmpinfo) {
          this.getResultShpere(response.result.streams.cmpinfo.uri)
          return
        }
        if (response.result.streams.output) {
          this.getResultShpere(response.result.streams.output.uri)
          return
        }
        if (response.result.status) {
          this.output = "Could not get output... please try again later"
        }
      }
    }, error => {


      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Gets result of compilation.
  getResultShpere(resultURL: string): void {
    this.spinnerService.loadingMessage = "Getting Output"


    this.codeEditorService.getResultShpere(resultURL).subscribe((response) => {
      this.output = response
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Add judge.
  addJudgeSphere(): void {
    this.spinnerService.loadingMessage = "Getting Output"


    let judge: any = {
      compilerId: 10,
      name: "master judge java addition",
      source: this.testCaseJudgeForm.get('testCaseJudge').value,
      typeId: 1
    }
    this.codeEditorService.addJudgeSphere(judge).subscribe((response) => {
      console.log(response)
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Add problem.
  addProblemSphere(): void {
    this.spinnerService.loadingMessage = "Getting Output"


    let problem: any = {
      name: "java problem addition",
      masterjudgeId: 2606,
      typeId: 4
    }
    this.codeEditorService.addProblemSphere(problem).subscribe((response) => {
      console.log(response)
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Add test case.
  addTestCaseSphere(): void {
    this.spinnerService.loadingMessage = "Getting Output"


    let testCase: any = {
      input: this.testCases[2].input,
      output: this.testCases[2].output,
      judgeId: 2605,
      problemID: 80121
    }
    this.codeEditorService.addTestCaseSphere(testCase).subscribe((response) => {
      console.log(response)
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Add submission.
  addSubmissionSphere(): void {
    this.spinnerService.loadingMessage = "Getting Output"


    let submission: any = {
      tests: "0,1,2",
      compilerId: 10,
      source: this.codeEditor.session.getValue(),
      problemId: 80121
    }
    this.codeEditorService.addSubmissionSphere(submission).subscribe((response) => {
      console.log(response)
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // **************************** Remote Code Compiler *****************************************

  // Send code to compiler API.
  sendCodeRemoteCode(languageNameAndCode: any): void {
    this.spinnerService.loadingMessage = "Getting Output"


    let sourceCodeAndInput: any = {
      sourceCode: languageNameAndCode.code
    }
    this.codeEditorService.sendCodeRemoteCode(sourceCodeAndInput).subscribe((response) => {
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // =========================================CONSTANT GETTERS=================================================
  get twilightTheme(): string {
    return TWILIGHT_THEME
  }
  get githubTheme(): string {
    return GIT_THEME
  }
  get nordDarkTheme(): string {
    return NORD_DARK_THEME
  }
  // Languages
  get langCCpp(): string {
    return LANG_C_CPP
  }
  get langGolang(): string {
    return LANG_GOLANG
  }
  get langJava(): string {
    return LANG_JAVA
  }
  get langJavascript(): string {
    return LANG_JAVASCRIPT
  }
  get langPhp(): string {
    return LANG_PHP
  }
  get langPython(): string {
    return LANG_PYTHON
  }
  get languageArray(): any {
    return LANGUAGE_ARRAY
  }

}

export interface ILanguageAndCode {
  languageName: string
  code: string
}
