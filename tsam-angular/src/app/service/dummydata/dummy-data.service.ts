import { Injectable } from '@angular/core';

@Injectable({
      providedIn: 'root'
})
export class DummyDataService {

      private ratingFeedbackQuestions: any[];

      constructor(
      ) {
            this.loadRequirementRating();
      }

      getRatingFeedbackQuestions() {
            return new Promise((resolve, reject) => {
                  resolve(this.ratingFeedbackQuestions);
            })
      }

      loadRequirementRating(): void {

            this.ratingFeedbackQuestions = [
                  {
                        "question": "When do you give increments / promotions?",
                        "columnName": "increment",
                        "maxScore": 4,
                        "options": [
                              {
                                    "key": 1,
                                    "value": "Not Fixed"
                              },
                              {
                                    "key": 2,
                                    "value": "Once in a year"
                              },
                              {
                                    "key": 4,
                                    "value": "After 6 months"
                              }
                        ],
                  },
                  {
                        "question": "When do you have holidays in the week?",
                        "columnName": "weeklyHoliday",
                        "maxScore": 4,
                        "options": [

                              {
                                    "key": 2,
                                    "value": "Only sunday"
                              },
                              {
                                    "key": 3,
                                    "value": "Alternate Saturdays"
                              },
                              {
                                    "key": 4,
                                    "value": "Weekends"
                              },
                              {
                                    "key": 1,
                                    "value": "Other"
                              },
                        ],
                  },
                  {
                        "question": "Any specific qualification requiremnt?",
                        "columnName": "qualification",
                        "maxScore": 4,
                        "options": [
                              {
                                    "key": 1,
                                    "value": "B.E/B.Tech in computers/IT only"
                              },
                              {
                                    "key": 2,
                                    "value": "Any B.E/B.Tech"
                              },
                              {
                                    "key": 3,
                                    "value": "Any B.Sc/BCA/B.E./B.Tech"
                              },
                              {
                                    "key": 4,
                                    "value": "Any Degree"
                              }
                        ],
                  },
                  {
                        "question": "Do you have a bond period?",
                        "columnName": "bondPeriod",
                        "maxScore": 10,
                        "options": [
                              {
                                    "key": 1,
                                    "value": "2 year bond or more"
                              },
                              {
                                    "key": 4,
                                    "value": "1 year bond or less"
                              },
                              {
                                    "key": 10,
                                    "value": "No bond"
                              },
                        ],
                  },

                  {
                        "question": "Do you have any gender preference?",
                        "columnName": "genderPreference",
                        "maxScore": 0,
                        "options": [
                              {
                                    "key": 0,
                                    "value": "Female"
                              },
                              {
                                    "key": 0, 
                                    "value": "Male"
                              },
                              {
                                    "key": 0,
                                    "value": "Any"
                              }
                        ],
                  },
                  {
                        "question": "Where do you want your applicants to be residing?",
                        "columnName": "talentLocation",
                        "maxScore": 6,
                        "options": [
                              {
                                    "key": 6,
                                    "value": "Any location(remote)"
                              },
                              {
                                    "key": 4,
                                    "value": "Any where in mumbai"
                              },
                              {
                                    "key": 1,
                                    "value": "Specific Line(Harbour,Western etc)"
                              }
                        ],
                  },
                  {
                        "question": "What will be the work shifts?",
                        "columnName": "workShift",
                        "maxScore": 6,
                        "options": [
                              {
                                    "key": 1,
                                    "value": "Rotational(Night/day)"
                              },
                              {
                                    "key": 1,
                                    "value": "Night"
                              },
                              {
                                    "key": 3,
                                    "value": "Fixed working hours"
                              },
                              {
                                    "key": 6,
                                    "value": "Flexible hours"
                              }
                        ],
                  },
                  {
                        "question": "What is the type of the company?",
                        "columnName": "companyType",
                        "maxScore": 5,
                        "options": [
                              {
                                    "key": 1,
                                    "value": "Startup"
                              },
                              {
                                    "key": 2,
                                    "value": "Funded"
                              },
                              {
                                    "key": 5,
                                    "value": "Public listed"
                              }
                        ],
                  },
                  {
                        "question": "What will be the joining period?",
                        "columnName": "joiningPeriod",
                        "maxScore": 3,
                        "options": [
                              {
                                    "key": 1,
                                    "value": "Immediate"
                              },
                              {
                                    "key": 2,
                                    "value": "30 days"
                              },
                              {
                                    "key": 3,
                                    "value": "2 months"
                              }
                        ]
                  }


                  // {
                  //       "question": "What kind of developers you are looking at?",
                  //       "columnName": "criteria8",
                  //       "maxScore": 5,
                  //       "options": [
                  //             {
                  //                   "key": 2,
                  //                   "value": "Junior level"
                  //             },
                  //             {
                  //                   "key": 3,
                  //                   "value": "Mid level"
                  //             },
                  //             {
                  //                   "key": 5, // not sure abt this score
                  //                   "value": "Senior level"
                  //             },
                  //             {
                  //                   "key": 1,
                  //                   "value": "Any"
                  //             }
                  //       ],
                  // },
                  // {
                  //       // Repeat // automatically calculated
                  //       "question": "Personality type",
                  //       "columnName": "criteria1",
                  //       "maxScore": 3,
                  //       "options": [
                  //             {
                  //                   "key": 1,
                  //                   "value": "Leader"
                  //             },
                  //             {
                  //                   "key": 2,
                  //                   "value": "Extrovert"
                  //             },
                  //             {
                  //                   "key": 2,
                  //                   "value": "Introvert"
                  //             },
                  //             {
                  //                   "key": 3,
                  //                   "value": "Any"
                  //             }
                  //       ],
                  // },
                  // {
                  //       // Repeat // automatically calculated
                  //       "question": "Talent type",
                  //       "columnName": "criteria2",
                  //       "maxScore": 3,
                  //       "options": [
                  //             {
                  //                   "key": 1,
                  //                   "value": "Outstanding"
                  //             },
                  //             {
                  //                   "key": 2,
                  //                   "value": "Good"
                  //             },
                  //             {
                  //                   "key": 3,
                  //                   "value": "Any"
                  //             }
                  //       ],
                  // },
                  // {
                  //       // Repeat // automatically calculated
                  //       "question": "Experience",
                  //       "columnName": "criteria6",
                  //       "maxScore": 10,
                  //       "options": [
                  //             {
                  //                   "key": 10,
                  //                   "value": "Fresher"
                  //             },
                  //             {
                  //                   "key": 8,
                  //                   "value": "1 - 3 years"
                  //             },
                  //             {
                  //                   "key": 6,
                  //                   "value": "3 - 6 years"
                  //             },
                  //             {
                  //                   "key": 4,
                  //                   "value": "6 - 8 years"
                  //             },
                  //             {
                  //                   "key": 2,
                  //                   "value": "8 - 10 years"
                  //             },
                  //             {
                  //                   "key": 1,
                  //                   "value": "10+ years"
                  //             },
                  //       ],
                  // },

            ]

      }


}
