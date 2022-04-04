import { TestBed } from '@angular/core/testing';

import { ProgrammingQuestionService } from './programming-question.service';

describe('ProgrammingQuestionService', () => {
  let service: ProgrammingQuestionService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProgrammingQuestionService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
