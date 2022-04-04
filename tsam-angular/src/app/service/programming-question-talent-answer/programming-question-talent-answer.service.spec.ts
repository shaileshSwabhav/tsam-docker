import { TestBed } from '@angular/core/testing';

import { ProgrammingQuestionTalentAnswerService } from './programming-question-talent-answer.service';

describe('ProrammingSolutionService', () => {
  let service: ProgrammingQuestionTalentAnswerService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProgrammingQuestionTalentAnswerService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
