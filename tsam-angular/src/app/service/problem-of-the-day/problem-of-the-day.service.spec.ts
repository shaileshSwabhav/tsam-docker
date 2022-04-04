import { TestBed } from '@angular/core/testing';

import { ProblemOfTheDayService } from './problem-of-the-day.service';

describe('ProblemOfTheDayService', () => {
  let service: ProblemOfTheDayService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProblemOfTheDayService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
