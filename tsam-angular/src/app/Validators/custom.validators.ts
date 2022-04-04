import { AbstractControl, FormGroup, ValidationErrors, ValidatorFn } from "@angular/forms"

// requirementDateValidator will validate if the date mentioned is atleast 10 days from the current date 
// if it is not, it will return inavlid as a map of the form {dateInvalid:true}
export function requirementDateValidator(control: AbstractControl): any {
    if (control.value) {
        let limitDate = new Date(new Date().setDate(new Date().getDate() + 10))
        let limitDay = limitDate.getDate()
        let limitMonth = limitDate.getMonth()
        let limitYear = limitDate.getFullYear()

        let formDate = new Date(control.value)
        let formDay = formDate.getDate()
        let formMonth = formDate.getMonth()
        let formYear = formDate.getFullYear()

        if (formYear < limitYear || (formYear == limitYear && formMonth < limitMonth)) {
            return { dateInvalid: true };
        }
        if (formYear > limitYear || (formYear == limitYear && formMonth > limitMonth)) {
            return null;
        }
        if (formDay >= limitDay) {
            return null;
        }
        return { dateInvalid: true };
    }
    return null;
}

export const fromAndToDateValidator: ValidatorFn = (control: AbstractControl): ValidationErrors | null => {
    const fromDate = control.get('fromDate');
    const toDate = control.get('toDate');

    console.log(fromDate.value > toDate.value);

    return fromDate && toDate && fromDate.value > toDate.value ? { toDateValidator: true } : { toDateValidator: false };
};

export const fromAndToTimeValidator: ValidatorFn = (control: AbstractControl): ValidationErrors | null => {
    const fromTime = control.get('fromTime');
    const toTime = control.get('toTime');
    if (fromTime.value == null || toTime.value == null) {
        return { toTimeValidator: false }
    }

    console.log(fromTime.value > toTime.value);

    return fromTime.value > toTime.value ? { toTimeValidator: true } : { toTimeValidator: false };
};

export function timeValidator(hour: string, min: string): ValidatorFn {
    return (form: FormGroup): { [key: string]: boolean } | null => {
        const timeHour = form.get(hour).value
        const timeMin = form.get(min).value
        if (timeHour == 0 && timeMin == 0) {
            const err = { toTimeValidator: true }
            form.get(hour).setErrors(err)
            form.get(min).setErrors(err)
            return
        }
        if (timeHour*60+timeMin<5){
            const err = { minMinutes: true}
            form.get(hour).setErrors(null)
            form.get(min).setErrors(err)
            return
        }

    }

}
// (control: AbstractControl): {[key:string]:any}|null{
//     const hour = control.get('timeHour');
//     const min = control.get('timeMin');
//     if (hour.value==0 && min.value==0) {
//         const err =  {toTimeValidator: true}
//         hour.setErrors(err)
//         return
//     }

//     // console.log(hour.value*60+ min.value>5);
//     // const err =  {toTimeValidator: true}
//     // hour.setErrors(err)
//     // return hour.value*60+ min.value<5 ? err :null;
// };